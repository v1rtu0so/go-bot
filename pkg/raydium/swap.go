package raydium

import (
	"context"
	"fmt"
	"log"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/helpers"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// SwapTokens performs a token swap on Raydium.
func SwapTokens(
	ctx context.Context,
	cfg *config.Config,
	client *rpc.Client,
	wallet *solana.Wallet,
	inputMint solana.PublicKey,
	outputMint solana.PublicKey,
	amountIn uint64,
	slippage float64,
) error {
	// Fetch the target pool
	pool, err := findPool(inputMint, outputMint)
	if err != nil {
		return fmt.Errorf("failed to find pool: %w", err)
	}

	// Build the swap transaction
	swapTx, err := buildSwapTransaction(ctx, cfg, wallet, pool, amountIn, slippage)
	if err != nil {
		return fmt.Errorf("failed to build swap transaction: %w", err)
	}

	// Submit the transaction
	txSignature, err := submitTransaction(ctx, client, swapTx, wallet)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	log.Printf("Transaction submitted successfully. Signature: %s", txSignature)
	return nil
}

// findPool identifies the appropriate Raydium pool for the given token pair.
func findPool(inputMint, outputMint solana.PublicKey) (*RaydiumPool, error) {
	// Fetch and parse pool data from a stored JSON file
	poolData := helpers.LoadPoolData("data/testdata/amm_pools.json") // Example path
	for _, pool := range poolData {
		if pool.InputMint == inputMint && pool.OutputMint == outputMint {
			return &pool, nil
		}
	}
	return nil, fmt.Errorf("no suitable pool found for tokens: %s -> %s", inputMint, outputMint)
}

// buildSwapTransaction constructs a transaction for a Raydium swap.
func buildSwapTransaction(
	ctx context.Context,
	cfg *config.Config,
	wallet *solana.Wallet,
	pool *RaydiumPool,
	amountIn uint64,
	slippage float64,
) (*solana.Transaction, error) {
	// Calculate minimum amount out considering slippage
	minAmountOut := calculateMinAmountOut(pool, amountIn, slippage)

	// Create instructions for the swap
	instructions, err := generateSwapInstructions(cfg, pool, amountIn, minAmountOut, wallet.PublicKey())
	if err != nil {
		return nil, fmt.Errorf("failed to generate swap instructions: %w", err)
	}

	// Build transaction
	tx := solana.NewTransaction(
		[]solana.Instruction{instructions},
		wallet.PublicKey(),
	)

	// Set blockhash and recent commitment level
	blockhash, err := rpcClient.GetLatestBlockhash(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blockhash: %w", err)
	}
	tx.SetBlockhash(blockhash.Blockhash)

	return tx, nil
}

// submitTransaction signs and submits the transaction to the Solana blockchain.
func submitTransaction(
	ctx context.Context,
	client *rpc.Client,
	tx *solana.Transaction,
	wallet *solana.Wallet,
) (string, error) {
	// Sign the transaction
	err := tx.Sign(wallet.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Serialize the transaction
	rawTx, err := tx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// Submit the transaction
	sig, err := client.SendTransaction(ctx, rawTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	// Confirm the transaction
	err = helpers.ConfirmTransaction(ctx, client, sig)
	if err != nil {
		return "", fmt.Errorf("transaction confirmation failed: %w", err)
	}

	return sig, nil
}

// calculateMinAmountOut computes the minimum output amount based on slippage.
func calculateMinAmountOut(pool *RaydiumPool, amountIn uint64, slippage float64) uint64 {
	amountOut := pool.CalculateOutput(amountIn)
	minAmountOut := uint64(float64(amountOut) * (1 - slippage))
	return minAmountOut
}

// generateSwapInstructions creates swap instructions for the transaction.
func generateSwapInstructions(
	cfg *config.Config,
	pool *RaydiumPool,
	amountIn, minAmountOut uint64,
	owner solana.PublicKey,
) (solana.Instruction, error) {
	// Placeholder for generating instructions
	// Ensure this matches Raydium's program specifications
	// Example structure: instruction := NewSwapInstruction(pool, amountIn, minAmountOut, owner)
	return nil, nil // TODO: Implement instruction generation
}
