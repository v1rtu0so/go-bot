package solana

import "encoding/binary"

// CreateAccountInstruction creates a `CreateAccount` instruction.
func CreateAccountInstruction(payer, newAccount [32]byte, lamports, space uint64) Instruction {
	programID := SystemProgramID()

	// Create the instruction data
	data := make([]byte, 16)
	binary.LittleEndian.PutUint64(data[:8], lamports)
	binary.LittleEndian.PutUint64(data[8:], space)

	return Instruction{
		ProgramID: programID,
		Accounts: []AccountMeta{
			{PublicKey: payer, IsSigner: true, IsWritable: true},
			{PublicKey: newAccount, IsSigner: true, IsWritable: true},
		},
		Data: data,
	}
}

// InitializeAccountInstruction creates an instruction to initialize an account.
func InitializeAccountInstruction(account, owner [32]byte) Instruction {
	programID := TokenProgramID()
	data := []byte{1} // Specific instruction layout for InitializeAccount
	return Instruction{
		ProgramID: programID,
		Accounts: []AccountMeta{
			{PublicKey: account, IsSigner: false, IsWritable: true},
			{PublicKey: owner, IsSigner: true, IsWritable: false},
		},
		Data: data,
	}
}

// SystemProgramID returns the hardcoded program ID for the System Program.
func SystemProgramID() [32]byte {
	var id [32]byte
	copy(id[:], "11111111111111111111111111111111")
	return id
}

// TokenProgramID returns the hardcoded program ID for the Token Program.
func TokenProgramID() [32]byte {
	var id [32]byte
	copy(id[:], "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	return id
}
