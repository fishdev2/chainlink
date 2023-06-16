package types

import "time"

type TransactionManagerConfig interface {
	BroadcasterConfig
	ConfirmerConfig
	ResenderConfig
	ReaperConfig

	SequenceAutoSync() bool
}

type TransactionManagerFeeConfig interface {
	BroadcasterFeeConfig
	ConfirmerFeeConfig
}

type TransactionManagerTransactionsConfig interface {
	BroadcasterTransactionsConfig
	ConfirmerTransactionsConfig
	ResenderTransactionsConfig
	ReaperTransactionsConfig

	ForwardersEnabled() bool
	MaxQueued() uint64
}

type BroadcasterConfig interface {
	IsL2() bool
}

type BroadcasterFeeConfig interface {
	MaxFeePrice() string     // logging value
	FeePriceDefault() string // logging value
}

type BroadcasterTransactionsConfig interface {
	MaxInFlight() uint32
}

type BroadcasterListenerConfig interface {
	FallbackPollInterval() time.Duration
}

type ConfirmerConfig interface {
	RPCDefaultBatchSize() uint32
	FinalityDepth() uint32
}

type ConfirmerFeeConfig interface {
	FeeBumpTxDepth() uint32
	FeeLimitDefault() uint32

	// from gas.Config
	FeeBumpThreshold() uint64
	MaxFeePrice() string // logging value
	FeeBumpPercent() uint16
}

type ConfirmerDatabaseConfig interface {
	// from pg.QConfig
	DefaultQueryTimeout() time.Duration
}

type ConfirmerTransactionsConfig interface {
	MaxInFlight() uint32
	ForwardersEnabled() bool
}

type ResenderConfig interface {
	RPCDefaultBatchSize() uint32
}

type ResenderTransactionsConfig interface {
	ResendAfterThreshold() time.Duration
	MaxInFlight() uint32
}

//go:generate mockery --quiet --name ReaperConfig --output ./mocks/ --case=underscore

// ReaperConfig is the config subset used by the reaper
type ReaperConfig interface {
	FinalityDepth() uint32
}

type ReaperTransactionsConfig interface {
	ReaperInterval() time.Duration
	ReaperThreshold() time.Duration
}
