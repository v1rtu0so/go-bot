// File: /mnt/dev/go-bot/pkg/raydium/listeners/amm_pool_listener.go

package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/raydium/parse"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// AMMPoolListener manages the WebSocket subscription to pool events
type AMMPoolListener struct {
	wsClient  *ws.Client
	eventChan chan *parse.ParsedAMMPool
	config    *config.Config
	parser    *parse.AMMParser
}

// NewAMMPoolListener creates a new listener instance
func NewAMMPoolListener(cfg *config.Config) (*AMMPoolListener, error) {
	wsClient, err := ws.Connect(context.Background(), cfg.WSConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	parser, err := parse.NewAMMParser(cfg.RaydiumAMMProgramID)
	if err != nil {
		wsClient.Close()
		return nil, fmt.Errorf("failed to create parser: %w", err)
	}

	return &AMMPoolListener{
		wsClient:  wsClient,
		eventChan: make(chan *parse.ParsedAMMPool, 100),
		config:    cfg,
		parser:    parser,
	}, nil
}

// Start begins listening for pool events
func (l *AMMPoolListener) Start(ctx context.Context) error {
	programID := solana.MustPublicKeyFromBase58(l.config.RaydiumAMMProgramID)

	sub, err := l.wsClient.LogsSubscribeMentions(
		programID,
		rpc.CommitmentProcessed,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}
	defer sub.Unsubscribe()

	log.Printf("Subscribed to Raydium AMM logs for program: %s", programID)

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Open log file for raw data
	file, err := os.OpenFile("logs/raydium_pool_events.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	for {
		select {
		case <-ctx.Done():
			return nil
		case resp, ok := <-sub.Response():
			if !ok {
				return fmt.Errorf("subscription channel closed")
			}
			if resp.Value.Err != nil {
				log.Printf("Error in log response: %v", resp.Value.Err)
				continue
			}

			// Save raw log data
			logEntry := struct {
				Timestamp time.Time
				Data      *ws.LogResult
			}{
				Timestamp: time.Now().UTC(),
				Data:      resp.Value,
			}
			if err := encoder.Encode(logEntry); err != nil {
				log.Printf("Error saving raw log: %v", err)
			}

			// Parse pool initialization
			pool, err := l.parser.ParsePoolInit(resp.Value)
			if err != nil {
				// Only log if it's not a "no initialize2 instruction found" error
				if err.Error() != "no initialize2 instruction found in logs" {
					log.Printf("Error parsing pool initialization: %v", err)
				}
				continue
			}

			// Send parsed pool data to channel
			select {
			case l.eventChan <- pool:
				log.Printf("New pool detected - ID: %s, Base: %s, Quote: %s",
					pool.ID, pool.BaseMint, pool.QuoteMint)
			default:
				log.Printf("Warning: Event channel full, dropping pool event for ID: %s", pool.ID)
			}
		}
	}
}

// GetEventChannel returns the channel for receiving parsed pool events
func (l *AMMPoolListener) GetEventChannel() chan *parse.ParsedAMMPool {
	return l.eventChan
}

// Close closes the WebSocket connection
func (l *AMMPoolListener) Close() {
	l.wsClient.Close()
}

// SaveRawLogs saves the raw log data to a file for analysis
func (l *AMMPoolListener) SaveRawLogs(logValue *ws.LogResult) error {
	logData := struct {
		Timestamp time.Time `json:"timestamp"`
		Signature string    `json:"signature"`
		Slot      uint64    `json:"slot"`
		Logs      []string  `json:"logs"`
	}{
		Timestamp: time.Now().UTC(),
		Signature: logValue.Signature.String(),
		Slot:      logValue.Context.Slot,
		Logs:      logValue.Logs,
	}

	file, err := os.OpenFile("logs/raydium_raw_logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open raw logs file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(logData); err != nil {
		return fmt.Errorf("failed to write log data: %w", err)
	}

	return nil
}
