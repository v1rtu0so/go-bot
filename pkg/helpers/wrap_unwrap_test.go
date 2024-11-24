package helpers

import (
	"context"
	"fmt"
	"log"

	"github.com/mr-tron/base58"
)

// UnwrapSOL closes the WSOL account, returning SOL back to the payer.
func UnwrapSOL(ctx context.Context, rpcURL string, payerPrivKey string, wsolAccount string) (string, error) {
	// Decode payer private key
	payerPrivBytes, err := base58.Decode(payerPrivKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode payer private key: %w", err)
	}
	payerPubKey := DerivePublicKey(payerPrivBytes)

	// Decode WSOL account public key
	wsolPubKey, err := base58.Decode(wsolAccount)
	if err != nil {
		return "", fmt.Errorf("failed to decode WSOL account public key: %w", err)
	}

	// Fetch the recent blockhash
	blockhash, err := GetRecentBlockhash(ctx, rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch recent blockhash: %w", err)
	}

	// Create raw instruction to close the WSOL account
	closeAccountIx := CloseAccountInstruction(wsolPubKey, payerPubKey, payerPubKey)

	// Construct the raw transaction
	rawTx := BuildRawTransaction([]Instruction{closeAccountIx}, blockhash, payerPubKey)

	// Sign the transaction
	signatures := SignTransaction(rawTx, []PrivateKey{payerPrivBytes})

	// Send the raw transaction
	txHash, err := SendRawTransaction(ctx, rpcURL, signatures)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("UnwrapSOL successful: txHash=%s", txHash)
	return txHash, nil
}
