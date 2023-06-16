package evm

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// TODO: configurable values or based on block time
const (
	// logRetention is the amount of time to retain logs for.
	// 5 minutes is the desired retention time for logs, but we add an extra 10 minutes buffer.
	logRetention              = (time.Minute * 5) + (time.Minute * 10)
	allowedLogsPerBlock       = 100
	logBlocksLookback   int64 = 256
	lookbackBuffer      int64 = 10
)

// TODO: configurable values or based on block time
var (
	blockRateLimit  = rate.Every(time.Second)
	blockLimitBurst = 32
	logsRateLimit   = rate.Every(time.Second)
	logsLimitBurst  = 4
	queryWorkers    = 4
)

// LogTriggerConfig is an alias for log trigger config.
type LogTriggerConfig = i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig

// upkeepFilterEntry holds the upkeep filter, rate limiter and last polled block.
type upkeepFilterEntry struct {
	id     *big.Int
	filter logpoller.Filter
	cfg    LogTriggerConfig
	// lastPollBlock is the last block number the logs were fetched for this upkeep
	lastPollBlock int64
	// blockLimiter is used to limit the number of blocks to fetch logs for an upkeep
	blockLimiter *rate.Limiter
}

// logEventProvider manages log filters for upkeeps and enables to read the log events.
type logEventProvider struct {
	lggr logger.Logger

	poller logpoller.LogPoller

	lock   sync.RWMutex
	active map[string]upkeepFilterEntry

	buffer *logEventBuffer
}

func NewLogEventProvider(lggr logger.Logger, poller logpoller.LogPoller) *logEventProvider {
	return &logEventProvider{
		lggr:   lggr.Named("KeepersRegistry.LogEventProvider"),
		buffer: newLogEventBuffer(lggr, int(logBlocksLookback)*3, allowedLogsPerBlock),
		poller: poller,
		lock:   sync.RWMutex{},
		active: make(map[string]upkeepFilterEntry),
	}
}

// Register creates a filter from the given upkeep and calls log poller to register it.
func (p *logEventProvider) RegisterFilter(upkeepID *big.Int, cfg LogTriggerConfig) error {
	if err := p.validateLogTriggerConfig(cfg); err != nil {
		return errors.Wrap(err, "invalid log trigger config")
	}
	filter := p.newLogFilter(upkeepID, cfg)

	// TODO: optimize locking, currently we lock the whole map while registering the filter
	p.lock.Lock()
	defer p.lock.Unlock()

	uid := upkeepID.String()
	if _, ok := p.active[uid]; ok {
		// TODO: check for updates
		return errors.Errorf("filter for upkeep with id %s already registered", uid)
	}
	if err := p.poller.RegisterFilter(filter); err != nil {
		return errors.Wrap(err, "failed to register upkeep filter")
	}
	p.active[uid] = upkeepFilterEntry{
		id:           upkeepID,
		filter:       filter,
		cfg:          cfg,
		blockLimiter: rate.NewLimiter(blockRateLimit, blockLimitBurst),
	}

	return nil
}

// Unregister removes the filter for the given upkeepID
func (p *logEventProvider) UnregisterFilter(upkeepID *big.Int) error {
	err := p.poller.UnregisterFilter(p.filterName(upkeepID), nil)
	if err == nil {
		p.lock.Lock()
		delete(p.active, upkeepID.String())
		p.lock.Unlock()
	}
	return errors.Wrap(err, "failed to unregister upkeep filter")
}

// GetLogs returns the logs in the given range.
func (p *logEventProvider) GetLogs(from, to int64) map[string][]logpoller.Log {
	logs := p.buffer.dequeueRange(from, to)

	results := make(map[string][]logpoller.Log, 0)
	for _, l := range logs {
		uid := l.id.String()
		upkeepLogs, ok := results[uid]
		if !ok {
			upkeepLogs = []logpoller.Log{}
		}
		upkeepLogs = append(upkeepLogs, l.log)
		results[uid] = upkeepLogs
	}

	return results
}

// FetchLogs fetches the logs for the given upkeeps.
func (p *logEventProvider) FetchLogs(ctx context.Context, force bool, ids ...*big.Int) error {
	latest, err := p.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}
	entries := p.getEntries(latest, force, ids...)

	p.lggr.Debugw("getting logs for entries", "latestBlock", latest, "entries", len(entries))

	err = p.fetchLogsConcurrent(latest, entries)
	if err != nil {
		return fmt.Errorf("failed to get logs: %w", err)
	}
	p.updateEntriesLastPoll(entries)

	return nil
}

func (p *logEventProvider) updateEntriesLastPoll(entries []*upkeepFilterEntry) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, entry := range entries {
		// for successful queries, the last poll block was updated
		orig := p.active[entry.id.String()]
		if entry.lastPollBlock == orig.lastPollBlock {
			continue
		}
		orig.lastPollBlock = entry.lastPollBlock
		p.active[entry.id.String()] = orig
	}
}

// getEntries returns the filters for the given upkeepIDs,
// returns empty filter for inactive upkeeps.
//
// TODO: group filters by contract address?
func (p *logEventProvider) getEntries(latestBlock int64, force bool, ids ...*big.Int) []*upkeepFilterEntry {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var filters []*upkeepFilterEntry
	for _, id := range ids {
		entry, ok := p.active[id.String()]
		if !ok { // entry not found, could be inactive upkeep
			p.lggr.Debugw("upkeep filter not found", "upkeep", id.String())
			filters = append(filters, &upkeepFilterEntry{id: id})
			continue
		}
		if !force && entry.lastPollBlock > latestBlock {
			p.lggr.Debugw("already polled latest block", "entry.lastPollBlock", entry.lastPollBlock, "latestBlock", latestBlock, "upkeep", id.String())
			filters = append(filters, &upkeepFilterEntry{id: id, lastPollBlock: entry.lastPollBlock})
			continue
		}
		// recreating the struct to be thread safe
		filters = append(filters, &upkeepFilterEntry{
			id:            id,
			filter:        p.newLogFilter(id, entry.cfg),
			lastPollBlock: entry.lastPollBlock,
			blockLimiter:  entry.blockLimiter,
		})
	}

	return filters
}

func (p *logEventProvider) fetchLogsConcurrent(latest int64, entries []*upkeepFilterEntry) (err error) {
	var wg sync.WaitGroup
	// using a set of worker goroutines to fetch logs for upkeeps
	for i := 0; i < len(entries); i += queryWorkers {
		end := i + queryWorkers
		if end > len(entries) {
			end = len(entries)
		}
		wg.Add(1)
		go func(i, end int) {
			defer wg.Done()
			errGetLogs := p.fetchLogs(latest, entries[i:end]...)
			if err != nil {
				p.lggr.Debugw("failed to get logs", "i", i, "end", end, "err", err)
				multierr.Append(err, errGetLogs)
			}
		}(i, end)
	}
	wg.Wait()

	return nil
}

// fetchLogs calls log poller to get the logs for the given upkeep entries.
func (p *logEventProvider) fetchLogs(latest int64, entries ...*upkeepFilterEntry) (err error) {
	mainLggr := p.lggr.With("latestBlock", latest)

	for _, entry := range entries {
		if len(entry.filter.Addresses) == 0 {
			continue
		}
		lggr := mainLggr.With("upkeep", entry.id.String(), "addrs", entry.filter.Addresses, "sigs", entry.filter.EventSigs)
		start := entry.lastPollBlock
		if start == 0 { // first time polling
			start = latest - logBlocksLookback
			entry.blockLimiter.SetBurst(int(logBlocksLookback + 1))
		}
		start = start - lookbackBuffer // adding a buffer to avoid missing logs
		if start < 0 {
			start = 0
		}
		resv := entry.blockLimiter.ReserveN(time.Now(), int(latest-start))
		if !resv.OK() {
			lggr.Warnw("log upkeep block limit exceeded")
			multierr.Append(err, fmt.Errorf("log upkeep block limit exceeded for upkeep %s", entry.id.String()))
			continue
		}
		lggr = lggr.With("startBlock", start)
		// TODO: TBD what function to use to get logs
		logs, err := p.poller.LogsWithSigs(start, latest, entry.filter.EventSigs, entry.filter.Addresses[0])
		if err != nil {
			resv.Cancel() // cancels limit reservation as we failed to get logs
			lggr.Warnw("failed to get logs", "err", err)
			multierr.Append(err, fmt.Errorf("failed to get logs for upkeep %s: %w", entry.id.String(), err))
			continue
		}

		// TODO: filter unfinalized logs if we fetch logs from the latest block

		// if this limiter's burst was set to the max,
		// we need to reset it
		if entry.blockLimiter.Burst() == int(logBlocksLookback+1) {
			entry.blockLimiter.SetBurst(blockLimitBurst)
		}
		added := p.buffer.enqueue(entry.id, logs...)
		// if we added logs or couldn't find, update the last poll block
		if added > 0 || len(logs) == 0 {
			entry.lastPollBlock = latest
		}
	}

	return nil
}

// newLogFilter creates logpoller.Filter from the given upkeep config
func (p *logEventProvider) newLogFilter(upkeepID *big.Int, cfg LogTriggerConfig) logpoller.Filter {
	sigs := p.getFiltersBySelector(cfg.FilterSelector, cfg.Topic1[:], cfg.Topic2[:], cfg.Topic3[:])
	sigs = append([]common.Hash{common.BytesToHash(cfg.Topic0[:])}, sigs...)
	return logpoller.Filter{
		Name:      p.filterName(upkeepID),
		EventSigs: sigs,
		Addresses: []common.Address{cfg.ContractAddress},
		Retention: logRetention,
	}
}

func (p *logEventProvider) validateLogTriggerConfig(cfg LogTriggerConfig) error {
	var zeroAddr common.Address
	var zeroBytes [32]byte
	if bytes.Equal(cfg.ContractAddress[:], zeroAddr[:]) {
		return errors.New("invalid contract address: zeroed")
	}
	if bytes.Equal(cfg.Topic0[:], zeroBytes[:]) {
		return errors.New("invalid topic0: zeroed")
	}
	return nil
}

// getFiltersBySelector the filters based on the filterSelector
func (p *logEventProvider) getFiltersBySelector(filterSelector uint8, filters ...[]byte) []common.Hash {
	var sigs []common.Hash
	var zeroBytes [32]byte
	for i, f := range filters {
		// bitwise AND the filterSelector with the index to check if the filter is needed
		mask := uint8(1 << uint8(i))
		a := filterSelector & mask
		if a == uint8(0) {
			continue
		}
		if bytes.Equal(f, zeroBytes[:]) {
			continue
		}
		sigs = append(sigs, common.BytesToHash(common.LeftPadBytes(f, 32)))
	}
	return sigs
}

func (p *logEventProvider) filterName(upkeepID *big.Int) string {
	return logpoller.FilterName("KeepersRegistry LogUpkeep", upkeepID.String())
}
