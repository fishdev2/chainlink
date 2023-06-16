package evm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
)

func TestLogEventBuffer_GetBlocksInRange(t *testing.T) {
	size := 3
	buf := newLogEventBuffer(logger.TestLogger(t), size, 10)

	buf.enqueue(big.NewInt(1), []logpoller.Log{
		{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
		{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
	}...)

	buf.enqueue(big.NewInt(2), []logpoller.Log{
		{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 2},
		{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 2},
	}...)

	tests := []struct {
		name string
		from int64
		to   int64
		want int
	}{
		{
			name: "all",
			from: 1,
			to:   3,
			want: 3,
		},
		{
			name: "partial",
			from: 1,
			to:   2,
			want: 2,
		},
		{
			name: "circular",
			from: 2,
			to:   4,
			want: 3,
		},
		{
			name: "zero start",
			from: 0,
			to:   2,
			want: 2,
		},
		{
			name: "invalid zero end",
			from: 0,
			to:   0,
		},
		{
			name: "invalid from larger than to",
			from: 4,
			to:   2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			blocks := buf.getBlocksInRange(int(tc.from), int(tc.to))
			require.Equal(t, tc.want, len(blocks))
			if tc.want > 0 {
				from := tc.from
				if from == 0 {
					from++
				}
				require.Equal(t, from, blocks[0].blockNumber)
				to := tc.to
				if to == 0 {
					to++
				} else if to > int64(size) {
					to = to % int64(size)
				}
				require.Equal(t, to, blocks[len(blocks)-1].blockNumber)
			}
		})
	}
}

func TestLogEventBuffer_EnqueueDequeue(t *testing.T) {
	t.Run("dequeue empty", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10)

		results := buf.dequeueRange(int64(1), int64(2))
		require.Equal(t, 0, len(results))
		results = buf.dequeue(2)
		require.Equal(t, 0, len(results))
	})

	t.Run("enqueue and dequeue range", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10)

		buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.lock.Lock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
		buf.lock.Unlock()
	})

	t.Run("enqueue overflow", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 2)

		require.Equal(t, 2, buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 1},
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 2},
		))
		buf.lock.Lock()
		require.Equal(t, 2, len(buf.blocks[0].logs))
		buf.lock.Unlock()
	})

	t.Run("enqueue and dequeue range twice", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10)

		buf.enqueue(big.NewInt(10),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 10},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 11},
		)
		buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		results := buf.dequeueRange(int64(1), int64(2))
		require.Equal(t, 2, len(results))
		verifyBlockNumbers(t, results, 1, 2)
		results = buf.dequeueRange(int64(1), int64(1))
		require.Equal(t, 0, len(results))
	})

	t.Run("enqueue and dequeue", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10)

		buf.enqueue(big.NewInt(10),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 10},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 11},
		)
		buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		results := buf.dequeue(8)
		require.Equal(t, 4, len(results))
		verifyBlockNumbers(t, results, 1, 2, 3, 3)
	})

	t.Run("enqueue and dequeue range circular", func(t *testing.T) {
		buf := newLogEventBuffer(logger.TestLogger(t), 3, 10)

		buf.enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 1, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
			logpoller.Log{BlockNumber: 3, TxHash: common.HexToHash("0x3"), LogIndex: 0},
		)
		buf.enqueue(big.NewInt(10),
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 10},
			logpoller.Log{BlockNumber: 4, TxHash: common.HexToHash("0x1"), LogIndex: 11},
		)

		results := buf.dequeueRange(int64(1), int64(1))
		require.Equal(t, 0, len(results))

		results = buf.dequeueRange(int64(3), int64(5))
		require.Equal(t, 3, len(results))
		verifyBlockNumbers(t, results, 3, 4, 4)
	})
}

func verifyBlockNumbers(t *testing.T, logs []fetchedLog, bns ...int64) {
	require.Equal(t, len(bns), len(logs))
	for i, log := range logs {
		require.Equal(t, bns[i], log.log.BlockNumber)
	}
}
