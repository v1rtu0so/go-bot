package transactions

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// CreateTransaction creates a Solana transaction with the provided instructions.
func CreateTransaction(
	ctx context.Context,
	client *rpc.Client,
	instructions []solana.Instruction,
	payer solana.PublicKey,
) (*solana.Transaction, error) {
	// Get the latest blockhash
	blockhash, err := client.GetLatestBlockhash(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest blockhash: %w", err)
	}

	// Create a new transaction
	tx, err := solana.NewTransaction(
		instructions,
		blockhash.Blockhash,
		payer,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}
