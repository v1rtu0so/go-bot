package amm

import (
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
)

func parseAmmAccountData(data []byte, programID string, poolID string) (*RaydiumAmmPool, error) {
	if len(data) < ExpectedLength {
		return nil, fmt.Errorf("insufficient data length for parsing: got %d, expected at least %d",
			len(data), ExpectedLength)
	}

	log.Printf("Parsing AMM account data: pool ID %s, data length %d", poolID, len(data))

	// The correct offsets based on Raydium AMM program account layout
	// Reference: https://api-v3.raydium.io/docs/#/amm/get_amm_pool
	pool := &RaydiumAmmPool{
		ID:            poolID,
		ProgramID:     programID,
		BaseMint:      solana.PublicKeyFromBytes(data[0:32]).String(),
		QuoteMint:     solana.PublicKeyFromBytes(data[32:64]).String(),
		LpMint:        solana.PublicKeyFromBytes(data[64:96]).String(),
		BaseVault:     solana.PublicKeyFromBytes(data[96:128]).String(),
		QuoteVault:    solana.PublicKeyFromBytes(data[128:160]).String(),
		LpVault:       solana.PublicKeyFromBytes(data[160:192]).String(),
		Authority:     solana.PublicKeyFromBytes(data[192:224]).String(),
		OpenOrders:    solana.PublicKeyFromBytes(data[224:256]).String(),
		TargetOrders:  solana.PublicKeyFromBytes(data[256:288]).String(),
		WithdrawQueue: solana.PublicKeyFromBytes(data[288:320]).String(),
		// Updated offsets for version and decimals
		Version:       uint8(data[321]), // Changed from 320 to 321
		BaseDecimals:  uint8(data[322]), // Changed from 324 to 322
		QuoteDecimals: uint8(data[323]), // Changed from 325 to 323
		LpDecimals:    uint8(data[324]), // Changed from 326 to 324
	}

	// Changed to uint8 instead of int for proper byte reading
	log.Printf("Parsed pool data: Base Mint: %s, Quote Mint: %s", pool.BaseMint, pool.QuoteMint)
	log.Printf("Vaults - Base: %s, Quote: %s", pool.BaseVault, pool.QuoteVault)
	log.Printf("Decimals - Base: %d, Quote: %d, LP: %d",
		pool.BaseDecimals, pool.QuoteDecimals, pool.LpDecimals)
	log.Printf("Version: %d", pool.Version)

	// Validate parsed data
	if err := validatePoolData(pool); err != nil {
		return nil, fmt.Errorf("invalid pool data: %w", err)
	}

	return pool, nil
}

func validatePoolData(pool *RaydiumAmmPool) error {
	// Basic validation
	if pool.ID == "" || pool.ProgramID == "" {
		return fmt.Errorf("missing required pool identifier")
	}

	if pool.BaseMint == "" || pool.QuoteMint == "" {
		return fmt.Errorf("missing required mint addresses")
	}

	// Validate public key formats
	keys := map[string]string{
		"BaseMint":      pool.BaseMint,
		"QuoteMint":     pool.QuoteMint,
		"LpMint":        pool.LpMint,
		"BaseVault":     pool.BaseVault,
		"QuoteVault":    pool.QuoteVault,
		"LpVault":       pool.LpVault,
		"OpenOrders":    pool.OpenOrders,
		"TargetOrders":  pool.TargetOrders,
		"WithdrawQueue": pool.WithdrawQueue,
	}

	for name, key := range keys {
		if key != "" && !isValidPublicKey(key) {
			return fmt.Errorf("invalid public key format for %s: %s", name, key)
		}
	}

	// Additional validation
	if pool.BaseDecimals < 0 || pool.QuoteDecimals < 0 || pool.LpDecimals < 0 {
		return fmt.Errorf("invalid decimals: base=%d, quote=%d, lp=%d",
			pool.BaseDecimals, pool.QuoteDecimals, pool.LpDecimals)
	}

	return nil
}

func isValidPublicKey(key string) bool {
	if key == "" {
		return true // Allow empty strings as some fields might be optional
	}
	_, err := solana.PublicKeyFromBase58(key)
	return err == nil
}
