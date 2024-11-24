package raydium

import (
	"fmt"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/raydium/pool/amm" // Correct reference

	"github.com/gagliardetto/solana-go"
)

// GenerateSwapInstructions generates the instructions for a token swap on Raydium.
func GenerateSwapInstructions(
	pool *amm.RaydiumAmmPool, // Use the existing RaydiumAmmPool type
	amountIn uint64,
	minAmountOut uint64,
	owner solana.PublicKey,
	cfg *config.Config,
) ([]solana.Instruction, error) {
	// Validate the input pool data
	if pool.ID == "" || pool.BaseMint == "" || pool.QuoteMint == "" {
		return nil, fmt.Errorf("invalid pool data")
	}

	// Construct accounts involved in the swap
	sourceTokenAccount := solana.MustPublicKeyFromBase58(pool.BaseMint)
	destTokenAccount := solana.MustPublicKeyFromBase58(pool.QuoteMint)
	poolAccount := solana.MustPublicKeyFromBase58(pool.ID)
	programID := solana.MustPublicKeyFromBase58(cfg.RaydiumAMMProgramID)

	// Construct the swap instruction for Raydium
	instruction := solana.NewInstruction(
		programID,
		solana.AccountMetaSlice{
			{PublicKey: sourceTokenAccount, IsWritable: true, IsSigner: false},
			{PublicKey: destTokenAccount, IsWritable: true, IsSigner: false},
			{PublicKey: poolAccount, IsWritable: false, IsSigner: false},
			{PublicKey: owner, IsWritable: false, IsSigner: true},
		},
		[]byte{}, // Replace with the actual payload for Raydium swap
	)

	return []solana.Instruction{instruction}, nil
}

// FetchPoolAndGenerateInstructions fetches the pool and generates swap instructions.
func FetchPoolAndGenerateInstructions(
	tokenAddress string,
	filePath string,
	amountIn uint64,
	minAmountOut uint64,
	owner solana.PublicKey,
	cfg *config.Config,
) ([]solana.Instruction, error) {
	// Fetch the pool from the JSON file
	pool, err := amm.FetchAmmPoolFromJSON(tokenAddress, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch AMM pool: %w", err)
	}

	// Generate swap instructions using the fetched pool
	return GenerateSwapInstructions(pool, amountIn, minAmountOut, owner, cfg)
}
