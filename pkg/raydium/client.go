package raydium

import (
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// RaydiumClient handles interactions with Raydium's AMM program.
type RaydiumClient struct {
	RPCClient *rpc.Client
	ProgramID solana.PublicKey
}

// NewRaydiumClient initializes a new Raydium client.
func NewRaydiumClient(rpcURL, programID string) (*RaydiumClient, error) {
	programPubKey, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		log.Printf("Invalid program ID: %v", err)
		return nil, err
	}

	return &RaydiumClient{
		RPCClient: rpc.New(rpcURL),
		ProgramID: programPubKey,
	}, nil
}
