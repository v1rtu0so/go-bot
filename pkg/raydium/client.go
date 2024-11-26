package raydium

import (
	"context"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"corvus_bot/pkg/raydium/pool/amm"
	"corvus_bot/pkg/raydium/pool/clmm"
)

// RaydiumClient is the main entry point for interacting with Raydium pools.
type RaydiumClient struct {
	RPCConnection string
	AMMProgramID  string
	CLMMProgramID string
	AMMDataPath   string
	CLMMDataPath  string
}

// NewRaydiumClient creates a new RaydiumClient instance.
func NewRaydiumClient(rpcConnection, ammProgramID, clmmProgramID, ammDataPath, clmmDataPath string) *RaydiumClient {
	return &RaydiumClient{
		RPCConnection: rpcConnection,
		AMMProgramID:  ammProgramID,
		CLMMProgramID: clmmProgramID,
		AMMDataPath:   ammDataPath,
		CLMMDataPath:  clmmDataPath,
	}
}

// FetchPoolData fetches the pool data for a given token address and pool type.
// It first attempts to fetch from JSON storage, then falls back to network fetch if needed.
func (rc *RaydiumClient) FetchPoolData(ctx context.Context, poolType, tokenAddress string) (interface{}, error) {
	client := rpc.New(rc.RPCConnection)
	wsolAddress := "So11111111111111111111111111111111111111112"

	switch poolType {
	case "AMM":
		log.Printf("Attempting to fetch AMM pool data for token: %s", tokenAddress)
		// Try with token as base, WSOL as quote
		pool, err := amm.FetchAmmPoolFromJSONOrNetwork(
			ctx,
			client,
			tokenAddress,
			wsolAddress,
			rc.AMMProgramID,
			rc.AMMDataPath,
		)
		if err == nil {
			return pool, nil
		}

		// If failed, try with WSOL as base, token as quote
		pool, err = amm.FetchAmmPoolFromJSONOrNetwork(
			ctx,
			client,
			wsolAddress,
			tokenAddress,
			rc.AMMProgramID,
			rc.AMMDataPath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch AMM pool: %w", err)
		}
		return pool, nil

	case "CLMM":
		log.Printf("Attempting to fetch CLMM pool data for token: %s", tokenAddress)
		pool, err := clmm.FetchClmmPoolFromJSONOrNetwork(
			ctx,
			client,
			tokenAddress,
			rc.CLMMProgramID,
			rc.CLMMDataPath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch CLMM pool: %w", err)
		}
		return pool, nil

	default:
		return nil, fmt.Errorf("unsupported pool type: %s", poolType)
	}
}

// PerformSwap executes a token swap using the specified pool type and data.
func (rc *RaydiumClient) PerformSwap(
	ctx context.Context,
	wallet *solana.Wallet,
	poolType string,
	poolData interface{},
	amountIn uint64,
	minAmountOut uint64,
) (solana.Signature, error) {
	client := rpc.New(rc.RPCConnection)

	switch poolType {
	case "AMM":
		ammPool, ok := poolData.(*amm.RaydiumAmmPool)
		if !ok {
			return solana.Signature{}, fmt.Errorf("invalid pool data for AMM pool")
		}
		return amm.SwapTokens(
			ctx,
			client,
			wallet,
			*ammPool,
			solana.MustPublicKeyFromBase58(ammPool.BaseMint),
			solana.MustPublicKeyFromBase58(ammPool.QuoteMint),
			amountIn,
			minAmountOut,
		)

	case "CLMM":
		clmmPool, ok := poolData.(*clmm.RaydiumClmmPool)
		if !ok {
			return solana.Signature{}, fmt.Errorf("invalid pool data for CLMM pool")
		}
		return clmm.SwapTokens(
			ctx,
			client,
			wallet,
			clmmPool,
			amountIn,
			minAmountOut,
			nil,
		)

	default:
		return solana.Signature{}, fmt.Errorf("unsupported pool type: %s", poolType)
	}
}

// ValidateAndPerformSwap orchestrates the entire swap process.
func (rc *RaydiumClient) ValidateAndPerformSwap(
	ctx context.Context,
	wallet *solana.Wallet,
	poolType, tokenAddress string,
	amountIn, minAmountOut uint64,
) (solana.Signature, error) {
	poolData, err := rc.FetchPoolData(ctx, poolType, tokenAddress)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to fetch pool data: %w", err)
	}

	return rc.PerformSwap(ctx, wallet, poolType, poolData, amountIn, minAmountOut)
}
