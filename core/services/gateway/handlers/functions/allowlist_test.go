package functions_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
)

func TestAllowlist_Basic(t *testing.T) {
	t.Parallel()

	addr1 := "9ed925d8206a4f88a2f643b28b3035b315753cd6"
	addr2 := "ea6721ac65bced841b8ec3fc5fedea6141a0ade4"
	addr3 := "84689acc87ff22841b8ec378300da5e141a99911"
	abiEncodedAddresses :=
		"0000000000000000000000000000000000000000000000000000000000000020" +
			"0000000000000000000000000000000000000000000000000000000000000002" +
			"000000000000000000000000" + addr1 +
			"000000000000000000000000" + addr2
	rawData, err := hex.DecodeString(abiEncodedAddresses)
	require.NoError(t, err)

	client := mocks.NewClient(t)
	client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(rawData, nil)
	allowlist, err := functions.NewOnchainAllowlist(client, common.Address{}, 0, logger.TestLogger(t))
	require.NoError(t, err)

	require.NoError(t, allowlist.UpdateFromChain())
	require.False(t, allowlist.Allow(common.Address{}))
	require.True(t, allowlist.Allow(common.HexToAddress(addr1)))
	require.True(t, allowlist.Allow(common.HexToAddress(addr2)))
	require.False(t, allowlist.Allow(common.HexToAddress(addr3)))
}
