package solana

import (
	"testing"

	"corvus_bot/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestSolanaClient(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("../../pkg/config/config.yaml")
	assert.NoError(t, err, "Failed to load configuration")

	// Initialize the Solana client
	client, err := NewSolanaClient(cfg.RPCConnection, cfg.WSConnection)
	assert.NoError(t, err, "Failed to initialize Solana client")

	// Test GetAccountBalance
	balance, err := client.GetAccountBalance(cfg.WSOLAddress) // Using WSOL address for test
	assert.NoError(t, err, "Failed to fetch account balance")
	assert.GreaterOrEqual(t, balance, uint64(0), "Balance should be zero or greater")

	// Example test for SendTransaction (requires controlled test environment)
	/*
		recentBlockhash := solana.Hash{} // Add a valid recent blockhash here
		tx, err := CreateTransaction(recentBlockhash, nil) // Add valid instructions
		assert.NoError(t, err, "Failed to create transaction")

		sig, err := client.SendTransaction(tx, cfg.PrivateKey)
		assert.NoError(t, err, "Failed to send transaction")
		assert.NotEmpty(t, sig, "Transaction signature should not be empty")
	*/
}
