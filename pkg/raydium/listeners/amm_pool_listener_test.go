// pkg/raydium/listeners/amm_pool_listener_test.go

package listeners

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"corvus_bot/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAMMPoolListener(t *testing.T) {
	// Create test directory
	testDir := t.TempDir()
	testConfigPath := filepath.Join(testDir, "test_config.yaml")

	// Create test config file
	configContent := []byte(`
rpc_connection: "http://localhost:8899"
ws_connection: "ws://localhost:8900"
raydium_amm_program_id: "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
testing: true
`)
	err := os.WriteFile(testConfigPath, configContent, 0644)
	require.NoError(t, err, "Failed to create test config file")

	// Load config
	cfg, err := config.LoadConfig(testConfigPath)
	require.NoError(t, err, "Failed to load config")

	// Create listener
	listener, err := NewAMMPoolListener(cfg)
	require.NoError(t, err, "Failed to create listener")
	defer listener.Close()

	// Start listening in background
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	go func() {
		err := listener.Start(ctx)
		if err != nil && err != context.DeadlineExceeded {
			t.Errorf("Listener error: %v", err)
		}
	}()

	// Test parsing functions
	t.Run("ParseInitializeParams", func(t *testing.T) {
		testLog := `Program log: initialize2: InitializeInstruction2 { nonce: 254, open_time: 1732745528, init_pc_amount: 648251900000, init_coin_amount: 206900000000000000 }`
		params, err := parseInitializeParams(testLog)
		require.NoError(t, err)
		require.NotNil(t, params)
		assert.Equal(t, uint8(254), params.Nonce)
		assert.Equal(t, int64(1732745528), params.OpenTime)
		assert.Equal(t, uint64(648251900000), params.InitPcAmount)
		assert.Equal(t, uint64(206900000000000000), params.InitCoinAmount)
	})

	// Test account extraction
	t.Run("ExtractAccount", func(t *testing.T) {
		testLog := "Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]"
		account, err := extractAccount(testLog, "TokenProgram")
		require.NoError(t, err)
		assert.Equal(t, "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA", account.Address)
		assert.Equal(t, RoleProgram, account.Role)
	})

	// Test program ID extraction
	t.Run("ExtractProgramID", func(t *testing.T) {
		testLog := "Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]"
		programID := extractProgramID(testLog)
		assert.Equal(t, "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA", programID)
	})

	// Test complete pool initialization parsing
	t.Run("ParsePoolInitialization", func(t *testing.T) {
		sampleLogs := []string{
			`Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]`,
			`Program log: initialize2: InitializeInstruction2 { nonce: 254, open_time: 1732745528, init_pc_amount: 648251900000, init_coin_amount: 206900000000000000 }`,
			`Program 11111111111111111111111111111111 invoke [2]`,
			`Program log: Initialize the associated token account`,
		}

		parsed, err := parsePoolAccounts(sampleLogs)
		require.NoError(t, err)
		require.NotNil(t, parsed)
		assert.Equal(t, "success", parsed.Status)
	})

	// Test human readable conversion
	t.Run("HumanReadableConversion", func(t *testing.T) {
		parsed := &ParsedPoolInitialization{
			Signature: "test_signature",
			Timestamp: time.Now().UTC(),
			Accounts: PoolAccounts{
				CoinMint: AccountMeta{
					Address: "coin_mint_address",
					TokenInfo: &TokenInfo{
						Symbol: "SOL",
					},
				},
				PCMint: AccountMeta{
					Address: "pc_mint_address",
					TokenInfo: &TokenInfo{
						Symbol: "USDC",
					},
				},
			},
			TokenDetails: struct {
				CoinMintDecimals uint8   `json:"coin_mint_decimals"`
				PCMintDecimals   uint8   `json:"pc_mint_decimals"`
				LPMintDecimals   uint8   `json:"lp_mint_decimals"`
				InitialPrice     float64 `json:"initial_price"`
				CoinAmount       float64 `json:"coin_amount"`
				PCAmount         float64 `json:"pc_amount"`
			}{
				InitialPrice: 100.0,
				CoinAmount:   10.0,
				PCAmount:     1000.0,
			},
		}

		readable := parsed.ToHumanReadable()
		assert.Equal(t, "SOL", readable.TokenA)
		assert.Equal(t, "USDC", readable.TokenB)
		assert.Equal(t, 10.0, readable.InitialLiquidity.TokenAAmount)
		assert.Equal(t, 1000.0, readable.InitialLiquidity.TokenBAmount)
	})
}

func TestParseInitializeParams(t *testing.T) {
	testCases := []struct {
		name     string
		log      string
		expected *InitializeParams
		hasError bool
	}{
		{
			name: "Valid initialize2 log",
			log:  `Program log: initialize2: InitializeInstruction2 { nonce: 254, open_time: 1732745528, init_pc_amount: 648251900000, init_coin_amount: 206900000000000000 }`,
			expected: &InitializeParams{
				Nonce:          254,
				OpenTime:       1732745528,
				InitPcAmount:   648251900000,
				InitCoinAmount: 206900000000000000,
			},
			hasError: false,
		},
		{
			name:     "Non-initialize2 log",
			log:      "Program log: Some other log",
			expected: nil,
			hasError: false,
		},
		{
			name:     "Invalid log format",
			log:      "Program log: initialize2: Invalid format",
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params, err := parseInitializeParams(tc.log)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, params)
			}
		})
	}
}

func TestLiveAMMPoolListener(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping live test in short mode")
	}

	// Use a real config with actual RPC endpoints
	cfg := &config.Config{
		RPCConnection:       "https://api.mainnet-beta.solana.com",
		WSConnection:        "wss://api.mainnet-beta.solana.com",
		RaydiumAMMProgramID: "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8",
	}

	listener, err := NewAMMPoolListener(cfg)
	require.NoError(t, err, "Failed to create listener")
	defer listener.Close()

	// Start listening in background
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	eventsChan := listener.GetEventChannel()
	eventCount := 0

	go func() {
		err := listener.Start(ctx)
		if err != nil && err != context.DeadlineExceeded {
			t.Errorf("Listener error: %v", err)
		}
	}()

	// Listen for events
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			t.Logf("Test timeout reached. Processed %d events", eventCount)
			if eventCount == 0 {
				t.Error("No events received during test period")
			}
			return
		case event := <-eventsChan:
			t.Logf("Received event: %s", event.Signature)
			eventCount++

			// Verify event structure
			assert.NotEmpty(t, event.Signature)
			assert.NotZero(t, event.Slot)
			assert.NotEmpty(t, event.Logs)

			// Verify the output file after receiving event
			if eventCount >= 1 {
				content, err := os.ReadFile("raydium_filtered_logs.json")
				require.NoError(t, err, "Failed to read output file")
				assert.NotEmpty(t, content, "Output file should not be empty")

				// Try to parse the content to verify JSON structure
				var savedEvent RawPoolEvent
				err = json.Unmarshal(content, &savedEvent)
				require.NoError(t, err, "Failed to parse saved event")
				assert.Equal(t, event.Signature, savedEvent.Signature)

				return
			}
		case <-ctx.Done():
			t.Logf("Context cancelled. Processed %d events", eventCount)
			return
		}
	}
}
