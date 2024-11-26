package clmm_test

import (
	"context"
	"testing"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/raydium/pool/clmm"
	"corvus_bot/pkg/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
)

func TestSwapTokens(t *testing.T) {
	// Set a mock RPC client globally
	utils.SetMockRPCClient(&utils.MockRPCClient{
		MockGetLatestBlockhash: func(ctx context.Context, commitment rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error) {
			return &rpc.GetLatestBlockhashResult{
				Value: &rpc.LatestBlockhashResult{
					Blockhash:            solana.MustHashFromBase58("5NKsd8FNdL3jkGjBmHNL6QQFrdkW7ZytnkZ1WFAj3YP9"), // Valid Base58
					LastValidBlockHeight: 123456,
				},
			}, nil
		},
		MockSendTransaction: func(ctx context.Context, tx *solana.Transaction) (solana.Signature, error) {
			return solana.MustSignatureFromBase58("1111111111111111111111111115KfvzcywpQZLKjb8QUGF2vhHXM8FhxzBc3Umi87SsynoE4DntWu"), nil
		},
	})

	// Initialize test data
	ctx := context.Background()
	privateKey, _ := solana.NewRandomPrivateKey()
	wallet := &solana.Wallet{PrivateKey: privateKey}
	pool := &clmm.RaydiumClmmPool{
		VaultA:         "5ZrXUACAbF3bsxRmwAVmsA8AadVZEs5HcVSqrEL9ukXR", // Valid Base58
		VaultB:         "3eA1N7VTJcv2k8NhEA6LjcQGksLRUBhHEnRYBL4U1waK", // Valid Base58
		MintProgramIDA: "4NDj5HjVUN9f8ZMWCcUJ6TAqZTayDQgBV9ZcyYvFF1RU", // Valid Base58
		MintProgramIDB: "GDdR1ZhWQUwUSL69TsvZjWg7FgL1hsA2azXwvNbx3pE8", // Valid Base58
	}
	cfg := &config.Config{
		Testing:              true,
		RaydiumCLMMProgramID: "4NDj5HjVUN9f8ZMWCcUJ6TAqZTayDQgBV9ZcyYvFF1RU", // Valid Base58
		RPCConnection:        "https://api.mainnet-beta.solana.com",          // Valid RPC URL
	}
	amountIn := uint64(1000)
	minAmountOut := uint64(900)

	// Call SwapTokens
	signature, err := clmm.SwapTokens(ctx, utils.GetRPCClient(cfg), wallet, pool, amountIn, minAmountOut, cfg)

	// Assertions
	assert.NoError(t, err, "expected no error during token swap")
	assert.Equal(t, solana.MustSignatureFromBase58("1111111111111111111111111115KfvzcywpQZLKjb8QUGF2vhHXM8FhxzBc3Umi87SsynoE4DntWu"), signature, "unexpected signature returned")
}
