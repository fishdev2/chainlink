package functions

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// OnchainAllowlist maintains an allowlist of addresses that if fetched from the blockchain (EVM-only).
// Updating the allowlist can be achieved by either manually calling the UpdateFromChain method
// or starting the service (Start()/Close()), which will periodically update the allowlist.
//
//go:generate mockery --quiet --name OnchainAllowlist --output ./mocks/ --case=underscore
type OnchainAllowlist interface {
	job.ServiceCtx

	Allow(common.Address) bool
	UpdateFromChain() error
}

type onchainAllowlist struct {
	client          evmclient.Client
	allowlist       map[common.Address]struct{}
	mu              sync.RWMutex
	contractAddress common.Address
	abiReturnType   abi.Arguments
	methodSignature []byte
	updateFrequency time.Duration
	shutdownCh      chan struct{}
	lggr            logger.Logger
}

const (
	AllowlistArgumentType = "address[]"
	AllowlistMethodHeader = "getAuthorizedSenders()"
)

func NewOnchainAllowlist(client evmclient.Client, contractAddress common.Address, updateFrequency time.Duration, lggr logger.Logger) (OnchainAllowlist, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	abiType, err := abi.NewType(AllowlistArgumentType, "", []abi.ArgumentMarshaling{})
	if err != nil {
		return nil, fmt.Errorf("unexpected error during abi.NewType: %s", err)
	}
	abiArgs := abi.Arguments([]abi.Argument{{Name: "addresses", Type: abiType}})
	methodSignature := crypto.Keccak256([]byte(AllowlistMethodHeader))[:4]
	return &onchainAllowlist{
		client:          client,
		allowlist:       make(map[common.Address]struct{}),
		contractAddress: contractAddress,
		abiReturnType:   abiArgs,
		methodSignature: methodSignature,
		updateFrequency: updateFrequency,
		shutdownCh:      make(chan struct{}),
		lggr:            lggr.Named("OnchainAllowlist"),
	}, nil
}

func (a *onchainAllowlist) Allow(address common.Address) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, ok := a.allowlist[address]
	return ok
}

func (a *onchainAllowlist) UpdateFromChain() error {
	rawResponse, err := a.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &a.contractAddress,
		Data: a.methodSignature,
	}, big.NewInt(-1))
	if err != nil {
		a.lggr.Errorw("error calling CallContract", "err", err)
		return err
	}
	unpacked, err := a.abiReturnType.Unpack(rawResponse)
	if err != nil || len(unpacked) != 1 {
		a.lggr.Errorw("error unpacking", "err", err, "len", len(unpacked))
		return err
	}
	addrList, ok := unpacked[0].([]common.Address)
	if !ok {
		a.lggr.Error("error casting to []common.Address")
		return errors.New("error casting to []common.Address")
	}
	newAllowlist := make(map[common.Address]struct{})
	for _, addr := range addrList {
		newAllowlist[addr] = struct{}{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.allowlist = newAllowlist
	a.lggr.Info("allowlist updated successfully", "len", len(addrList))
	return nil
}

func (a *onchainAllowlist) updatePeriodically() {
	ticker := time.NewTicker(a.updateFrequency)
	defer ticker.Stop()
	for {
		select {
		case <-a.shutdownCh:
			return
		case <-ticker.C:
			a.UpdateFromChain()
		}
	}
}

func (a *onchainAllowlist) Start(context.Context) error {
	go a.updatePeriodically()
	return nil
}

func (a *onchainAllowlist) Close() error {
	a.shutdownCh <- struct{}{}
	return nil
}
