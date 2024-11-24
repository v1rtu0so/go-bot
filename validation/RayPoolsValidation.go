package validation

import (
	"encoding/json"
	"fmt"
	"os"
)

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

type ApiClmmConfigurationItem struct {
	ID              string `json:"id"`
	Index           int    `json:"index"`
	ProtocolFeeRate int    `json:"protocolFeeRate"`
	TradeFeeRate    int    `json:"tradeFeeRate"`
	TickSpacing     int    `json:"tickSpacing"`
	FundFeeRate     int    `json:"fundFeeRate"`
	FundOwner       string `json:"fundOwner"`
}

type RewardInfo struct {
	ID        int    `json:"id"`
	Mint      string `json:"mint"`
	ProgramID string `json:"programId"`
}

type ApiClmmPoolsItemStatistic struct {
	ID         int `json:"id"`
	Volume     int `json:"volume"`
	VolumeFee  int `json:"volumeFee"`
	FeeA       int `json:"feeA"`
	FeeB       int `json:"feeB"`
	FeeApr     int `json:"feeApr"`
	RewardAprA int `json:"rewardAprA"`
	RewardAprB int `json:"rewardAprB"`
	RewardAprC int `json:"rewardAprC"`
	Apr        int `json:"apr"`
	PriceMin   int `json:"priceMin"`
	PriceMax   int `json:"priceMax"`
}

type RaydiumClmmPool struct {
	ID                 string                     `json:"id"`
	MintProgramIDA     string                     `json:"mintProgramIdA"`
	MintProgramIDB     string                     `json:"mintProgramIdB"`
	MintA              string                     `json:"mintA"`
	MintB              string                     `json:"mintB"`
	VaultA             string                     `json:"vaultA"`
	VaultB             string                     `json:"vaultB"`
	MintDecimalsA      int                        `json:"mintDecimalsA"`
	MintDecimalsB      int                        `json:"mintDecimalsB"`
	AmmConfig          ApiClmmConfigurationItem   `json:"ammConfig"`
	RewardInfos        []RewardInfo               `json:"rewardInfos"`
	TVL                int                        `json:"tvl"`
	Day                *ApiClmmPoolsItemStatistic `json:"day,omitempty"`
	Week               *ApiClmmPoolsItemStatistic `json:"week,omitempty"`
	Month              *ApiClmmPoolsItemStatistic `json:"month,omitempty"`
	LookupTableAccount string                     `json:"lookupTableAccount"`
}

// validateJSON checks the JSON data for a given type.
func validateJSON(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&target); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

func main() {
	// Validate AMM JSON
	ammPools := []RaydiumAmmPool{}
	ammFilePath := "../testdata/amm_pools.json"
	if err := validateJSON(ammFilePath, &ammPools); err != nil {
		fmt.Printf("AMM JSON validation failed: %s\n", err)
	} else {
		fmt.Println("AMM JSON validation passed!")
	}

	// Validate CLMM JSON
	clmmPools := []RaydiumClmmPool{}
	clmmFilePath := "../testdata/clmm_pools.json"
	if err := validateJSON(clmmFilePath, &clmmPools); err != nil {
		fmt.Printf("CLMM JSON validation failed: %s\n", err)
	} else {
		fmt.Println("CLMM JSON validation passed!")
	}

	// Additional validation checks (e.g., required fields)
	for _, pool := range ammPools {
		if pool.ID == "" || pool.BaseMint == "" || pool.ProgramID == "" {
			fmt.Printf("Invalid AMM Pool: %+v\n", pool)
		}
	}

	for _, pool := range clmmPools {
		if pool.ID == "" || pool.MintProgramIDA == "" || pool.MintA == "" {
			fmt.Printf("Invalid CLMM Pool: %+v\n", pool)
		}
	}
}
