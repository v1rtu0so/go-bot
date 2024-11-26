package clmm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// Constants for expected data length
const expectedAccountDataLength = 500 // Replace with the actual expected size of the account data

// FetchClmmPoolFromJSONOrNetwork fetches a CLMM pool from a JSON file or directly via the Solana RPC network.
func FetchClmmPoolFromJSONOrNetwork(
	ctx context.Context,
	client *rpc.Client,
	tokenAddress, programID, filePath string,
) (*RaydiumClmmPool, error) {
	// Try to fetch the pool from the JSON file
	pool, err := FetchClmmPoolFromJSON(tokenAddress, filePath)
	if err == nil {
		fmt.Println("CLMM pool data retrieved from JSON file.")
		return pool, nil
	}

	// If not found in JSON, fetch dynamically from the network
	fmt.Println("CLMM pool data not found in JSON. Fetching dynamically from the Solana network...")
	pool, err = FetchClmmPoolByID(ctx, client, tokenAddress, programID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CLMM pool from network: %w", err)
	}

	// Store the fetched pool data into the JSON file for future use
	err = appendToJSON(filePath, pool)
	if err != nil {
		return nil, fmt.Errorf("failed to store fetched CLMM pool data: %w", err)
	}

	return pool, nil
}

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

// FetchClmmPoolByID fetches a CLMM pool dynamically from the Solana RPC network.
func FetchClmmPoolByID(ctx context.Context, client *rpc.Client, tokenAddress, programID string) (*RaydiumClmmPool, error) {
	accountInfo, err := client.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(tokenAddress))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account info: %w", err)
	}

	// Extract the binary data from accountInfo.Value.Data.Content
	if accountInfo.Value == nil || accountInfo.Value.Data == nil || len(accountInfo.Value.Data.GetBinary()) == 0 {
		return nil, fmt.Errorf("no data found in the account info")
	}
	rawData := accountInfo.Value.Data.GetBinary()

	// Parse the account data to extract pool details
	pool, err := parseClmmAccountData(rawData, programID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CLMM account data: %w", err)
	}

	return pool, nil
}

// parseClmmAccountData parses the raw account data to extract the pool details.
func parseClmmAccountData(accountData []byte, programID string) (*RaydiumClmmPool, error) {
	if len(accountData) < expectedAccountDataLength {
		return nil, fmt.Errorf("account data too short, expected at least %d bytes", expectedAccountDataLength)
	}

	// Simulate parsed pool (replace with real decoding logic)
	pool := &RaydiumClmmPool{
		ID:             programID,
		MintProgramIDA: "RealMintProgramA", // Replace with parsed data
		MintProgramIDB: "RealMintProgramB",
		MintA:          "RealMintA",
		MintB:          "RealMintB",
		VaultA:         "RealVaultA",
		VaultB:         "RealVaultB",
		MintDecimalsA:  6,
		MintDecimalsB:  6,
		AmmConfig: ApiClmmConfigurationItem{
			ID:              "AmmConfigID",
			Index:           1,
			ProtocolFeeRate: 0,
			TradeFeeRate:    25,
			TickSpacing:     64,
			FundFeeRate:     5,
			FundOwner:       "FundOwnerAccount",
		},
		RewardInfos: []RewardInfo{
			{ID: 1, Mint: "RewardMint1", ProgramID: "RewardProgramID1"},
			{ID: 2, Mint: "RewardMint2", ProgramID: "RewardProgramID2"},
		},
		TVL:                123456,
		Day:                &ApiClmmPoolsItemStatistic{Volume: 100000, FeeApr: 15},
		LookupTableAccount: "LookupTableAccount",
	}

	return pool, nil
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

	// Append the new pool
	pools = append(pools, *pool)

	// Write back the updated list to the JSON file
	file.Seek(0, 0)
	if err := json.NewEncoder(file).Encode(pools); err != nil {
		return fmt.Errorf("failed to encode updated data: %w", err)
	}

	return nil
}
