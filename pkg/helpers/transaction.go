package helpers

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// BuildTransaction creates a new Solana transaction using the provided instructions.
func BuildTransaction(ctx context.Context, client *rpc.Client, payer solana.PrivateKey, instructions []solana.Instruction) (*solana.Transaction, error) {
	// Fetch the latest blockhash
	blockhashResp, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest blockhash: %w", err)
	}

	// Create the transaction
	tx, err := solana.NewTransaction(
		instructions,
		blockhashResp.Value.Blockhash, // Use latest blockhash
		solana.TransactionPayer(payer.PublicKey()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %w", err)
	}

	// Sign the transaction with the payer
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(payer.PublicKey()) {
				return &payer
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return tx, nil
}
