package raydium

import (
	"context"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// FetchAMMPool retrieves the account data for a specific Raydium AMM Pool.
func (rc *RaydiumClient) FetchAMMPool(ctx context.Context, poolAddress string) (*rpc.GetAccountInfoResult, error) {
	log.Printf("Fetching AMM Pool account: %s", poolAddress)

	accountPubKey, err := solana.PublicKeyFromBase58(poolAddress)
	if err != nil {
		log.Printf("Invalid pool address: %v", err)
		return nil, err
	}

	accountInfo, err := rc.RPCClient.GetAccountInfo(ctx, accountPubKey)
	if err != nil {
		log.Printf("Error fetching AMM Pool account: %v", err)
		return nil, err
	}

	if accountInfo.Value == nil || accountInfo.Value.Data == nil {
		log.Println("No account data found.")
		return nil, nil
	}

	// Decode account data (this part will depend on the Raydium pool's data layout)
	log.Printf("Account Owner: %s", accountInfo.Value.Owner)
	log.Printf("Lamports: %d", accountInfo.Value.Lamports)
	log.Printf("Data Length: %d", len(accountInfo.Value.Data.GetBinary()))

	return accountInfo, nil
}
