package amm

// RaydiumAmmPool represents the structure of an AMM pool as fetched from the on-chain program
type RaydiumAmmPool struct {
	ID            string `json:"id"`            // Pool address
	ProgramID     string `json:"programId"`     // Raydium AMM Program ID
	BaseMint      string `json:"baseMint"`      // Address of base mint
	QuoteMint     string `json:"quoteMint"`     // Address of quote mint
	LpMint        string `json:"lpMint"`        // LP token mint
	BaseVault     string `json:"baseVault"`     // Base token vault
	QuoteVault    string `json:"quoteVault"`    // Quote token vault
	LpVault       string `json:"lpVault"`       // LP token vault
	Authority     string `json:"authority"`     // Authority for the pool
	OpenOrders    string `json:"openOrders"`    // Serum OpenOrders account
	TargetOrders  string `json:"targetOrders"`  // Serum TargetOrders account
	WithdrawQueue string `json:"withdrawQueue"` // Withdraw queue for the pool
	Version       uint8  `json:"version"`       // Changed from int to uint8
	BaseDecimals  uint8  `json:"baseDecimals"`  // Changed from int to uint8
	QuoteDecimals uint8  `json:"quoteDecimals"` // Changed from int to uint8
	LpDecimals    uint8  `json:"lpDecimals"`    // Changed from int to uint8
}

// RaydiumAPIResponse represents the structure of the Raydium v3 API response
type RaydiumAPIResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Data []struct {
			Type      string `json:"type"`
			ID        string `json:"id"`
			ProgramID string `json:"programId"`
			MintA     struct {
				Address  string `json:"address"`
				Decimals int    `json:"decimals"`
			} `json:"mintA"`
			MintB struct {
				Address  string `json:"address"`
				Decimals int    `json:"decimals"`
			} `json:"mintB"`
			LpMint struct {
				Address  string `json:"address"`
				Decimals int    `json:"decimals"`
			} `json:"lpMint"`
			MarketID string `json:"marketId"` // OpenOrders account
		} `json:"data"`
	} `json:"data"`
}
