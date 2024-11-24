package helpers

import (
	"context"
	"corvus_bot/pkg/config"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

// WrapSOL wraps the specified amount of SOL into WSOL.
func WrapSOL(ctx context.Context, client *rpc.Client, payer solana.PrivateKey, amountLamports uint64, cfg *config.Config) error {
	wsolMint := solana.MustPublicKeyFromBase58(cfg.WSOLAddress)

	// Check if the ATA exists
	ata, _, err := solana.FindAssociatedTokenAddress(payer.PublicKey(), wsolMint)
	if err != nil {
		return fmt.Errorf("failed to derive ATA: %w", err)
	}
	accountInfo, err := client.GetAccountInfo(ctx, ata)
	if err != nil || accountInfo == nil {
		log.Printf("Associated Token Account not found. Creating ATA...")

		// Create ATA instruction
		createATAInstr := associatedtokenaccount.NewCreateInstruction(
			payer.PublicKey(),
			payer.PublicKey(),
			wsolMint,
		).Build()

		// Build and send the transaction for creating ATA
		tx, err := BuildTransaction(ctx, client, payer, []solana.Instruction{createATAInstr})
		if err != nil {
			return fmt.Errorf("failed to create transaction for ATA: %w", err)
		}

		_, err = client.SendTransaction(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to create ATA: %w", err)
		}

		log.Println("ATA created successfully.")
	}

	// Transfer SOL to WSOL ATA
	transferInstruction := system.NewTransferInstruction(
		amountLamports,
		payer.PublicKey(),
		ata,
	).Build()

	// Sync WSOL account to update balance
	syncInstruction := token.NewSyncNativeInstruction(ata).Build()

	// Build and send the transaction for wrapping
	tx, err := BuildTransaction(ctx, client, payer, []solana.Instruction{transferInstruction, syncInstruction})
	if err != nil {
		return fmt.Errorf("failed to create transaction for wrap: %w", err)
	}

	_, err = client.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to wrap SOL: %w", err)
	}

	log.Printf("Successfully wrapped %d lamports into WSOL", amountLamports)
	return nil
}

// UnwrapSOL unwraps WSOL back into SOL by closing the WSOL account.
func UnwrapSOL(ctx context.Context, client *rpc.Client, payer solana.PrivateKey, cfg *config.Config) error {
	wsolMint := solana.MustPublicKeyFromBase58(cfg.WSOLAddress)

	// Derive the Associated Token Account (ATA)
	ata, _, err := solana.FindAssociatedTokenAddress(payer.PublicKey(), wsolMint)
	if err != nil {
		return fmt.Errorf("failed to derive ATA: %w", err)
	}

	// Check if the ATA exists
	accountInfo, err := client.GetAccountInfo(ctx, ata)
	if err != nil {
		return fmt.Errorf("failed to fetch ATA info: %w", err)
	}
	if accountInfo == nil {
		return fmt.Errorf("no WSOL account found to unwrap")
	}

	// Close the WSOL account
	closeInstruction := token.NewCloseAccountInstruction(
		ata,               // Account to close
		payer.PublicKey(), // Destination for SOL (payer)
		payer.PublicKey(), // Account owner
		nil,               // Additional signers
	).Build()

	// Build and send the transaction for unwrapping
	tx, err := BuildTransaction(ctx, client, payer, []solana.Instruction{closeInstruction})
	if err != nil {
		return fmt.Errorf("failed to create transaction for unwrapping: %w", err)
	}

	_, err = client.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to unwrap WSOL: %w", err)
	}

	log.Println("Successfully unwrapped WSOL back into SOL.")
	return nil
}
