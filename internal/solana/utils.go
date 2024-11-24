package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// CreateTransaction creates a new Solana transaction with instructions.
func CreateTransaction(recentBlockhash solana.Hash, instructions []*solana.Instruction) (*solana.Transaction, error) {
	// Convert []*solana.Instruction to []solana.Instruction
	var insts []solana.Instruction
	for _, inst := range instructions {
		insts = append(insts, *inst)
	}

	tx, err := solana.NewTransaction(
		insts,
		recentBlockhash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return tx, nil
}
