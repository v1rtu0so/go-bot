package amm

import (
	"context"
	"encoding/binary"
	"fmt"

	"corvus_bot/pkg/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// SwapTokens performs a token swap within the Raydium AMM pool.
func SwapTokens(
	ctx context.Context,
	client utils.RPCClientInterface, // Use RPCClientInterface instead of *config.Config
	wallet *solana.Wallet,
	pool RaydiumAmmPool,
	inputMint, outputMint solana.PublicKey,
	inputAmount, minOutputAmount uint64,
) (solana.Signature, error) {
	// Convert string-based public keys in the pool to solana.PublicKey
	baseVault, err := solana.PublicKeyFromBase58(pool.BaseVault)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("invalid BaseVault public key: %w", err)
	}

	quoteVault, err := solana.PublicKeyFromBase58(pool.QuoteVault)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("invalid QuoteVault public key: %w", err)
	}

	programID, err := solana.PublicKeyFromBase58(pool.ProgramID)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("invalid ProgramID public key: %w", err)
	}

	// Fetch the latest blockhash
	blockhashResult, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to fetch latest blockhash: %w", err)
	}

	// Construct the data payload for the Raydium swap instruction
	data := make([]byte, 16)
	binary.LittleEndian.PutUint64(data[0:], inputAmount)     // Input amount
	binary.LittleEndian.PutUint64(data[8:], minOutputAmount) // Minimum output amount

	// Construct the instruction
	swapInstruction := solana.NewInstruction(
		programID,
		solana.AccountMetaSlice{
			solana.NewAccountMeta(baseVault, true, false),          // Base token vault
			solana.NewAccountMeta(quoteVault, true, false),         // Quote token vault
			solana.NewAccountMeta(wallet.PublicKey(), false, true), // Payer wallet
		},
		data,
	)

	// Create the transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{swapInstruction},
		blockhashResult.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Sign the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(wallet.PublicKey()) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send the transaction
	signature, err := client.SendTransaction(ctx, tx)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signature, nil
}
