package raydium

import (
	"context"
	"fmt"
	"log"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/raydium/pool/amm"
	"corvus_bot/pkg/raydium/pool/clmm"
	"corvus_bot/pkg/utils"

	"github.com/gagliardetto/solana-go"
)

// Swap performs a token swap using Raydium pools (AMM or CLMM).
func Swap(
	ctx context.Context,
	wallet *solana.Wallet,
	poolID string,
	amountIn uint64,
	minAmountOut uint64,
	poolType string,
	cfg *config.Config,
) (solana.Signature, error) {
	if poolType != "AMM" && poolType != "CLMM" {
		return solana.Signature{}, fmt.Errorf("invalid pool type: %s (must be 'AMM' or 'CLMM')", poolType)
	}

	client := utils.GetRPCClient(cfg)
	var signature solana.Signature
	var err error

	switch poolType {
	case "AMM":
		pool, fetchErr := amm.FetchAmmPoolFromJSON(poolID, poolID, "./data/testdata/amm_pools.json")
		if fetchErr != nil {
			return solana.Signature{}, fmt.Errorf("failed to fetch AMM pool: %w", fetchErr)
		}

		signature, err = amm.SwapTokens(
			ctx,
			client,
			wallet,
			*pool,
			solana.MustPublicKeyFromBase58(pool.BaseMint),
			solana.MustPublicKeyFromBase58(pool.QuoteMint),
			amountIn,
			minAmountOut,
		)

	case "CLMM":
		pool, fetchErr := clmm.FetchClmmPoolFromJSON(poolID, "./data/testdata/clmm_pools.json")
		if fetchErr != nil {
			return solana.Signature{}, fmt.Errorf("failed to fetch CLMM pool: %w", fetchErr)
		}

		signature, err = clmm.SwapTokens(
			ctx,
			client,
			wallet,
			pool,
			amountIn,
			minAmountOut,
			cfg,
		)
	}

	if err != nil {
		return solana.Signature{}, fmt.Errorf("swap failed: %w", err)
	}

	log.Printf("Swap successful! Transaction signature: %s\n", signature)
	return signature, nil
}
