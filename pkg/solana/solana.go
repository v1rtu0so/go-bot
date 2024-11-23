package solana

import (
	"context"
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
)

// SolanaClient is a wrapper around the Solana RPC client
type SolanaClient struct {
	Client *rpc.Client
}

// NewClient initializes a new Solana RPC client
func NewClient(rpcURL string) *SolanaClient {
	return &SolanaClient{
		Client: rpc.New(rpcURL),
	}
}

// LoadPrivateKey loads a private key from a file
func LoadPrivateKey(filePath string) (solana.PrivateKey, error) {
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return solana.PrivateKey{}, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := solana.PrivateKeyFromBase58(string(keyData))
	if err != nil {
		return solana.PrivateKey{}, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// GetRecentBlockhash retrieves the most recent blockhash from the Solana blockchain
func (sc *SolanaClient) GetRecentBlockhash(ctx context.Context) (*rpc.GetRecentBlockhashResult, error) {
	blockhash, err := sc.Client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent blockhash: %w", err)
	}

	return blockhash, nil
}

// SendTransaction sends a signed transaction to the Solana blockchain
func (sc *SolanaClient) SendTransaction(ctx context.Context, tx *solana.Transaction) (string, error) {
	txID, err := sc.Client.SendTransaction(ctx, tx)
	if err != nil {
		if rpcErr, ok := err.(*jsonrpc.RPCError); ok {
			return "", fmt.Errorf("rpc error: %s (code: %d)", rpcErr.Message, rpcErr.Code)
		}
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return txID, nil
}
