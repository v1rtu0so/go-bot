package clmm

import (
	"fmt"
	"math"
)

// CalculateClmmTvl calculates the total value locked (TVL) in a CLMM pool.
// This is a placeholder implementation and should be expanded with real token price data.
func CalculateClmmTvl(vaultAReserve, vaultBReserve float64, priceA, priceB float64) (float64, error) {
	if priceA <= 0 || priceB <= 0 {
		return 0, fmt.Errorf("invalid token prices: priceA=%f, priceB=%f", priceA, priceB)
	}

	vaultAValue := vaultAReserve * priceA
	vaultBValue := vaultBReserve * priceB

	return vaultAValue + vaultBValue, nil
}

// CalculateClmmFeeRate calculates the fee rate for a given CLMM pool.
func CalculateClmmFeeRate(protocolFeeRate, tradeFeeRate float64) float64 {
	return math.Max(protocolFeeRate, tradeFeeRate)
}

// FormatClmmPoolData formats the key data of a CLMM pool for display or logging purposes.
func FormatClmmPoolData(pool RaydiumClmmPool) string {
	return fmt.Sprintf(
		"Pool ID: %s\nMint A: %s\nMint B: %s\nProgram ID A: %s\nProgram ID B: %s\n",
		pool.ID, pool.MintA, pool.MintB, pool.MintProgramIDA, pool.MintProgramIDB,
	)
}
