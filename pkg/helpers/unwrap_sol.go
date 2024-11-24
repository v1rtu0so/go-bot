package helpers

import (
	"context"
	"corvus_bot/internal/solana"
	"corvus_bot/pkg/config"
	"fmt"
	"log"

	"github.com/mr-tron/base58"
)

// UnwrapSOL unwraps WSOL back to SOL.
func UnwrapSOL(ctx context.Context, cfg *config.Config, wsolAccount string) (string, error) {
	// Decode private key
	payerPrivKey, err := base58.Decode(cfg.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode payer private key: %w", err)
	}
	payerPubKey := solana.DerivePublicKey(payerPrivKey)

	// Decode WSOL account
	wsolPubKey, err := base58.Decode(wsolAccount)
	if err != nil {
		return "", fmt.Errorf("failed to decode WSOL account: %w", err)
	}

	// Fetch recent blockhash
	blockhash, err := solana.GetRecentBlockhash(ctx, cfg.RPCConnection)
	if err != nil {
		return "", fmt.Errorf("failed to fetch recent blockhash: %w", err)
	}

	// Convert public keys to [32]byte
	payerPubKeyArray := solana.ToArray32(payerPubKey)
	wsolPubKeyArray := solana.ToArray32(wsolPubKey)

	// Create raw transaction instructions
	closeAccountIx := solana.CloseAccountInstruction(wsolPubKeyArray, payerPubKeyArray, payerPubKeyArray)

	// Build raw transaction
	rawTx := solana.BuildRawTransaction([]solana.Instruction{closeAccountIx}, blockhash, payerPubKeyArray)

	// Sign transaction
	signedTx, err := solana.SignTransaction(&rawTx, [][]byte{payerPrivKey})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	txHash, err := solana.SendRawTransaction(ctx, cfg.RPCConnection, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("UnwrapSOL successful: txHash=%s", txHash)
	return txHash, nil
}
