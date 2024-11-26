package amm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func FetchAmmPoolFromJSONOrNetwork(ctx context.Context, client *rpc.Client, baseMint, quoteMint, programID, filePath string) (*RaydiumAmmPool, error) {
	// First try to fetch from JSON
	pool, err := FetchAmmPoolFromJSON(baseMint, quoteMint, filePath)
	if err == nil {
		log.Printf("Found pool in JSON storage: %s", pool.ID)
		return pool, nil
	}

	log.Println("Pool data not found in JSON. Fetching dynamically from the Raydium API...")

	// Get pool ID from API
	poolID, poolProgramID, err := FetchAmmPoolIDFromAPI(baseMint, quoteMint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pool ID from API: %w", err)
	}

	log.Printf("Found pool ID from API: %s, fetching on-chain data...", poolID)

	// Fetch pool account data from chain
	accountInfo, err := client.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(poolID))
	if err != nil {
		return nil, fmt.Errorf("failed to get account info from chain: %w", err)
	}

	if accountInfo == nil || accountInfo.Value == nil {
		return nil, fmt.Errorf("no account data found for pool ID: %s", poolID)
	}

	rawData := accountInfo.Value.Data.GetBinary()
	log.Printf("Retrieved %d bytes of account data from chain", len(rawData))

	// Parse account data
	pool, err = parseAmmAccountData(rawData, poolProgramID, poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse on-chain data: %w", err)
	}

	// Log pool data for verification
	log.Printf("Successfully parsed pool data. ID: %s, Base: %s, Quote: %s",
		pool.ID, pool.BaseMint, pool.QuoteMint)

	// Store the data
	err = StorePoolData(pool, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to store pool data: %w", err)
	}

	return pool, nil
}

func FetchAmmPoolIDFromAPI(baseMint, quoteMint string) (string, string, error) {
	url := fmt.Sprintf("https://api-v3.raydium.io/pools/info/mint?mint1=%s&mint2=%s&poolType=standard&poolSortField=default&sortType=desc&pageSize=1&page=1",
		baseMint, quoteMint)

	log.Printf("Fetching from API: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp RaydiumAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", "", fmt.Errorf("failed to parse API response: %w", err)
	}

	if !apiResp.Success || len(apiResp.Data.Data) == 0 {
		return "", "", fmt.Errorf("no pool found")
	}

	poolData := apiResp.Data.Data[0]
	log.Printf("API returned pool ID: %s with program ID: %s",
		poolData.ID, poolData.ProgramID)
	return poolData.ID, poolData.ProgramID, nil
}
