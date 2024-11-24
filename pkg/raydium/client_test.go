package raydium

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchAndStoreAmmPool(t *testing.T) {
	client := NewRaydiumClient("../../testdata/amm_pools.json", "../../testdata/clmm_pools.json")
	tokenAddress := "2qEHjDLDLbuBgRYvsxhc5D6uDWAivNFZGan56P1tpump" // Test token

	// Fetch and store an AMM pool
	pool, err := client.FetchAndStorePoolData("amm", tokenAddress)
	assert.NoError(t, err, "Fetching AMM pool failed")
	assert.NotNil(t, pool, "AMM pool should not be nil")
}

func TestFetchAndStoreClmmPool(t *testing.T) {
	client := NewRaydiumClient("../../testdata/amm_pools.json", "../../testdata/clmm_pools.json")
	tokenAddress := "CLMMSampleToken" // Test token

	// Fetch and store a CLMM pool
	pool, err := client.FetchAndStorePoolData("clmm", tokenAddress)
	assert.NoError(t, err, "Fetching CLMM pool failed")
	assert.NotNil(t, pool, "CLMM pool should not be nil")
}
