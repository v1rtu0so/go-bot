package raydium

import (
	"corvus_bot/pkg/raydium/pool/amm"
	"corvus_bot/pkg/raydium/pool/clmm"
	"fmt"
)

type RaydiumClient struct {
	ammFilePath  string
	clmmFilePath string
}

// NewRaydiumClient creates a new Raydium client with file paths for AMM and CLMM data.
func NewRaydiumClient(ammFilePath, clmmFilePath string) *RaydiumClient {
	return &RaydiumClient{
		ammFilePath:  ammFilePath,
		clmmFilePath: clmmFilePath,
	}
}

// FetchAndStorePoolData fetches and stores pool data based on the type and token address.
func (c *RaydiumClient) FetchAndStorePoolData(poolType, tokenAddress string) (interface{}, error) {
	switch poolType {
	case "amm":
		return amm.FetchAndStorePool(tokenAddress, c.ammFilePath)
	case "clmm":
		return clmm.FetchAndStorePool(tokenAddress, c.clmmFilePath)
	default:
		return nil, fmt.Errorf("invalid pool type: %s", poolType)
	}
}
