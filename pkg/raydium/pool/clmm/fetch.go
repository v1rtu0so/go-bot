package clmm

import (
	"encoding/json"
	"fmt"
	"os"
)

// FetchClmmPoolFromJSON fetches a CLMM pool from the provided JSON file.
func FetchClmmPoolFromJSON(tokenAddress, filePath string) (*RaydiumClmmPool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var pools []RaydiumClmmPool
	err = json.NewDecoder(file).Decode(&pools)
	if err != nil && err.Error() != "EOF" {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	for _, pool := range pools {
		if pool.ID == tokenAddress {
			return &pool, nil
		}
	}

	return nil, fmt.Errorf("pool not found in JSON")
}

// FetchClmmPoolByID dynamically fetches a CLMM pool based on the token address.
func FetchClmmPoolByID(tokenAddress string) (*RaydiumClmmPool, error) {
	// Simulated fetched data
	return &RaydiumClmmPool{
		ID:             tokenAddress,
		MintProgramIDA: "MintProgramIDAExample",
		MintProgramIDB: "MintProgramIDBExample",
		MintA:          "MintAExample",
		MintB:          "MintBExample",
		VaultA:         "VaultAExample",
		VaultB:         "VaultBExample",
		MintDecimalsA:  6,
		MintDecimalsB:  6,
		AmmConfig: ApiClmmConfigurationItem{
			ID:              "AmmConfigIDExample",
			Index:           1,
			ProtocolFeeRate: 0,
			TradeFeeRate:    25,
			TickSpacing:     64,
			FundFeeRate:     5,
			FundOwner:       "ExampleFundOwner",
		},
		RewardInfos: []RewardInfo{
			{ID: 1, Mint: "RewardMintExample1", ProgramID: "RewardProgramIDExample1"},
			{ID: 2, Mint: "RewardMintExample2", ProgramID: "RewardProgramIDExample2"},
		},
		TVL: 1000000,
		Day: &ApiClmmPoolsItemStatistic{
			ID:         1,
			Volume:     500000,
			VolumeFee:  5000,
			FeeA:       250,
			FeeB:       100,
			FeeApr:     20,
			RewardAprA: 10,
			RewardAprB: 5,
			RewardAprC: 2,
			Apr:        37,
			PriceMin:   1,
			PriceMax:   10,
		},
		LookupTableAccount: "ExampleLookupTableAccount",
	}, nil
}

// FetchAndStorePool dynamically fetches a pool and appends it to the JSON file if not already present.
func FetchAndStorePool(tokenAddress, filePath string) (*RaydiumClmmPool, error) {
	existingPool, err := FetchClmmPoolFromJSON(tokenAddress, filePath)
	if err == nil {
		fmt.Println("CLMM pool data already exists.")
		return existingPool, nil
	}

	fmt.Println("CLMM pool data not found. Fetching dynamically...")
	newPool, err := FetchClmmPoolByID(tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CLMM pool: %w", err)
	}

	err = appendToJSON(filePath, newPool)
	if err != nil {
		return nil, fmt.Errorf("failed to store CLMM pool data: %w", err)
	}

	return newPool, nil
}

// appendToJSON appends a new pool to an existing JSON file.
func appendToJSON(filePath string, pool *RaydiumClmmPool) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var pools []RaydiumClmmPool
	err = json.NewDecoder(file).Decode(&pools)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	pools = append(pools, *pool)

	file.Seek(0, 0)
	if err := json.NewEncoder(file).Encode(pools); err != nil {
		return fmt.Errorf("failed to encode updated data: %w", err)
	}

	return nil
}
