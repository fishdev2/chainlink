package evm

import (
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// fetchedLog holds the log and the ID of the upkeep
type fetchedLog struct {
	id  *big.Int
	log logpoller.Log
}

// fetchedBlock holds the logs fetched for a block
type fetchedBlock struct {
	blockNumber int64
	logs        []fetchedLog
}

// Has returns true if the block has the log
func (b fetchedBlock) Has(id *big.Int, log logpoller.Log) bool {
	for _, l := range b.logs {
		if l.id.Cmp(id) != 0 {
			continue
		}
		if l.log.BlockNumber == log.BlockNumber && l.log.TxHash == log.TxHash && l.log.LogIndex == log.LogIndex {
			return true
		}
	}
	return false
}

// logEventBuffer is a circular/ring buffer of fetched logs.
// Each entry in the buffer represents a block,
// and holds the logs fetched for that block.
type logEventBuffer struct {
	lggr logger.Logger

	lock sync.RWMutex

	// size is the number of blocks supported by the buffer
	size int32
	// maxLogs is the number of logs allowed per block
	maxLogs int32
	// blocks is the circular buffer of fetched blocks
	blocks []fetchedBlock
	// latestBlock is the latest block number seen
	latestBlock int64
}

func newLogEventBuffer(lggr logger.Logger, size, maxLogs int) *logEventBuffer {
	return &logEventBuffer{
		lggr:    lggr.Named("KeepersRegistry.LogEventBuffer"),
		size:    int32(size),
		maxLogs: int32(maxLogs),
		blocks:  make([]fetchedBlock, size),
	}
}

func (b *logEventBuffer) latestBlockSeen() int64 {
	return atomic.LoadInt64(&b.latestBlock)
}

func (b *logEventBuffer) bufferSize() int {
	return int(atomic.LoadInt32(&b.size))
}

func (b *logEventBuffer) blockMaxLogs() int {
	return int(atomic.LoadInt32(&b.maxLogs))
}

// enqueue adds logs (if not exist) to the buffer, returning the number of logs added
func (b *logEventBuffer) enqueue(id *big.Int, logs ...logpoller.Log) int {
	b.lock.Lock()
	defer b.lock.Unlock()

	lggr := b.lggr.With("id", id.String())

	maxLogs := b.blockMaxLogs()

	latestBlock := b.latestBlock
	added := 0
	for _, log := range logs {
		if log.BlockNumber == 0 {
			continue
		}
		i := b.blockNumberIndex(log.BlockNumber)
		block := b.blocks[i]
		if block.blockNumber != log.BlockNumber {
			if block.blockNumber > 0 {
				lggr.Debugw("Overriding block", "currentBlock", block.blockNumber, "newBlock", log.BlockNumber)
			}
			block.blockNumber = log.BlockNumber
			block.logs = nil
		}
		if len(block.logs)+1 > maxLogs {
			lggr.Debugw("Reached max logs number, dropping log", "blockNumber", log.BlockNumber,
				"txHash", log.TxHash, "logIndex", log.LogIndex)
			continue
		}
		if !block.Has(id, log) {
			block.logs = append(block.logs, fetchedLog{id: id, log: log})
			b.blocks[i] = block
			added++
			if log.BlockNumber > latestBlock {
				latestBlock = log.BlockNumber
			}
		}
	}

	if latestBlock != b.latestBlock {
		b.latestBlock = latestBlock
	}
	if added > 0 {
		lggr.Debugw("Added logs to buffer", "addedLogs", added, "latestBlock", latestBlock)
	}

	return added
}

// dequeue returns the logs in range [latestBlock-blocks, latestBlock]
func (b *logEventBuffer) dequeue(blocks int) []fetchedLog {
	latestBlock := b.latestBlockSeen()
	if latestBlock == 0 {
		return nil
	}
	if blocks > int(latestBlock) {
		blocks = int(latestBlock) - 1
	}
	return b.dequeueRange(latestBlock-int64(blocks), latestBlock)
}

// dequeueRange returns the logs between start and end inclusive.
func (b *logEventBuffer) dequeueRange(start, end int64) []fetchedLog {
	b.lock.Lock()
	defer b.lock.Unlock()

	blocksInRange := b.getBlocksInRange(int(start), int(end))

	var results []fetchedLog
	for _, block := range blocksInRange {
		if block.blockNumber < start || block.blockNumber > end {
			continue
		}
		results = append(results, block.logs...)
		// clean logs once we read them
		block.logs = nil
		b.blocks[b.blockNumberIndex(block.blockNumber)] = block
	}

	return results
}

// getBlocksInRange returns the blocks between start and end.
// NOTE: this function should be called with the lock held
func (b *logEventBuffer) getBlocksInRange(start, end int) []fetchedBlock {
	var blocksInRange []fetchedBlock
	if end < start || end == 0 {
		// invalid range
		return blocksInRange
	}
	size := b.bufferSize()
	if start == 0 {
		// we reduce start by 1 to make it easier to calculate the index,
		// but we need to ensure we don't go below 0.
		start++
	}
	if start == end {
		// ensure we have at least one block in range
		end++
	}
	start = (start - 1) % size
	end = end % size
	if start < end {
		return b.blocks[start:end]
	}

	blocksInRange = append(blocksInRange, b.blocks[start:]...)
	blocksInRange = append(blocksInRange, b.blocks[:end]...)

	return blocksInRange
}

// blockNumberIndex returns the index of the block in the buffer
// NOTE: this function should be called with the lock held
func (t *logEventBuffer) blockNumberIndex(bn int64) int {
	return int(bn-1) % t.bufferSize()
}
