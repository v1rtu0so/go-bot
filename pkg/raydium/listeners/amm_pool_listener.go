// pkg/raydium/listeners/amm_pool_listener.go

package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"corvus_bot/pkg/config"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// Add new constants for known program IDs
const (
	TokenProgramID           = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	SystemProgramID          = "11111111111111111111111111111111"
	AssociatedTokenProgramID = "ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"
	SerumProgramID           = "srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX"
)

// TransactionAccount represents an account involved in the transaction
type TransactionAccount struct {
	Pubkey     string `json:"pubkey"`
	IsSigner   bool   `json:"is_signer"`
	IsWritable bool   `json:"is_writable"`
}

// AccountRole represents the role of an account in the pool
type AccountRole string

const (
	RoleMint         AccountRole = "MINT"
	RoleTokenAccount AccountRole = "TOKEN_ACCOUNT"
	RoleProgram      AccountRole = "PROGRAM"
	RoleMarket       AccountRole = "MARKET"
	RoleAuthority    AccountRole = "AUTHORITY"
	RoleUser         AccountRole = "USER"
)

// AccountMeta represents a program account with its metadata
type AccountMeta struct {
	Address     string      `json:"address"`
	IsSigner    bool        `json:"is_signer"`
	IsMutable   bool        `json:"is_mutable"`
	Name        string      `json:"name"`
	ProgramName string      `json:"program_name,omitempty"`
	Role        AccountRole `json:"role"`
	TokenInfo   *TokenInfo  `json:"token_info,omitempty"`
}

// TokenInfo represents token-specific metadata
type TokenInfo struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Decimals  uint8  `json:"decimals"`
	Supply    uint64 `json:"supply,omitempty"`
	Authority string `json:"authority,omitempty"`
}

// RawPoolEvent represents a pool event from the blockchain
type RawPoolEvent struct {
	Signature  string                    `json:"signature"`
	Slot       uint64                    `json:"slot"`
	BlockTime  int64                     `json:"blockTime"`
	Logs       []string                  `json:"logs"`
	Timestamp  time.Time                 `json:"timestamp"`
	InitParams *InitializeParams         `json:"init_params,omitempty"`
	ParsedData *ParsedPoolInitialization `json:"parsed_data,omitempty"`
}

// InitializeParams represents the parameters from initialize2 instruction
type InitializeParams struct {
	Nonce          uint8  `json:"nonce"`
	OpenTime       int64  `json:"open_time"`
	InitPcAmount   uint64 `json:"init_pc_amount"`
	InitCoinAmount uint64 `json:"init_coin_amount"`
}

// PoolAccounts represents all the accounts involved in pool initialization
type PoolAccounts struct {
	TokenProgram           AccountMeta `json:"token_program"`
	AssociatedTokenProgram AccountMeta `json:"associated_token_program"`
	SystemProgram          AccountMeta `json:"system_program"`
	AMM                    AccountMeta `json:"amm"`
	AMMAuthority           AccountMeta `json:"amm_authority"`
	AMMOpenOrders          AccountMeta `json:"amm_open_orders"`
	LPMint                 AccountMeta `json:"lp_mint"`
	CoinMint               AccountMeta `json:"coin_mint"`
	PCMint                 AccountMeta `json:"pc_mint"`
	PoolCoinTokenAccount   AccountMeta `json:"pool_coin_token_account"`
	PoolPCTokenAccount     AccountMeta `json:"pool_pc_token_account"`
	PoolWithdrawQueue      AccountMeta `json:"pool_withdraw_queue"`
	AMMTargetOrders        AccountMeta `json:"amm_target_orders"`
	SerumMarket            AccountMeta `json:"serum_market"`
	SerumProgram           AccountMeta `json:"serum_program"`
	UserWallet             AccountMeta `json:"user_wallet"`
}

// ParsedPoolInitialization represents the complete parsed pool initialization data
type ParsedPoolInitialization struct {
	Signature    string           `json:"signature"`
	Timestamp    time.Time        `json:"timestamp"`
	Slot         uint64           `json:"slot"`
	BlockTime    int64            `json:"block_time"`
	Accounts     PoolAccounts     `json:"accounts"`
	InitParams   InitializeParams `json:"init_params"`
	TokenDetails struct {
		CoinMintDecimals uint8   `json:"coin_mint_decimals"`
		PCMintDecimals   uint8   `json:"pc_mint_decimals"`
		LPMintDecimals   uint8   `json:"lp_mint_decimals"`
		InitialPrice     float64 `json:"initial_price"`
		CoinAmount       float64 `json:"coin_amount"`
		PCAmount         float64 `json:"pc_amount"`
	} `json:"token_details"`
	Status string `json:"status"`
}

// HumanReadablePool provides a clean, readable representation of the pool
type HumanReadablePool struct {
	Signature        string    `json:"signature"`
	Timestamp        time.Time `json:"timestamp"`
	TokenA           string    `json:"token_a"`
	TokenB           string    `json:"token_b"`
	InitialLiquidity struct {
		TokenAAmount float64 `json:"token_a_amount"`
		TokenBAmount float64 `json:"token_b_amount"`
	} `json:"initial_liquidity"`
	InitialPrice float64 `json:"initial_price"`
	Creator      string  `json:"creator"`
	PoolID       string  `json:"pool_id"`
	Market       string  `json:"market"`
}

// AMMPoolListener manages the WebSocket subscription to pool events
type AMMPoolListener struct {
	wsClient  *ws.Client
	eventChan chan *RawPoolEvent
	config    *config.Config
}

// NewAMMPoolListener creates a new listener instance
func NewAMMPoolListener(cfg *config.Config) (*AMMPoolListener, error) {
	wsClient, err := ws.Connect(context.Background(), cfg.WSConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	return &AMMPoolListener{
		wsClient:  wsClient,
		eventChan: make(chan *RawPoolEvent, 100),
		config:    cfg,
	}, nil
}

// Start begins listening for pool events
func (l *AMMPoolListener) Start(ctx context.Context) error {
	programID := solana.MustPublicKeyFromBase58(l.config.RaydiumAMMProgramID)

	sub, err := l.wsClient.LogsSubscribeMentions(
		programID,
		rpc.CommitmentProcessed,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}
	defer sub.Unsubscribe()

	log.Printf("Subscribed to Raydium AMM logs for program: %s", programID)

	file, err := os.OpenFile("raydium_filtered_logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	for {
		select {
		case <-ctx.Done():
			return nil
		case resp, ok := <-sub.Response():
			if !ok {
				return fmt.Errorf("subscription channel closed")
			}
			if resp.Value.Err != nil {
				log.Printf("Error in log response: %v", resp.Value.Err)
				continue
			}

			if isPoolInitEvent(resp.Value.Logs) {
				var initParams *InitializeParams
				for _, logLine := range resp.Value.Logs {
					if params, err := parseInitializeParams(logLine); err == nil && params != nil {
						initParams = params
						break
					}
				}

				parsedData, err := parsePoolAccounts(resp.Value.Logs)
				if err != nil {
					log.Printf("Error parsing pool accounts: %v", err)
				}

				event := &RawPoolEvent{
					Signature:  resp.Value.Signature.String(),
					Slot:       resp.Context.Slot,
					BlockTime:  time.Now().Unix(),
					Logs:       resp.Value.Logs,
					Timestamp:  time.Now().UTC(),
					InitParams: initParams,
					ParsedData: parsedData,
				}

				select {
				case l.eventChan <- event:
				default:
					log.Printf("Warning: Event channel full, dropping event %s", event.Signature)
				}

				if err := encoder.Encode(event); err != nil {
					log.Printf("Error saving event to file: %v", err)
				}

				log.Printf("Captured pool initialization event: %s", event.Signature)
			}
		}
	}
}

// GetEventChannel returns the channel for receiving pool events
func (l *AMMPoolListener) GetEventChannel() chan *RawPoolEvent {
	return l.eventChan
}

// Close closes the WebSocket connection
func (l *AMMPoolListener) Close() {
	l.wsClient.Close()
}

// Helper functions

func isPoolInitEvent(logs []string) bool {
	for _, log := range logs {
		if strings.Contains(log, "initialize2") {
			return true
		}
	}
	return false
}

func parseInitializeParams(log string) (*InitializeParams, error) {
	if !strings.Contains(log, "initialize2:") {
		return nil, nil
	}

	var params InitializeParams
	start := strings.Index(log, "{")
	end := strings.Index(log, "}")

	if start == -1 || end == -1 {
		return nil, fmt.Errorf("invalid initialize2 log format")
	}

	paramsStr := log[start : end+1]
	paramsStr = strings.ReplaceAll(paramsStr, " ", "")

	nonceRe := regexp.MustCompile(`nonce:(\d+)`)
	openTimeRe := regexp.MustCompile(`open_time:(\d+)`)
	pcAmountRe := regexp.MustCompile(`init_pc_amount:(\d+)`)
	coinAmountRe := regexp.MustCompile(`init_coin_amount:(\d+)`)

	if m := nonceRe.FindStringSubmatch(paramsStr); len(m) > 1 {
		nonce, _ := strconv.ParseUint(m[1], 10, 8)
		params.Nonce = uint8(nonce)
	}
	if m := openTimeRe.FindStringSubmatch(paramsStr); len(m) > 1 {
		params.OpenTime, _ = strconv.ParseInt(m[1], 10, 64)
	}
	if m := pcAmountRe.FindStringSubmatch(paramsStr); len(m) > 1 {
		params.InitPcAmount, _ = strconv.ParseUint(m[1], 10, 64)
	}
	if m := coinAmountRe.FindStringSubmatch(paramsStr); len(m) > 1 {
		params.InitCoinAmount, _ = strconv.ParseUint(m[1], 10, 64)
	}

	return &params, nil
}

// Update parsePoolAccounts to directly handle account mapping without helper function
func parsePoolAccounts(logs []string) (*ParsedPoolInitialization, error) {
	parsed := &ParsedPoolInitialization{
		Status:   "success",
		Accounts: PoolAccounts{},
	}

	accountMap := make(map[string]*TransactionAccount)

	// First pass: Collect all account addresses and their properties
	for _, log := range logs {
		// Parse program invocations
		if strings.Contains(log, "invoke") {
			programID := extractProgramID(log)
			if programID != "" {
				accountMap[programID] = &TransactionAccount{
					Pubkey:     programID,
					IsSigner:   false,
					IsWritable: false,
				}
			}
			continue
		}

		// Parse InitializeAccount instructions
		if strings.Contains(log, "Instruction: InitializeAccount") {
			if mint, _ := extractTokenAccounts(log); mint != "" {
				accountMap[mint] = &TransactionAccount{
					Pubkey:     mint,
					IsSigner:   false,
					IsWritable: true,
				}
			}
			continue
		}

		// Parse InitializeMint instructions
		if strings.Contains(log, "Instruction: InitializeMint") {
			// Extract mint address from subsequent log lines
			continue
		}

		// Parse Transfer instructions
		if strings.Contains(log, "Instruction: Transfer") {
			source, dest, _ := extractTransferDetails(log)
			if source != "" {
				accountMap[source] = &TransactionAccount{
					Pubkey:     source,
					IsSigner:   false,
					IsWritable: true,
				}
			}
			if dest != "" {
				accountMap[dest] = &TransactionAccount{
					Pubkey:     dest,
					IsSigner:   false,
					IsWritable: true,
				}
			}
			continue
		}
	}

	// Second pass: Map accounts to their roles directly
	for pubkey, account := range accountMap {
		switch {
		case pubkey == TokenProgramID:
			parsed.Accounts.TokenProgram = AccountMeta{
				Address:     pubkey,
				IsSigner:    account.IsSigner,
				IsMutable:   account.IsWritable,
				ProgramName: "Token Program",
				Role:        RoleProgram,
			}
		case pubkey == SystemProgramID:
			parsed.Accounts.SystemProgram = AccountMeta{
				Address:     pubkey,
				IsSigner:    account.IsSigner,
				IsMutable:   account.IsWritable,
				ProgramName: "System Program",
				Role:        RoleProgram,
			}
		case pubkey == AssociatedTokenProgramID:
			parsed.Accounts.AssociatedTokenProgram = AccountMeta{
				Address:     pubkey,
				IsSigner:    account.IsSigner,
				IsMutable:   account.IsWritable,
				ProgramName: "Associated Token Program",
				Role:        RoleProgram,
			}
		case pubkey == SerumProgramID:
			parsed.Accounts.SerumProgram = AccountMeta{
				Address:     pubkey,
				IsSigner:    account.IsSigner,
				IsMutable:   account.IsWritable,
				ProgramName: "Serum Program",
				Role:        RoleProgram,
			}
		}
	}

	// Extract and set pool parameters
	for _, logLine := range logs {
		if strings.Contains(logLine, "initialize2:") {
			if params, err := parseInitializeParams(logLine); err == nil && params != nil {
				parsed.InitParams = *params
				break
			}
		}
	}

	// Parse ray_log if present
	for _, logLine := range logs {
		if strings.Contains(logLine, "ray_log:") {
			if rayData := parseRayLog(logLine); rayData != nil {
				// Update parsed data with ray_log information
			}
		}
	}

	return parsed, nil
}

func extractProgramID(log string) string {
	re := regexp.MustCompile(`Program (\w+) invoke`)
	if matches := re.FindStringSubmatch(log); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractAccount(log string, currentProgram string) (*AccountMeta, error) {
	addrRe := regexp.MustCompile(`([1-9A-HJ-NP-Za-km-z]{32,44})`)
	matches := addrRe.FindStringSubmatch(log)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no account address found")
	}

	return &AccountMeta{
		Address:     matches[1],
		IsSigner:    strings.Contains(log, "signer"),
		IsMutable:   strings.Contains(log, "mut"),
		ProgramName: currentProgram,
		Role:        determineAccountRole(log),
	}, nil
}

func determineAccountRole(log string) AccountRole {
	switch {
	case strings.Contains(log, "InitializeMint"):
		return RoleMint
	case strings.Contains(log, "InitializeAccount"):
		return RoleTokenAccount
	case strings.Contains(log, "authority"):
		return RoleAuthority
	case strings.Contains(log, "market"):
		return RoleMarket
	default:
		return RoleProgram
	}
}

func (p *ParsedPoolInitialization) ToHumanReadable() *HumanReadablePool {
	return &HumanReadablePool{
		Signature: p.Signature,
		Timestamp: p.Timestamp,
		TokenA:    p.Accounts.CoinMint.TokenInfo.Symbol,
		TokenB:    p.Accounts.PCMint.TokenInfo.Symbol,
		InitialLiquidity: struct {
			TokenAAmount float64 `json:"token_a_amount"`
			TokenBAmount float64 `json:"token_b_amount"`
		}{
			TokenAAmount: p.TokenDetails.CoinAmount,
			TokenBAmount: p.TokenDetails.PCAmount,
		},
		InitialPrice: p.TokenDetails.InitialPrice,
		Creator:      p.Accounts.UserWallet.Address,
		PoolID:       p.Accounts.AMM.Address,
		Market:       p.Accounts.SerumMarket.Address,
	}
}

// Add new helper function to extract token accounts
func extractTokenAccounts(log string) (mint string, owner string) {
	// Extract addresses from InitializeAccount instructions
	accountRe := regexp.MustCompile(`Initialize (?:token|mint) account (\w{32,44})`)
	if matches := accountRe.FindStringSubmatch(log); len(matches) > 1 {
		return matches[1], ""
	}
	return "", ""
}

// Add new helper function to extract transfer details
func extractTransferDetails(log string) (source string, destination string, amount uint64) {
	// Extract transfer details from Transfer instructions
	return "", "", 0
}

// Update parseRayLog to remove unused parameter
func parseRayLog(rawLog string) map[string]interface{} {
	if !strings.Contains(rawLog, "ray_log:") {
		return nil
	}

	// Extract and decode the base64 data after "ray_log:"
	parts := strings.Split(rawLog, "ray_log: ")
	if len(parts) != 2 {
		return nil
	}

	// TODO: Implement ray_log parsing
	return nil
}
