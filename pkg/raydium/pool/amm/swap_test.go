package amm

import (
	"context"
	"testing"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
)

func TestSwapTokens(t *testing.T) {
	// Mock configuration with testing enabled
	cfg := &config.Config{
		Testing: true,
	}

	// Mock RPC client
	utils.SetMockRPCClient(&utils.MockRPCClient{
		MockGetLatestBlockhash: func(ctx context.Context, commitment rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error) {
			return &rpc.GetLatestBlockhashResult{
				Value: &rpc.LatestBlockhashResult{
					Blockhash:            solana.MustHashFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"), // Valid base58
					LastValidBlockHeight: 123456,
				},
			}, nil
		},
		MockSendTransaction: func(ctx context.Context, tx *solana.Transaction) (solana.Signature, error) {
			return solana.MustSignatureFromBase58("67ZwYesdsTXe9noVdf8Z1DqMvXPErsWL2y3QiGM3Lg9t1S63bzBTxkkieUWbCg1iCjscy1Lk1TH8uYpMNbcBaHdA"), nil
		},
	})

	// Create a mock wallet
	privateKey, _ := solana.NewRandomPrivateKey()
	wallet := &solana.Wallet{
		PrivateKey: privateKey,
	}

	// Create a mock pool
	pool := RaydiumAmmPool{
		BaseVault:  "8A8R5PA2mNe5dcdbMtGncd2KrCqAEogKcs39XhDhMdSA", // Valid base58
		QuoteVault: "G5oAZj84EXR7TPxY7wAsEoA7R9ZqrStESp3Dz8cHkSo3", // Valid base58
		ProgramID:  "Fg6PaFpoGXkYsidMpWTK6W2BeZ7FEfcYkg9ZpTTDdmz8", // Valid base58
	}

	// Define input and output parameters
	inputAmount := uint64(1000)
	minOutputAmount := uint64(900)

	// Retrieve the correct client based on the mock config
	client := utils.GetRPCClient(cfg)

	// Call SwapTokens using the correct RPC client
	sig, err := SwapTokens(
		context.Background(),
		client, // Pass the RPC client instead of cfg
		wallet,
		pool,
		solana.MustPublicKeyFromBase58(pool.BaseVault),
		solana.MustPublicKeyFromBase58(pool.QuoteVault),
		inputAmount,
		minOutputAmount,
	)

	// Assertions to verify correctness
	assert.NoError(t, err, "expected no error during token swap")
	assert.Equal(t, solana.MustSignatureFromBase58("67ZwYesdsTXe9noVdf8Z1DqMvXPErsWL2y3QiGM3Lg9t1S63bzBTxkkieUWbCg1iCjscy1Lk1TH8uYpMNbcBaHdA"), sig, "unexpected signature returned")
}
