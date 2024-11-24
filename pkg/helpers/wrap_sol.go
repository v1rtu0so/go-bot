package helpers

import (
	"context"
	"corvus_bot/internal/solana"
	"fmt"
	"log"

	"github.com/mr-tron/base58"
)

// WrapSOL creates a raw transaction to wrap SOL into WSOL.
func WrapSOL(ctx context.Context, rpcURL string, payerPrivKey string, amount uint64) (string, string, error) {
	// Decode payer private key
	payerPrivBytes, err := base58.Decode(payerPrivKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode payer private key: %w", err)
	}
	payerPubKey := solana.DerivePublicKey(payerPrivBytes)

	// Generate a new account for WSOL
	wsolPrivBytes := solana.GeneratePrivateKey()
	wsolPubKey := solana.DerivePublicKey(wsolPrivBytes)

	// Fetch the recent blockhash
	blockhash, err := solana.GetRecentBlockhash(ctx, rpcURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch recent blockhash: %w", err)
	}

	// Example rent exemption: adjust dynamically if needed
	rentExemption := uint64(2039280)

	// Create raw instructions
	createAccountIx := solana.CreateAccountInstruction(payerPubKey, wsolPubKey, rentExemption+amount, 165)
	initializeAccountIx := solana.InitializeAccountInstruction(wsolPubKey, payerPubKey)

	// Build transaction
	rawTx := solana.BuildRawTransaction([]solana.Instruction{createAccountIx, initializeAccountIx}, blockhash, payerPubKey)

	// Sign transaction
	signedTx := solana.SignTransaction(rawTx, [][]byte{payerPrivBytes, wsolPrivBytes})

	// Send transaction
	txHash, err := solana.SendRawTransaction(ctx, rpcURL, signedTx)
	if err != nil {
		return "", "", fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("WrapSOL successful: txHash=%s, wsolAccount=%s", txHash, base58.Encode(wsolPubKey))
	return txHash, base58.Encode(wsolPubKey), nil
}
