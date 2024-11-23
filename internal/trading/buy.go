package trading

import (
	"context"
	"fmt"
	"log"
	"math"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/raydium"
	"corvus_bot/pkg/solana"
)

// Buy executes a buy operation
func Buy(cfg *config.Config, tokenMintAddress string, amount, slippage float64) error {
	// Use the RPCConnection from the config
	client := solana.NewClient(cfg.RPCConnection)

	// Use the PrivateKey from the config
	privateKey, err := solana.PrivateKeyFromString(cfg.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	// Fetch Raydium pool information
	pool, err := raydium.FetchPoolInfo(client.Client, tokenMintAddress)
	if err != nil {
		return fmt.Errorf("failed to fetch pool info: %w", err)
	}

	// Calculate token amounts
	swapAmount := uint64(amount * math.Pow10(int(pool.TokenDecimals)))
	minOutputAmount := uint64(float64(swapAmount) * (1 - slippage/100))

	// Create transaction
	tx, err := solana.NewTransaction()
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Add Raydium swap instruction
	err = raydium.AddSwapInstruction(tx, pool, privateKey.PublicKey(), swapAmount, minOutputAmount, true)
	if err != nil {
		return fmt.Errorf("failed to add swap instruction: %w", err)
	}

	// Set recent blockhash
	blockhash, err := client.GetRecentBlockhash(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch recent blockhash: %w", err)
	}
	tx.SetRecentBlockHash(blockhash.Blockhash)

	// Sign the transaction
	err = tx.Sign([]solana.PrivateKey{privateKey})
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send the transaction
	txID, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("Buy transaction sent: %s", txID)
	return nil
}
