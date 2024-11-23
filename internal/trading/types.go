package trading

import (
	"context"
	"fmt"

	"corvus_bot/pkg/raydium"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// createSwapTransaction creates a swap transaction (buy or sell) on Raydium.
func createSwapTransaction(
	client *rpc.Client,
	pool *raydium.PoolInfo,
	privateKey solana.PrivateKey,
	swapAmount, minOutputAmount uint64,
	isBuy bool,
) (*solana.Transaction, error) {
	tx, err := solana.NewTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Add the swap instruction
	err = raydium.AddSwapInstruction(tx, pool, privateKey.PublicKey(), swapAmount, minOutputAmount, isBuy)
	if err != nil {
		return nil, fmt.Errorf("failed to add swap instruction: %w", err)
	}

	// Set the recent blockhash
	recentBlockhash, err := client.GetRecentBlockhash(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get recent blockhash: %w", err)
	}
	tx.SetRecentBlockHash(recentBlockhash.Blockhash)

	// Sign the transaction
	err = tx.Sign([]solana.PrivateKey{privateKey})
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return tx, nil
}
