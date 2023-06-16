package wsrpc

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/connectivity"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const MaxConsecutiveTransmitFailures = 5

type Client interface {
	services.ServiceCtx
	pb.MercuryClient
}

type Conn interface {
	WaitForReady(ctx context.Context) bool
	GetState() connectivity.State
	Close()
}

type client struct {
	utils.StartStopOnce

	csaKey       csakey.KeyV2
	serverPubKey []byte
	serverURL    string

	logger logger.Logger
	conn   Conn
	client pb.MercuryClient

	consecutiveTimeoutCnt atomic.Int32
	wg                    sync.WaitGroup
	chStop                utils.StopChan
	chResetTransport      chan struct{}
}

// Consumers of wsrpc package should not usually call NewClient directly, but instead use the Pool
func NewClient(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) Client {
	return newClient(lggr, clientPrivKey, serverPubKey, serverURL)
}

func newClient(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) *client {
	return &client{
		csaKey:           clientPrivKey,
		serverPubKey:     serverPubKey,
		serverURL:        serverURL,
		logger:           lggr.Named("WSRPC"),
		chResetTransport: make(chan struct{}, 1),
	}
}

func (w *client) Start(_ context.Context) error {
	return w.StartOnce("WSRPC Client", func() error {
		if err := w.dial(context.Background()); err != nil {
			return err
		}
		w.wg.Add(1)
		go w.runloop()
		return nil
	})
}

// NOTE: Dial is non-blocking, and will retry on an exponential backoff
// in the background until close is called, or context is cancelled.
// This is why we use the background context, not the start context here.
//
// Any transmits made while client is still trying to dial will fail
// with error.
func (w *client) dial(ctx context.Context, opts ...wsrpc.DialOption) error {
	conn, err := wsrpc.DialWithContext(ctx, w.serverURL,
		append(opts,
			wsrpc.WithTransportCreds(w.csaKey.Raw().Bytes(), w.serverPubKey),
			wsrpc.WithLogger(w.logger),
		)...,
	)
	if err != nil {
		setLivenessMetric(false)
		return errors.Wrap(err, "failed to dial wsrpc client")
	}
	setLivenessMetric(true)
	w.conn = conn
	w.client = pb.NewMercuryClient(conn)
	return nil
}

func (w *client) runloop() {
	defer w.wg.Done()
	for {
		select {
		case <-w.chStop:
			return
		case <-w.chResetTransport:
			// Using channel here ensures we only have one reset in process at
			// any given time
			w.resetTransport()
		}
	}
}

// resetTransport disconnects and reconnects to the mercury server
func (w *client) resetTransport() {
	ok := w.IfStarted(func() {
		w.conn.Close() // Close is safe to call multiple times
	})
	if !ok {
		panic("resetTransport should never be called unless client is in 'started' state")
	}
	ctx, cancel := w.chStop.Ctx(context.Background())
	defer cancel()
	b := utils.NewRedialBackoff()
	for {
		// Will block until successful dial, or context is canceled (i.e. on close)
		if err := w.dial(ctx, wsrpc.WithBlock()); err != nil {
			w.logger.Errorw("ResetTransport failed to redial", "err", err)
			time.Sleep(b.Duration())
		} else {
			w.logger.Info("ResetTransport successfully redialled")
			break
		}
	}
}

func (w *client) Close() error {
	return w.StopOnce("WSRPC Client", func() error {
		close(w.chStop)
		w.conn.Close()
		w.wg.Wait()
		return nil
	})
}

func (w *client) Name() string {
	return "EVM.Mercury.WSRPCClient"
}

func (w *client) HealthReport() map[string]error {
	return map[string]error{w.Name(): w.Healthy()}
}

// Healthy if connected
func (w *client) Healthy() (err error) {
	if err = w.StartStopOnce.Healthy(); err != nil {
		return err
	}
	state := w.conn.GetState()
	if state != connectivity.Ready {
		return errors.Errorf("client state should be %s; got %s", connectivity.Ready, state)
	}
	return nil
}

func (w *client) waitForReady(ctx context.Context) (err error) {
	ok := w.IfStarted(func() {
		if ready := w.conn.WaitForReady(ctx); !ready {
			err = errors.Errorf("websocket client not ready; got state: %v", w.conn.GetState())
			return
		}
	})
	if !ok {
		return errors.New("client is not started")
	}
	return
}

func (w *client) Transmit(ctx context.Context, req *pb.TransmitRequest) (resp *pb.TransmitResponse, err error) {
	lggr := w.logger.With("req.Payload", hexutil.Encode(req.Payload))
	lggr.Debug("Transmit")
	start := time.Now()
	if err := w.waitForReady(ctx); err != nil {
		return nil, errors.Wrap(err, "Transmit failed")
	}
	resp, err = w.client.Transmit(ctx, req)
	if errors.Is(err, context.DeadlineExceeded) {
		cnt := w.consecutiveTimeoutCnt.Add(1)
		if cnt == MaxConsecutiveTransmitFailures {
			w.logger.Errorf("Timed out on %d consecutive transmits, resetting transport", cnt)
			// HACK: if we get 5+ request timeouts in a row, close
			// and re-open the websocket connection. This *shouldn't* be
			// necessary in theory (wsrpc should handle it for us) but it acts
			// as a "belts and braces" approach to ensure we get a websocket
			// connection back up and running again if it gets itself into a
			// bad state
			select {
			case w.chResetTransport <- struct{}{}:
			default:
				// TODO: how to handle the case where a reset is already happening? How can this even happen?
			}
		}
	} else {
		w.consecutiveTimeoutCnt.Store(0)
	}
	if err != nil {
		lggr.Warnw("Transmit failed", "err", err, "req", req, "resp", resp)
		incRequestStatusMetric(statusFailed)
	} else {
		lggr.Debugw("Transmit succeeded", "resp", resp)
		incRequestStatusMetric(statusSuccess)
		setRequestLatencyMetric(float64(time.Since(start).Milliseconds()))
	}
	return
}

func (w *client) LatestReport(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error) {
	lggr := w.logger.With("req.FeedId", hexutil.Encode(req.FeedId))
	lggr.Debug("LatestReport")
	if err := w.waitForReady(ctx); err != nil {
		return nil, errors.Wrap(err, "LatestReport failed")
	}
	resp, err = w.client.LatestReport(ctx, req)
	if err != nil {
		lggr.Errorw("LatestReport failed", "err", err, "req", req, "resp", resp)
	} else if resp.Error != "" {
		lggr.Errorw("LatestReport failed; mercury server returned error", "err", resp.Error, "req", req, "resp", resp)
	} else {
		lggr.Debugw("LatestReport succeeded", "resp", resp)
	}
	return
}
