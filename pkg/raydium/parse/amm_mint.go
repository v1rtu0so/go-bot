//mnt/dev/go-bot/pkg/raydium/parse/amm_mint.go

package parse

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// AMMInitInstructionData represents the IDL structure for initialize2
type AMMInitInstructionData struct {
	Nonce          uint8
	OpenTime       uint64
	InitPcAmount   uint64
	InitCoinAmount uint64
}

// AMMPoolAccounts represents the account structure from the IDL
type AMMPoolAccounts struct {
	TokenProgram    solana.PublicKey // account[0]
	AssociatedToken solana.PublicKey // account[1]
	SystemProgram   solana.PublicKey // account[2]
	Rent            solana.PublicKey // account[3]
	Pool            solana.PublicKey // account[4]  - Main pool account
	Authority       solana.PublicKey // account[5]
	OpenOrders      solana.PublicKey // account[6]
	LPMint          solana.PublicKey // account[7]
	CoinMint        solana.PublicKey // account[8]  - Base token
	PCMint          solana.PublicKey // account[9]  - Quote token
	CoinVault       solana.PublicKey // account[10] - Base vault
	PCVault         solana.PublicKey // account[11] - Quote vault
	WithdrawQueue   solana.PublicKey // account[12]
	TargetOrders    solana.PublicKey // account[13]
	LPTokenAccount  solana.PublicKey // account[14]
}

// ParsedAMMPool represents the final parsed pool data
type ParsedAMMPool struct {
	ID            string
	ProgramID     string
	Version       uint8
	BaseMint      string
	QuoteMint     string
	LPMint        string
	BaseVault     string
	QuoteVault    string
	Authority     string
	OpenOrders    string
	TargetOrders  string
	WithdrawQueue string
	BaseDecimals  uint8
	QuoteDecimals uint8
	LPDecimals    uint8
	InitialBase   uint64
	InitialQuote  uint64
}

// AMMParser handles parsing of AMM initialization instructions
type AMMParser struct {
	programID solana.PublicKey
}

// NewAMMParser creates a new parser instance
func NewAMMParser(programID string) (*AMMParser, error) {
	pid, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("invalid program ID: %w", err)
	}
	return &AMMParser{programID: pid}, nil
}

// ParsePoolInit parses pool initialization from transaction logs
func (p *AMMParser) ParsePoolInit(logMsg *ws.LogResult) (*ParsedAMMPool, error) {
	// First validate we have an initialize2 instruction
	var initLog string
	for _, log := range logMsg.Logs {
		if strings.Contains(log, "initialize2") {
			initLog = log
			break
		}
	}
	if initLog == "" {
		return nil, fmt.Errorf("no initialize2 instruction found in logs")
	}

	// Parse instruction data from the log
	instructionData, err := p.parseInstructionData(initLog)
	if err != nil {
		return nil, fmt.Errorf("failed to parse instruction data: %w", err)
	}

	// Parse accounts from the transaction
	accounts, err := p.parseAccounts(logMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse accounts: %w", err)
	}

	// Construct the parsed pool data
	pool := &ParsedAMMPool{
		ID:            accounts.Pool.String(),
		ProgramID:     p.programID.String(),
		Version:       1, // Initialize2 is v1
		BaseMint:      accounts.CoinMint.String(),
		QuoteMint:     accounts.PCMint.String(),
		LPMint:        accounts.LPMint.String(),
		BaseVault:     accounts.CoinVault.String(),
		QuoteVault:    accounts.PCVault.String(),
		Authority:     accounts.Authority.String(),
		OpenOrders:    accounts.OpenOrders.String(),
		TargetOrders:  accounts.TargetOrders.String(),
		WithdrawQueue: accounts.WithdrawQueue.String(),
		InitialBase:   instructionData.InitCoinAmount,
		InitialQuote:  instructionData.InitPcAmount,
	}

	if err := p.validatePoolData(pool); err != nil {
		return nil, fmt.Errorf("pool validation failed: %w", err)
	}

	log.Printf("Successfully parsed AMM pool initialization: ID=%s, Base=%s, Quote=%s",
		pool.ID, pool.BaseMint, pool.QuoteMint)

	return pool, nil
}

// parseInstructionData extracts initialize2 parameters from the log
func (p *AMMParser) parseInstructionData(log string) (*AMMInitInstructionData, error) {
	// Extract base64 data from log if present
	parts := strings.Split(log, "data: ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid instruction log format")
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode instruction data: %w", err)
	}

	if len(data) < 25 { // 1 (nonce) + 8 (openTime) + 8 (pcAmount) + 8 (coinAmount)
		return nil, fmt.Errorf("insufficient instruction data length")
	}

	return &AMMInitInstructionData{
		Nonce:          data[0],
		OpenTime:       binary.LittleEndian.Uint64(data[1:9]),
		InitPcAmount:   binary.LittleEndian.Uint64(data[9:17]),
		InitCoinAmount: binary.LittleEndian.Uint64(data[17:25]),
	}, nil
}

// parseAccounts extracts and validates account addresses
func (p *AMMParser) parseAccounts(logMsg *ws.LogResult) (*AMMPoolAccounts, error) {
	if len(logMsg.Messages) < 15 {
		return nil, fmt.Errorf("insufficient account keys: got %d, need 15", len(logMsg.Messages))
	}

	accounts := &AMMPoolAccounts{}
	// Map accounts based on IDL order
	keys := logMsg.Messages[:15] // Take first 15 accounts as per IDL
	for i, key := range keys {
		pubKey, err := solana.PublicKeyFromBase58(key)
		if err != nil {
			return nil, fmt.Errorf("invalid public key at index %d: %w", i, err)
		}

		switch i {
		case 0:
			accounts.TokenProgram = pubKey
		case 1:
			accounts.AssociatedToken = pubKey
		case 2:
			accounts.SystemProgram = pubKey
		case 3:
			accounts.Rent = pubKey
		case 4:
			accounts.Pool = pubKey
		case 5:
			accounts.Authority = pubKey
		case 6:
			accounts.OpenOrders = pubKey
		case 7:
			accounts.LPMint = pubKey
		case 8:
			accounts.CoinMint = pubKey
		case 9:
			accounts.PCMint = pubKey
		case 10:
			accounts.CoinVault = pubKey
		case 11:
			accounts.PCVault = pubKey
		case 12:
			accounts.WithdrawQueue = pubKey
		case 13:
			accounts.TargetOrders = pubKey
		case 14:
			accounts.LPTokenAccount = pubKey
		}
	}

	return accounts, nil
}

// validatePoolData performs comprehensive validation of parsed pool data
func (p *AMMParser) validatePoolData(pool *ParsedAMMPool) error {
	// Validate mint addresses
	if pool.BaseMint == pool.QuoteMint {
		return fmt.Errorf("base and quote mints cannot be identical")
	}

	// Validate vaults are different
	if pool.BaseVault == pool.QuoteVault {
		return fmt.Errorf("base and quote vaults cannot be identical")
	}

	// Check for zero addresses
	zeroAddr := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
	checks := map[string]string{
		"BaseMint":   pool.BaseMint,
		"QuoteMint":  pool.QuoteMint,
		"BaseVault":  pool.BaseVault,
		"QuoteVault": pool.QuoteVault,
		"Authority":  pool.Authority,
	}

	for name, addr := range checks {
		pubKey, err := solana.PublicKeyFromBase58(addr)
		if err != nil {
			return fmt.Errorf("invalid %s address: %w", name, err)
		}
		if pubKey.Equals(zeroAddr) {
			return fmt.Errorf("%s cannot be system program address", name)
		}
	}

	return nil
}
