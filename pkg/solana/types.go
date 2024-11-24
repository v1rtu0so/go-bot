package solana

import (
	"bytes"
	"encoding/binary"
)

// Instruction represents a raw Solana instruction.
type Instruction struct {
	ProgramID [32]byte
	Accounts  []AccountMeta
	Data      []byte
}

// AccountMeta represents metadata about an account.
type AccountMeta struct {
	PublicKey  [32]byte
	IsSigner   bool
	IsWritable bool
}

// Transaction represents a Solana transaction.
type Transaction struct {
	RecentBlockhash [32]byte
	Instructions    []Instruction
	Signatures      [][]byte
	Payer           [32]byte
}

// Serialize serializes the transaction into a byte array.
func (tx *Transaction) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Serialize recent blockhash
	buf.Write(tx.RecentBlockhash[:])

	// Serialize instructions
	for _, inst := range tx.Instructions {
		buf.Write(inst.ProgramID[:])
		buf.Write(encodeAccountMeta(inst.Accounts))
		buf.Write(inst.Data)
	}

	return buf.Bytes(), nil
}

func encodeAccountMeta(accounts []AccountMeta) []byte {
	var buf bytes.Buffer
	for _, account := range accounts {
		buf.Write(account.PublicKey[:])
		binary.Write(&buf, binary.LittleEndian, account.IsSigner)
		binary.Write(&buf, binary.LittleEndian, account.IsWritable)
	}
	return buf.Bytes()
}
