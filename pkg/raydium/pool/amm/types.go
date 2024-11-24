package amm

// RaydiumAmmPool represents the structure of an AMM pool.
type RaydiumAmmPool struct {
	ID                 string `json:"id"`
	BaseMint           string `json:"baseMint"`
	QuoteMint          string `json:"quoteMint"`
	LpMint             string `json:"lpMint"`
	BaseDecimals       int    `json:"baseDecimals"`
	QuoteDecimals      int    `json:"quoteDecimals"`
	LpDecimals         int    `json:"lpDecimals"`
	Version            int    `json:"version"`
	ProgramID          string `json:"programId"`
	Authority          string `json:"authority"`
	OpenOrders         string `json:"openOrders"`
	TargetOrders       string `json:"targetOrders"`
	BaseVault          string `json:"baseVault"`
	QuoteVault         string `json:"quoteVault"`
	WithdrawQueue      string `json:"withdrawQueue"`
	LpVault            string `json:"lpVault"`
	MarketVersion      int    `json:"marketVersion"`
	MarketProgramID    string `json:"marketProgramId"`
	MarketID           string `json:"marketId"`
	MarketAuthority    string `json:"marketAuthority"`
	MarketBaseVault    string `json:"marketBaseVault"`
	MarketQuoteVault   string `json:"marketQuoteVault"`
	MarketBids         string `json:"marketBids"`
	MarketAsks         string `json:"marketAsks"`
	MarketEventQueue   string `json:"marketEventQueue"`
	LookupTableAccount string `json:"lookupTableAccount"`
}
