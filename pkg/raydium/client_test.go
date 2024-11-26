package raydium

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

const (
	testRPCEndpoint = "http://127.0.0.1:8899"
	testAMMProgram  = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
	testCLMMProgram = "CLMMv1A5qmqh3FQKL7e1ppT88uZqGpuZBtjPxkSwgmt2"
	testTokenAddr   = "2qEHjDLDLbuBgRYvsxhc5D6uDWAivNFZGan56P1tpump" // PNUT token
)

func TestFetchPoolData(t *testing.T) {
	// Create data/testdata directory if it doesn't exist
	testDataDir := filepath.Join(".", "..", "..", "data", "testdata")
	err := os.MkdirAll(testDataDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test data directory: %v", err)
	}

	// Initialize file paths
	ammDataPath := filepath.Join(testDataDir, "amm_pools.json")
	clmmDataPath := filepath.Join(testDataDir, "clmm_pools.json")

	// Create or reset the JSON files
	createEmptyJSONFile(t, ammDataPath)
	createEmptyJSONFile(t, clmmDataPath)

	// Initialize the Raydium client
	client := NewRaydiumClient(
		testRPCEndpoint,
		testAMMProgram,
		testCLMMProgram,
		ammDataPath,
		clmmDataPath,
	)

	// Test both pool types
	poolTypes := []string{"AMM", "CLMM"}
	for _, poolType := range poolTypes {
		t.Run(fmt.Sprintf("FetchPoolData_%s", poolType), func(t *testing.T) {
			ctx := context.Background()

			// First attempt - should fetch from network
			t.Logf("Attempting to fetch %s pool data for token: %s", poolType, testTokenAddr)
			pool, err := client.FetchPoolData(ctx, poolType, testTokenAddr)

			if err != nil {
				t.Logf("Error fetching pool data: %v", err)
				return
			}

			assert.NotNil(t, pool, "Pool data should not be nil")
			t.Logf("Successfully fetched %s pool data from network", poolType)

			// Second attempt - should fetch from storage
			t.Logf("Attempting to fetch %s pool data from storage", poolType)
			poolFromStorage, err := client.FetchPoolData(ctx, poolType, testTokenAddr)

			if err != nil {
				t.Logf("Error fetching pool data from storage: %v", err)
				return
			}

			assert.NotNil(t, poolFromStorage, "Pool data from storage should not be nil")
			assert.Equal(t, pool, poolFromStorage, "Pool data from storage should match network data")
			t.Logf("Successfully retrieved %s pool data from storage", poolType)
		})
	}
}

func TestFetchPoolDataInvalidType(t *testing.T) {
	client := NewRaydiumClient(
		testRPCEndpoint,
		testAMMProgram,
		testCLMMProgram,
		filepath.Join(".", "..", "..", "data", "testdata", "amm_pools.json"),
		filepath.Join(".", "..", "..", "data", "testdata", "clmm_pools.json"),
	)

	_, err := client.FetchPoolData(context.Background(), "INVALID", testTokenAddr)
	assert.Error(t, err, "Should return error for invalid pool type")
	assert.Contains(t, err.Error(), "unsupported pool type")
}

func TestPerformSwap(t *testing.T) {
	ctx := context.Background()

	// Create a test wallet
	privateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		t.Fatalf("Failed to create test wallet: %v", err)
	}
	wallet := &solana.Wallet{PrivateKey: privateKey}

	client := NewRaydiumClient(
		testRPCEndpoint,
		testAMMProgram,
		testCLMMProgram,
		filepath.Join(".", "..", "..", "data", "testdata", "amm_pools.json"),
		filepath.Join(".", "..", "..", "data", "testdata", "clmm_pools.json"),
	)

	// First fetch pool data
	poolData, err := client.FetchPoolData(ctx, "AMM", testTokenAddr)
	if err != nil {
		t.Fatalf("Failed to fetch pool data: %v", err)
	}

	// Test swap with minimal amounts
	amountIn := uint64(1000000)    // 0.001 SOL
	minAmountOut := uint64(900000) // 0.0009 SOL (10% slippage)

	_, err = client.PerformSwap(ctx, wallet, "AMM", poolData, amountIn, minAmountOut)
	if err != nil {
		t.Logf("Swap failed as expected with test wallet: %v", err)
	}
}

// Helper function to create or reset a JSON file with empty array
func createEmptyJSONFile(t *testing.T, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	_, err = file.Write([]byte("[]"))
	if err != nil {
		t.Fatalf("Failed to write to file %s: %v", filePath, err)
	}
}
