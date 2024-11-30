//mnt/dev/go-bot/pkg/raydium/parse/amm_mint_test.go

package parse

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testProgramID = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"

func TestAMMParser(t *testing.T) {
	// Create parser instance
	parser, err := NewAMMParser(testProgramID)
	require.NoError(t, err)

	// Load test data
	logData, err := loadTestLogData(t)
	require.NoError(t, err)

	// Parse pool initialization
	pool, err := parser.ParsePoolInit(logData)
	require.NoError(t, err)

	// Verify parsed data
	assert.NotEmpty(t, pool.ID)
	assert.Equal(t, testProgramID, pool.ProgramID)
	assert.NotEmpty(t, pool.BaseMint)
	assert.NotEmpty(t, pool.QuoteMint)
	assert.NotEmpty(t, pool.BaseVault)
	assert.NotEmpty(t, pool.QuoteVault)

	// Verify amounts
	assert.Greater(t, pool.InitialBase, uint64(0))
	assert.Greater(t, pool.InitialQuote, uint64(0))
}

func TestParseInvalidLog(t *testing.T) {
	parser, err := NewAMMParser(testProgramID)
	require.NoError(t, err)

	// Test with invalid log data
	invalidLog := &ws.LogResult{
		Signature: solana.Signature{},
		Logs:      []string{"some random log"},
	}

	_, err = parser.ParsePoolInit(invalidLog)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no initialize2 instruction found")
}

func TestValidatePoolData(t *testing.T) {
	parser, err := NewAMMParser(testProgramID)
	require.NoError(t, err)

	tests := []struct {
		name    string
		pool    *ParsedAMMPool
		wantErr bool
	}{
		{
			name: "identical mints",
			pool: &ParsedAMMPool{
				BaseMint:  "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
				QuoteMint: "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
			},
			wantErr: true,
		},
		{
			name: "identical vaults",
			pool: &ParsedAMMPool{
				BaseMint:   "So11111111111111111111111111111111111111112",
				QuoteMint:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				BaseVault:  "3uQytDKNd5H6XK8FhTei4wrkxZzkcnecMpVxkJzEGNnG",
				QuoteVault: "3uQytDKNd5H6XK8FhTei4wrkxZzkcnecMpVxkJzEGNnG",
			},
			wantErr: true,
		},
		{
			name: "invalid public key",
			pool: &ParsedAMMPool{
				BaseMint: "invalid-key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.validatePoolData(tt.pool)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// loadTestLogData loads test data from file
func loadTestLogData(t *testing.T) (*ws.LogResult, error) {
	// Read test data file
	data, err := os.ReadFile("testdata/raydium_init_log.json")
	require.NoError(t, err)

	var logResult ws.LogResult
	err = json.Unmarshal(data, &logResult)
	require.NoError(t, err)

	return &logResult, nil
}
