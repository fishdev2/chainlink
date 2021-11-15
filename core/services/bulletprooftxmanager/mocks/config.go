// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	big "math/big"

	chains "github.com/smartcontractkit/chainlink/core/chains"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Config is an autogenerated mock type for the Config type
type Config struct {
	mock.Mock
}

// BlockHistoryEstimatorBatchSize provides a mock function with given fields:
func (_m *Config) BlockHistoryEstimatorBatchSize() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// BlockHistoryEstimatorBlockDelay provides a mock function with given fields:
func (_m *Config) BlockHistoryEstimatorBlockDelay() uint16 {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// BlockHistoryEstimatorBlockHistorySize provides a mock function with given fields:
func (_m *Config) BlockHistoryEstimatorBlockHistorySize() uint16 {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// BlockHistoryEstimatorTransactionPercentile provides a mock function with given fields:
func (_m *Config) BlockHistoryEstimatorTransactionPercentile() uint16 {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// ChainType provides a mock function with given fields:
func (_m *Config) ChainType() chains.ChainType {
	ret := _m.Called()

	var r0 chains.ChainType
	if rf, ok := ret.Get(0).(func() chains.ChainType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(chains.ChainType)
	}

	return r0
}

// EthTxReaperInterval provides a mock function with given fields:
func (_m *Config) EthTxReaperInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// EthTxReaperThreshold provides a mock function with given fields:
func (_m *Config) EthTxReaperThreshold() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// EthTxResendAfterThreshold provides a mock function with given fields:
func (_m *Config) EthTxResendAfterThreshold() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// EvmEIP1559DynamicFees provides a mock function with given fields:
func (_m *Config) EvmEIP1559DynamicFees() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EvmFinalityDepth provides a mock function with given fields:
func (_m *Config) EvmFinalityDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmGasBumpPercent provides a mock function with given fields:
func (_m *Config) EvmGasBumpPercent() uint16 {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// EvmGasBumpThreshold provides a mock function with given fields:
func (_m *Config) EvmGasBumpThreshold() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// EvmGasBumpTxDepth provides a mock function with given fields:
func (_m *Config) EvmGasBumpTxDepth() uint16 {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// EvmGasBumpWei provides a mock function with given fields:
func (_m *Config) EvmGasBumpWei() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// EvmGasFeeCap provides a mock function with given fields:
func (_m *Config) EvmGasFeeCap() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// EvmGasLimitDefault provides a mock function with given fields:
func (_m *Config) EvmGasLimitDefault() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// EvmGasLimitMultiplier provides a mock function with given fields:
func (_m *Config) EvmGasLimitMultiplier() float32 {
	ret := _m.Called()

	var r0 float32
	if rf, ok := ret.Get(0).(func() float32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float32)
	}

	return r0
}

// EvmGasPriceDefault provides a mock function with given fields:
func (_m *Config) EvmGasPriceDefault() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// EvmGasTipCapDefault provides a mock function with given fields:
func (_m *Config) EvmGasTipCapDefault() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// EvmGasTipCapMinimum provides a mock function with given fields:
func (_m *Config) EvmGasTipCapMinimum() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// EvmMaxGasPriceWei provides a mock function with given fields:
func (_m *Config) EvmMaxGasPriceWei() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// EvmMaxInFlightTransactions provides a mock function with given fields:
func (_m *Config) EvmMaxInFlightTransactions() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmMaxQueuedTransactions provides a mock function with given fields:
func (_m *Config) EvmMaxQueuedTransactions() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// EvmMinGasPriceWei provides a mock function with given fields:
func (_m *Config) EvmMinGasPriceWei() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// EvmNonceAutoSync provides a mock function with given fields:
func (_m *Config) EvmNonceAutoSync() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EvmRPCDefaultBatchSize provides a mock function with given fields:
func (_m *Config) EvmRPCDefaultBatchSize() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// GasEstimatorMode provides a mock function with given fields:
func (_m *Config) GasEstimatorMode() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// KeySpecificMaxGasPriceWei provides a mock function with given fields: addr
func (_m *Config) KeySpecificMaxGasPriceWei(addr common.Address) *big.Int {
	ret := _m.Called(addr)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(common.Address) *big.Int); ok {
		r0 = rf(addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// TriggerFallbackDBPollInterval provides a mock function with given fields:
func (_m *Config) TriggerFallbackDBPollInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}
