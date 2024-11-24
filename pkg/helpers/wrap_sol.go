package helpers

import (
	"context"
	"corvus_bot/internal/solana"
	"corvus_bot/pkg/config"
	"fmt"
	"log"

	"github.com/mr-tron/base58"
)

// WrapSOL wraps SOL into WSOL.
func WrapSOL(ctx context.Context, cfg *config.Config, amount uint64) (string, string, error) {
	// Decode private key
	payerPrivKey, err := base58.Decode(cfg.PrivateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode payer private key: %w", err)
	}
	payerPubKey := solana.DerivePublicKey(payerPrivKey)

	// Generate a new keypair for WSOL
	wsolPrivKey := solana.GeneratePrivateKey()
	wsolPubKey := solana.DerivePublicKey(wsolPrivKey)

	// Fetch recent blockhash
	blockhash, err := solana.GetRecentBlockhash(ctx, cfg.RPCConnection)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch recent blockhash: %w", err)
	}

	// Convert public keys to [32]byte
	payerPubKeyArray := solana.ToArray32(payerPubKey)
	wsolPubKeyArray := solana.ToArray32(wsolPubKey)

	// Rent exemption value
	rentExemption := uint64(2039280) // Adjust based on the latest cluster requirements

	// Create raw transaction instructions
	createAccountIx := solana.CreateAccountInstruction(payerPubKeyArray, wsolPubKeyArray, rentExemption+amount, 165)
	initializeAccountIx := solana.InitializeAccountInstruction(wsolPubKeyArray, solana.WSOLMintPubKey(), payerPubKeyArray)

	// Build raw transaction
	rawTx := solana.BuildRawTransaction([]solana.Instruction{createAccountIx, initializeAccountIx}, blockhash, payerPubKeyArray)

	// Sign transaction
	signedTx, err := solana.SignTransaction(&rawTx, [][]byte{payerPrivKey, wsolPrivKey})
	if err != nil {
		return "", "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	txHash, err := solana.SendRawTransaction(ctx, cfg.RPCConnection, signedTx)
	if err != nil {
		return "", "", fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("WrapSOL successful: txHash=%s, wsolAccount=%s", txHash, base58.Encode(wsolPubKey))
	return txHash, base58.Encode(wsolPubKey), nil
}
