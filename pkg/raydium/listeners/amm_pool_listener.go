// pkg/raydium/listeners/amm_pool_listener.go

package listeners

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"

	"corvus_bot/pkg/config"
)

// RawPoolEvent represents the raw data captured from a pool creation event
type RawPoolEvent struct {
	Signature solana.Signature
	Logs      []string
	Slot      uint64
	BlockTime int64
	TxData    *rpc.GetTransactionResult
}

// AMMPoolListener handles WebSocket connections for Raydium AMM pool creation events
type AMMPoolListener struct {
	rpcClient *rpc.Client
	wsClient  *ws.Client
	programID solana.PublicKey
	eventChan chan *RawPoolEvent
}

// NewAMMPoolListener creates a new instance of the AMM pool listener
func NewAMMPoolListener(cfg *config.Config) (*AMMPoolListener, error) {
	programID, err := solana.PublicKeyFromBase58(cfg.RaydiumAMMProgramID)
	if err != nil {
		return nil, fmt.Errorf("invalid Raydium AMM program ID: %w", err)
	}

	wsClient, err := ws.Connect(context.Background(), cfg.WSConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	return &AMMPoolListener{
		rpcClient: rpc.New(cfg.RPCConnection),
		wsClient:  wsClient,
		programID: programID,
		eventChan: make(chan *RawPoolEvent, 100), // Buffer for 100 events
	}, nil
}

// Start begins listening for AMM pool creation events
func (l *AMMPoolListener) Start(ctx context.Context) error {
	log.Printf("Starting AMM pool listener for program: %s", l.programID.String())

	sub, err := l.wsClient.LogsSubscribeMentions(
		l.programID,
		rpc.CommitmentProcessed,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
		l.wsClient.Close()
		close(l.eventChan)
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				got, err := sub.Recv(ctx)
				if err != nil {
					log.Printf("Error receiving message: %v", err)
					continue
				}
				if got != nil && len(got.Value.Logs) > 0 {
					// Look for pool creation events
					if l.isPoolCreationEvent(got.Value.Logs) {
						l.capturePoolEvent(ctx, got)
					}
				}
			}
		}
	}()

	return nil
}

// isPoolCreationEvent checks if logs contain pool initialization
func (l *AMMPoolListener) isPoolCreationEvent(logs []string) bool {
	for _, log := range logs {
		if strings.Contains(log, "initialize2") {
			return true
		}
	}
	return false
}

// capturePoolEvent fetches complete transaction data for a pool creation event
func (l *AMMPoolListener) capturePoolEvent(ctx context.Context, log *ws.LogResult) {
	opts := &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentProcessed,
	}

	tx, err := l.rpcClient.GetTransaction(ctx, log.Value.Signature, opts)
	if err != nil {
		log.Printf("Error fetching transaction %s: %v", log.Value.Signature, err)
		return
	}

	event := &RawPoolEvent{
		Signature: log.Value.Signature,
		Logs:      log.Value.Logs,
		Slot:      tx.Slot,
		TxData:    tx,
	}

	if tx.BlockTime != nil {
		event.BlockTime = tx.BlockTime.UnixSeconds()
	}

	l.eventChan <- event
}

// GetEventChannel returns the channel for receiving pool events
func (l *AMMPoolListener) GetEventChannel() <-chan *RawPoolEvent {
	return l.eventChan
}

// Stop gracefully shuts down the listener
func (l *AMMPoolListener) Stop() {
	if l.wsClient != nil {
		l.wsClient.Close()
	}
}
