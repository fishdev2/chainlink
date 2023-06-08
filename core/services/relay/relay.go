package relay

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type Network string

var (
	EVM             Network = "evm"
	Cosmos          Network = "cosmos"
	Solana          Network = "solana"
	StarkNet        Network = "starknet"
	SupportedRelays         = map[Network]struct{}{
		EVM:      {},
		Cosmos:   {},
		Solana:   {},
		StarkNet: {},
	}
)

type LogPollerCapable interface {
	RegisterLogFilters(args types.RelayArgs, q pg.Queryer) (err error)
	UnregisterLogFilters(args types.RelayArgs, q pg.Queryer) (err error)
}

type JobHooks interface {
	OnCreateJob(arg types.RelayArgs, q pg.Queryer) (err error)
	OnDeleteJob(arg types.RelayArgs, q pg.Queryer) (err error)
}

// RelayerExt is a subset of [loop.Relayer] for adapting [types.Relayer], typically with a ChainSet. See [relayerAdapter].
type RelayerExt interface {
	services.ServiceCtx

	ChainStatus(ctx context.Context, id string) (types.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error)

	NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error)

	SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error
}

var _ loop.Relayer = (*relayerAdapter)(nil)

// relayerAdapter adapts a [types.Relayer] and [RelayerExt] to implement [loop.Relayer].
type relayerAdapter struct {
	types.Relayer
	JobHooks
	RelayerExt
}

// NewRelayerAdapter returns a [loop.Relayer] adapted from a [types.Relayer] and [RelayerExt].
func NewRelayerAdapter(r types.Relayer, e RelayerExt) loop.Relayer {
	return &relayerAdapter{Relayer: r, RelayerExt: e}
}

func (r *relayerAdapter) NewConfigProvider(ctx context.Context, rargs types.RelayArgs) (types.ConfigProvider, error) {
	return r.Relayer.NewConfigProvider(rargs)
}

func (r *relayerAdapter) NewMedianProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MedianProvider, error) {
	return r.Relayer.NewMedianProvider(rargs, pargs)
}

func (r *relayerAdapter) NewMercuryProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MercuryProvider, error) {
	return r.Relayer.NewMercuryProvider(rargs, pargs)
}

func (r *relayerAdapter) Start(ctx context.Context) error {
	var ms services.MultiStart
	return ms.Start(ctx, r.RelayerExt, r.Relayer)
}

func (r *relayerAdapter) Close() error {
	return services.CloseAll(r.Relayer, r.RelayerExt)
}

func (r *relayerAdapter) Name() string {
	return fmt.Sprintf("%s-%s", r.Relayer.Name(), r.RelayerExt.Name())
}

func (r *relayerAdapter) Ready() (err error) {
	return errors.Join(r.Relayer.Ready(), r.RelayerExt.Ready())
}

func (r *relayerAdapter) HealthReport() map[string]error {
	hr := make(map[string]error)
	maps.Copy(r.Relayer.HealthReport(), hr)
	maps.Copy(r.RelayerExt.HealthReport(), hr)
	return hr
}

type ErrLogFiltersNotSupported struct {
	relayName string
}

func (e ErrLogFiltersNotSupported) Error() string {
	return fmt.Sprintf("Log filtering is not supported by relay %s", e.relayName)
}

func (r *relayerAdapter) OnCreateJob(rargs types.RelayArgs, q pg.Queryer) (err error) {
	relay, ok := r.Relayer.(LogPollerCapable)
	if !ok {
		return nil
	}
	return relay.RegisterLogFilters(rargs, q)

}

func (r *relayerAdapter) OnDeleteJob(rargs types.RelayArgs, q pg.Queryer) (err error) {
	relay, ok := r.Relayer.(LogPollerCapable)
	if !ok {
		return nil
	}
	return relay.UnregisterLogFilters(rargs, q)
}
