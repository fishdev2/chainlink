package functions

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type FunctionsHandlerConfig struct {
	AllowlistCheckEnabled       bool   `json:"allowlistCheckEnabled"`
	AllowlistChainID            int64  `json:"allowlistChainID"`
	AllowlistContractAddress    string `json:"allowlistContractAddress"`
	AllowlistUpdateFrequencySec int    `json:"allowlistUpdateFrequencySec"`
}

type functionsHandler struct {
	handlerConfig *FunctionsHandlerConfig
	donConfig     *config.DONConfig
	don           handlers.DON
	allowlist     OnchainAllowlist
	lggr          logger.Logger
}

var _ handlers.Handler = (*functionsHandler)(nil)

func NewFunctionsHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, chains evm.ChainSet, lggr logger.Logger) (handlers.Handler, error) {
	var cfg FunctionsHandlerConfig
	if err := json.Unmarshal(handlerConfig, &cfg); err != nil {
		return nil, err
	}
	var allowlist OnchainAllowlist
	if cfg.AllowlistCheckEnabled {
		chain, err := chains.Get(big.NewInt(cfg.AllowlistChainID))
		if err != nil {
			return nil, err
		}
		if !common.IsHexAddress(cfg.AllowlistContractAddress) {
			return nil, errors.New("allowlistContractAddress is not a valid hex address")
		}
		if cfg.AllowlistUpdateFrequencySec <= 0 {
			return nil, errors.New("allowlistUpdateFrequencySec must be positive")
		}
		checkFreq := time.Duration(cfg.AllowlistUpdateFrequencySec) * time.Second
		allowlist, err = NewOnchainAllowlist(chain.Client(), common.HexToAddress(cfg.AllowlistContractAddress), checkFreq, lggr)
		if err != nil {
			return nil, err
		}
	}
	return &functionsHandler{
		handlerConfig: &cfg,
		donConfig:     donConfig,
		don:           don,
		allowlist:     allowlist,
		lggr:          lggr,
	}, nil
}

func (h *functionsHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	if err := msg.Validate(); err != nil {
		h.lggr.Debug("received invalid message", "err", err)
		return err
	}
	sender := common.HexToAddress(msg.Body.Sender)
	if h.allowlist != nil && !h.allowlist.Allow(sender) {
		h.lggr.Debug("received a message from a non-allowlisted address", "sender", msg.Body.Sender)
		return errors.New("sender not allowlisted")
	}
	h.lggr.Debug("received a valid message", "sender", msg.Body.Sender)
	return nil
}

func (h *functionsHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	return nil
}

func (h *functionsHandler) Start(ctx context.Context) error {
	return h.allowlist.Start(ctx)
}

func (h *functionsHandler) Close() error {
	return h.allowlist.Close()
}
