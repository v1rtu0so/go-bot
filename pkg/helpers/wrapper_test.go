package helpers

import (
	"context"
	"log"
	"testing"

	"corvus_bot/pkg/config"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
)

// TestWrapperMainnetWithConfig tests wrapping and unwrapping SOL using configuration.
func TestWrapperMainnetWithConfig(t *testing.T) {
	// Load the configuration
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Decode the private key from Base58
	privateKey, err := solana.PrivateKeyFromBase58(cfg.PrivateKey)
	if err != nil {
		t.Fatalf("Failed to decode Base58 private key: %v", err)
	}

	// Extract the public key
	payer := privateKey.PublicKey()

	// Initialize the RPC client using the configuration
	client := rpc.New(cfg.RPCConnection)

	// Debugging: Log wallet public key
	log.Printf("Test wallet public key: %s", payer)

	// Fetch balance and ensure it's sufficient
	balance, err := client.GetBalance(context.Background(), payer, rpc.CommitmentProcessed)
	if err != nil {
		t.Fatalf("Failed to fetch wallet balance: %v", err)
	}
	log.Printf("Initial wallet balance: %.9f SOL", float64(balance.Value)/float64(solana.LAMPORTS_PER_SOL))

	// Ensure the wallet has enough SOL for testing
	minimumBalance := uint64(0.1 * float64(solana.LAMPORTS_PER_SOL))
	if balance.Value < minimumBalance {
		t.Fatalf("Insufficient SOL balance in the wallet: %.9f SOL", float64(balance.Value)/float64(solana.LAMPORTS_PER_SOL))
	}

	// Test wrapping SOL into WSOL
	amountToWrap := uint64(0.05 * float64(solana.LAMPORTS_PER_SOL)) // Wrapping 0.05 SOL
	err = WrapSOL(context.Background(), client, privateKey, amountToWrap, cfg)
	assert.NoError(t, err, "Failed to wrap SOL into WSOL")

	// Test unwrapping WSOL back into SOL
	err = UnwrapSOL(context.Background(), client, privateKey, cfg)
	assert.NoError(t, err, "Failed to unwrap WSOL back into SOL")
}
