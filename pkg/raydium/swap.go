package raydium

import (
	"context"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// SwapTokens swaps tokenA for tokenB in a Raydium pool.
func (rc *RaydiumClient) SwapTokens(ctx context.Context, wallet solana.PrivateKey, tokenA, tokenB string, amount uint64) (string, error) {
	log.Printf("Swapping %d of %s for %s", amount, tokenA, tokenB)

	// Fetch the recent blockhash
	blockhash, err := rc.RPCClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Printf("Failed to fetch recent blockhash: %v", err)
		return "", err
	}

	// Construct the transaction instruction for the swap
	txInstruction := solana.NewInstruction(
		rc.ProgramID,
		[]*solana.AccountMeta{
			solana.NewAccountMeta(wallet.PublicKey(), true, true),
		},
		[]byte{0x01}, // Placeholder: Actual instruction data for a Raydium swap
	)

	// Construct the transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{txInstruction},
		blockhash.Value.Blockhash,
	)
	if err != nil {
		log.Printf("Failed to construct transaction: %v", err)
		return "", err
	}

	// Sign the transaction
	tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(wallet.PublicKey()) {
			return &wallet
		}
		return nil
	})

	// Send the transaction
	sig, err := rc.RPCClient.SendTransaction(ctx, tx)
	if err != nil {
		log.Printf("Failed to send transaction: %v", err)
		return "", err
	}

	return sig.String(), nil
}
