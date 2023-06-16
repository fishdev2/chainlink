package txmgr

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
)

// Config encompasses config used by txmgr package
// Unless otherwise specified, these should support changing at runtime
//
//go:generate mockery --quiet --recursive --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go
type Config interface {
	ChainType() coreconfig.ChainType
	EvmFinalityDepth() uint32
	EvmNonceAutoSync() bool
	EvmRPCDefaultBatchSize() uint32
	KeySpecificMaxGasPriceWei(addr common.Address) *assets.Wei
}

type FeeConfig interface {
	EIP1559DynamicFees() bool
	BumpPercent() uint16
	BumpThreshold() uint64
	BumpTxDepth() uint32
	LimitDefault() uint32
	PriceDefault() *assets.Wei
	TipCapMin() *assets.Wei
	PriceMax() *assets.Wei
	PriceMin() *assets.Wei
}

type DatabaseConfig interface {
	DefaultQueryTimeout() time.Duration
	LogSQL() bool
}

type ListenerConfig interface {
	FallbackPollInterval() time.Duration
}

type (
	EvmTxmConfig         txmgrtypes.TransactionManagerConfig
	EvmTxmFeeConfig      txmgrtypes.TransactionManagerFeeConfig
	EvmBroadcasterConfig txmgrtypes.BroadcasterConfig
	EvmConfirmerConfig   txmgrtypes.ConfirmerConfig
	EvmResenderConfig    txmgrtypes.ResenderConfig
	EvmReaperConfig      txmgrtypes.ReaperConfig
)

var _ EvmTxmConfig = (*evmTxmConfig)(nil)

type evmTxmConfig struct {
	Config
}

func NewEvmTxmConfig(c Config) *evmTxmConfig {
	return &evmTxmConfig{c}
}

func (c evmTxmConfig) SequenceAutoSync() bool { return c.EvmNonceAutoSync() }

func (c evmTxmConfig) IsL2() bool { return c.ChainType().IsL2() }

func (c evmTxmConfig) RPCDefaultBatchSize() uint32 { return c.EvmRPCDefaultBatchSize() }

func (c evmTxmConfig) FinalityDepth() uint32 { return c.EvmFinalityDepth() }

var _ EvmTxmFeeConfig = (*evmTxmFeeConfig)(nil)

type evmTxmFeeConfig struct {
	FeeConfig
}

func NewEvmTxmFeeConfig(c FeeConfig) *evmTxmFeeConfig {
	return &evmTxmFeeConfig{c}
}

func (c evmTxmFeeConfig) MaxFeePrice() string { return c.PriceMax().String() }

func (c evmTxmFeeConfig) FeePriceDefault() string { return c.PriceDefault().String() }

func (c evmTxmFeeConfig) FeeBumpTxDepth() uint32 { return c.BumpTxDepth() }

func (c evmTxmFeeConfig) FeeLimitDefault() uint32 { return c.LimitDefault() }

func (c evmTxmFeeConfig) FeeBumpThreshold() uint64 { return c.BumpThreshold() }

func (c evmTxmFeeConfig) FeeBumpPercent() uint16 { return c.BumpPercent() }
