package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// SolanaClient handles interactions with the Solana blockchain.
type SolanaClient struct {
	RPCClient *rpc.Client
	WSClient  *ws.Client
}

// NewSolanaClient initializes the Solana client with RPC and WebSocket endpoints.
func NewSolanaClient(rpcURL, wsURL string) (*SolanaClient, error) {
	rpcClient := rpc.New(rpcURL)

	wsClient, err := ws.Connect(context.Background(), wsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Solana WebSocket: %w", err)
	}

	return &SolanaClient{
		RPCClient: rpcClient,
		WSClient:  wsClient,
	}, nil
}

// GetAccountBalance retrieves the balance of a Solana account in lamports.
func (c *SolanaClient) GetAccountBalance(address string) (uint64, error) {
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return 0, fmt.Errorf("invalid Solana address: %w", err)
	}

	balance, err := c.RPCClient.GetBalance(context.Background(), pubKey, rpc.CommitmentFinalized)
	if err != nil {
		return 0, fmt.Errorf("failed to get account balance: %w", err)
	}

	return balance.Value, nil
}

// SendTransaction signs and sends a transaction to the Solana blockchain.
func (c *SolanaClient) SendTransaction(tx *solana.Transaction, privateKey string) (string, error) {
	privKey, err := solana.PrivateKeyFromBase58(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	signatures, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key == privKey.PublicKey() {
			return &privKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	if len(signatures) == 0 {
		return "", fmt.Errorf("no signatures produced for transaction")
	}

	sig, err := c.RPCClient.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return sig.String(), nil
}
