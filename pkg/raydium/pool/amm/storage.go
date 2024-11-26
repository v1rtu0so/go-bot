package amm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func FetchAmmPoolFromJSON(baseMint, quoteMint, filePath string) (*RaydiumAmmPool, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create empty file with array structure
			if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
				return nil, fmt.Errorf("failed to create initial JSON file: %w", err)
			}
			return nil, fmt.Errorf("pool not found in JSON")
		}
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	// Read and parse JSON file
	var pools []RaydiumAmmPool
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pools); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Search for pool with matching mints in either direction
	for _, pool := range pools {
		if (pool.BaseMint == baseMint && pool.QuoteMint == quoteMint) ||
			(pool.BaseMint == quoteMint && pool.QuoteMint == baseMint) {
			log.Printf("Found pool in JSON storage: %s", pool.ID)
			return &pool, nil
		}
	}

	return nil, fmt.Errorf("pool not found in JSON")
}

func StorePoolData(pool *RaydiumAmmPool, filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read existing pools
	var pools []RaydiumAmmPool
	file, err := os.Open(filePath)
	if err == nil {
		if err := json.NewDecoder(file).Decode(&pools); err != nil && err.Error() != "EOF" {
			file.Close()
			return fmt.Errorf("failed to decode existing pools: %w", err)
		}
		file.Close()
	}

	// Check for existing pool and update or append
	found := false
	for i, p := range pools {
		if p.ID == pool.ID {
			pools[i] = *pool
			found = true
			log.Printf("Updated existing pool in storage: %s", pool.ID)
			break
		}
	}

	if !found {
		pools = append(pools, *pool)
		log.Printf("Added new pool to storage: %s", pool.ID)
	}

	// Write updated pools back to file
	file, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file for writing: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(pools); err != nil {
		return fmt.Errorf("failed to write pools to file: %w", err)
	}

	return nil
}
