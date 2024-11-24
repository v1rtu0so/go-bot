package helpers

import (
	"context"
	"fmt"
	"log"

	"github.com/mr-tron/base58"
)

func UnwrapSOL(ctx context.Context, rpcURL string, payerPrivKey string, wsolAccount string) (string, error) {
	// Convert private key and derive public key
	payerPrivBytes, err := base58.Decode(payerPrivKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %w", err)
	}
	payerPubKey := DerivePublicKey(payerPrivBytes)

	// Convert WSOL account to bytes
	wsolPubKey, err := base58.Decode(wsolAccount)
	if err != nil {
		return "", fmt.Errorf("failed to decode WSOL account: %w", err)
	}

	// Fetch recent blockhash
	blockhash, err := GetRecentBlockhash(ctx, rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch recent blockhash: %w", err)
	}

	// Construct raw instruction to close WSOL account
	closeAccountIx := CloseAccountInstruction(wsolPubKey, payerPubKey, payerPubKey)

	// Create raw transaction
	rawTx := BuildRawTransaction([]Instruction{closeAccountIx}, blockhash, payerPubKey)

	// Sign transaction
	signatures := SignTransaction(rawTx, []PrivateKey{payerPrivBytes})

	// Send transaction
	txHash, err := SendRawTransaction(ctx, rpcURL, signatures)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("UnwrapSOL transaction hash: %s", txHash)
	return txHash, nil
}
