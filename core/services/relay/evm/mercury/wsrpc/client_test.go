package wsrpc

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	mocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
)

func Test_Client_Transmit(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	req := &pb.TransmitRequest{}

	t.Run("sends on reset channel after MaxConsecutiveTransmitFailures timed out transmits", func(t *testing.T) {
		calls := 0
		wsrpcClient := &mocks.MockWSRPCClient{
			TransmitF: func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
				calls++
				return nil, context.DeadlineExceeded
			},
		}
		conn := &mocks.MockConn{
			Ready: true,
		}
		c := newClient(lggr, csakey.KeyV2{}, nil, "")
		c.conn = conn
		c.client = wsrpcClient
		c.StartOnce("Mock WSRPC Client", func() error { return nil })
		for i := 1; i < MaxConsecutiveTransmitFailures; i++ {
			_, err := c.Transmit(ctx, req)
			require.EqualError(t, err, "context deadline exceeded")
		}
		assert.Equal(t, 4, calls)
		select {
		case <-c.chResetTransport:
			t.Fatal("unexpected send on chResetTransport")
		default:
		}
		_, err := c.Transmit(ctx, req)
		require.EqualError(t, err, "context deadline exceeded")
		assert.Equal(t, 5, calls)
		select {
		case <-c.chResetTransport:
		default:
			t.Fatal("expected send on chResetTransport")
		}
	})
}
func Test_WSRPCClient_LatestReport(t *testing.T) {
	t.Fatal("TODO")
}
