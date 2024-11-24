package solana

import "encoding/binary"

// CreateAccountInstruction creates an instruction for the System Program to create a new account.
func CreateAccountInstruction(fromPubKey, newPubKey [32]byte, lamports uint64, space uint64) Instruction {
	return Instruction{
		ProgramID: SystemProgramID(),
		Accounts: []AccountMeta{
			{PublicKey: fromPubKey, IsSigner: true, IsWritable: true},
			{PublicKey: newPubKey, IsSigner: true, IsWritable: true},
		},
		Data: SystemProgramCreateAccount(lamports, space, TokenProgramID()),
	}
}

// InitializeAccountInstruction creates an instruction to initialize a token account.
func InitializeAccountInstruction(accountPubKey, mintPubKey, ownerPubKey [32]byte) Instruction {
	return Instruction{
		ProgramID: TokenProgramID(),
		Accounts: []AccountMeta{
			{PublicKey: accountPubKey, IsSigner: false, IsWritable: true},
			{PublicKey: mintPubKey, IsSigner: false, IsWritable: false},
			{PublicKey: ownerPubKey, IsSigner: false, IsWritable: false},
			{PublicKey: SysVarRentPubKey(), IsSigner: false, IsWritable: false},
		},
		Data: TokenProgramInitializeAccount(),
	}
}

// CloseAccountInstruction creates an instruction to close a token account.
func CloseAccountInstruction(accountPubKey, destPubKey, ownerPubKey [32]byte) Instruction {
	return Instruction{
		ProgramID: TokenProgramID(),
		Accounts: []AccountMeta{
			{PublicKey: accountPubKey, IsSigner: false, IsWritable: true},
			{PublicKey: destPubKey, IsSigner: false, IsWritable: true},
			{PublicKey: ownerPubKey, IsSigner: true, IsWritable: false},
		},
		Data: TokenProgramCloseAccount(),
	}
}

// SystemProgramCreateAccount creates the data payload for a System Program create account instruction.
func SystemProgramCreateAccount(lamports uint64, space uint64, programID [32]byte) []byte {
	data := make([]byte, 8+8+32)
	binary.LittleEndian.PutUint64(data[:8], lamports)
	binary.LittleEndian.PutUint64(data[8:16], space)
	copy(data[16:], programID[:])
	return data
}

// TokenProgramInitializeAccount creates the data payload for a token account initialization instruction.
func TokenProgramInitializeAccount() []byte {
	return []byte{1} // Initialize Account instruction ID is usually 1 for SPL tokens.
}

// TokenProgramCloseAccount creates the data payload for a close account instruction.
func TokenProgramCloseAccount() []byte {
	return []byte{9} // Close Account instruction ID is usually 9 for SPL tokens.
}

// SystemProgramID returns the public key for the System Program.
func SystemProgramID() [32]byte {
	var id [32]byte
	copy(id[:], "11111111111111111111111111111111")
	return id
}

// TokenProgramID returns the public key for the SPL Token Program.
func TokenProgramID() [32]byte {
	var id [32]byte
	copy(id[:], "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	return id
}

// SysVarRentPubKey returns the public key for the Rent system variable.
func SysVarRentPubKey() [32]byte {
	var id [32]byte
	copy(id[:], "SysvarRent111111111111111111111111111111111")
	return id
}
