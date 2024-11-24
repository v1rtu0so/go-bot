package transactions

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// SignTransaction signs a transaction using the provided private key.
func SignTransaction(tx *solana.Transaction, wallet *solana.Wallet) error {
	err := tx.Sign(wallet.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}
	return nil
}

// SendTransaction sends a signed transaction to the Solana blockchain.
func SendTransaction(ctx context.Context, client *rpc.Client, tx *solana.Transaction) (string, error) {
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

	return sig.String(), nil
}
