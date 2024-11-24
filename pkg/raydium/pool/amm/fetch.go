package amm

import (
	"encoding/json"
	"fmt"
	"os"
)

// FetchAmmPoolFromJSON fetches an AMM pool from the provided JSON file.
func FetchAmmPoolFromJSON(tokenAddress, filePath string) (*RaydiumAmmPool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var pools []RaydiumAmmPool
	err = json.NewDecoder(file).Decode(&pools)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	for _, pool := range pools {
		if pool.ID == tokenAddress {
			return &pool, nil
		}
	}

	return nil, fmt.Errorf("pool not found in JSON")
}

// FetchAmmPoolByID dynamically fetches an AMM pool based on the token address.
func FetchAmmPoolByID(tokenAddress string) (*RaydiumAmmPool, error) {
	// Simulated fetched data
	return &RaydiumAmmPool{
		ID:                 tokenAddress,
		BaseMint:           "BaseMintExample",
		QuoteMint:          "QuoteMintExample",
		LpMint:             "LpMintExample",
		BaseDecimals:       9,
		QuoteDecimals:      9,
		LpDecimals:         6,
		Version:            4,
		ProgramID:          "ExampleProgramID",
		Authority:          "ExampleAuthority",
		OpenOrders:         "ExampleOpenOrders",
		TargetOrders:       "ExampleTargetOrders",
		BaseVault:          "ExampleBaseVault",
		QuoteVault:         "ExampleQuoteVault",
		WithdrawQueue:      "ExampleWithdrawQueue",
		LpVault:            "ExampleLpVault",
		MarketVersion:      3,
		MarketProgramID:    "ExampleMarketProgramID",
		MarketID:           "ExampleMarketID",
		MarketAuthority:    "ExampleMarketAuthority",
		MarketBaseVault:    "ExampleMarketBaseVault",
		MarketQuoteVault:   "ExampleMarketQuoteVault",
		MarketBids:         "ExampleMarketBids",
		MarketAsks:         "ExampleMarketAsks",
		MarketEventQueue:   "ExampleMarketEventQueue",
		LookupTableAccount: "ExampleLookupTableAccount",
	}, nil
}

// FetchAndStorePool fetches a pool dynamically and stores it in the JSON file.
func FetchAndStorePool(tokenAddress, filePath string) (*RaydiumAmmPool, error) {
	// Check for existing pool in JSON
	existingPool, err := FetchAmmPoolFromJSON(tokenAddress, filePath)
	if err == nil {
		fmt.Println("AMM pool data already exists.")
		return existingPool, nil
	}

	// Dynamically fetch the pool
	fmt.Println("AMM pool data not found. Fetching dynamically...")
	newPool, err := FetchAmmPoolByID(tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch AMM pool: %w", err)
	}

	// Append the new pool data to the JSON file
	err = appendToJSON(filePath, newPool)
	if err != nil {
		return nil, fmt.Errorf("failed to store AMM pool data: %w", err)
	}

	return newPool, nil
}

// appendToJSON appends a new pool to an existing JSON file.
func appendToJSON(filePath string, pool *RaydiumAmmPool) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var pools []RaydiumAmmPool
	err = json.NewDecoder(file).Decode(&pools)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Append the new pool
	pools = append(pools, *pool)

	// Write back the updated list to the JSON file
	file.Seek(0, 0)
	if err := json.NewEncoder(file).Encode(pools); err != nil {
		return fmt.Errorf("failed to encode updated data: %w", err)
	}

	return nil
}
