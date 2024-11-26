package amm

import (
	"fmt"
	"math"
)

// CalculatePoolTvl calculates the total value locked (TVL) in the AMM pool
func CalculatePoolTvl(baseReserve, quoteReserve float64, basePrice, quotePrice float64) (float64, error) {
	if basePrice <= 0 || quotePrice <= 0 {
		return 0, fmt.Errorf("invalid token prices: basePrice=%f, quotePrice=%f", basePrice, quotePrice)
	}

	baseValue := baseReserve * basePrice
	quoteValue := quoteReserve * quotePrice

	return baseValue + quoteValue, nil
}

// CalculatePoolFeeRate calculates the fee rate for a given AMM pool
func CalculatePoolFeeRate(protocolFeeRate, tradeFeeRate float64) float64 {
	return math.Max(protocolFeeRate, tradeFeeRate)
}

// FormatPoolData formats the key data of an AMM pool for display or logging
func FormatPoolData(pool RaydiumAmmPool) string {
	return fmt.Sprintf(
		"Pool ID: %s\nBase Mint: %s\nQuote Mint: %s\nVersion: %d\nProgram ID: %s\n",
		pool.ID, pool.BaseMint, pool.QuoteMint, pool.Version, pool.ProgramID,
	)
}
