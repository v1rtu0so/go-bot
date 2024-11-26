package clmm

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// SwapTokens performs a token swap within the Raydium CLMM pool.
func SwapTokens(
	ctx context.Context,
	client utils.RPCClientInterface, // Use the interface for RPC client
	wallet *solana.Wallet,
	pool *RaydiumClmmPool,
	amountIn uint64,
	minAmountOut uint64,
	cfg *config.Config,
) (solana.Signature, error) {
	if pool == nil {
		return solana.Signature{}, fmt.Errorf("pool data is nil")
	}

	// Validate pool data
	if pool.VaultA == "" || pool.VaultB == "" || pool.MintProgramIDA == "" || pool.MintProgramIDB == "" {
		return solana.Signature{}, fmt.Errorf("invalid CLMM pool data")
	}

	// Parse account public keys
	vaultA := solana.MustPublicKeyFromBase58(pool.VaultA)
	vaultB := solana.MustPublicKeyFromBase58(pool.VaultB)
	owner := wallet.PublicKey()

	// Retrieve the CLMM program ID from the configuration
	clmmProgramID := solana.MustPublicKeyFromBase58(cfg.RaydiumCLMMProgramID)

	// Encode instruction data for the swap
	instructionData, err := encodeSwapInstructionData(amountIn, minAmountOut)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to encode instruction data: %w", err)
	}

	// Construct the swap instruction
	swapInstruction := solana.NewInstruction(
		clmmProgramID,
		solana.AccountMetaSlice{
			{PublicKey: vaultA, IsWritable: true, IsSigner: false},
			{PublicKey: vaultB, IsWritable: true, IsSigner: false},
			{PublicKey: owner, IsWritable: false, IsSigner: true},
		},
		instructionData,
	)

	// Fetch the latest blockhash
	latestBlockhash, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get latest blockhash: %w", err)
	}

	// Create the transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{swapInstruction},
		latestBlockhash.Value.Blockhash,
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

// encodeSwapInstructionData encodes the swap parameters into instruction data.
func encodeSwapInstructionData(amountIn, minAmountOut uint64) ([]byte, error) {
	var buf bytes.Buffer

	discriminator := uint8(1) // Replace with actual opcode for "Swap"
	if err := binary.Write(&buf, binary.LittleEndian, discriminator); err != nil {
		return nil, fmt.Errorf("failed to write discriminator: %w", err)
	}

	if err := binary.Write(&buf, binary.LittleEndian, amountIn); err != nil {
		return nil, fmt.Errorf("failed to write amountIn: %w", err)
	}

	if err := binary.Write(&buf, binary.LittleEndian, minAmountOut); err != nil {
		return nil, fmt.Errorf("failed to write minAmountOut: %w", err)
	}

	return buf.Bytes(), nil
}
