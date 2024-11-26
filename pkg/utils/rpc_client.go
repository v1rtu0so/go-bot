package utils

import (
	"context"
	"corvus_bot/pkg/config"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// RPCClientInterface defines the required methods for an RPC client.
type RPCClientInterface interface {
	GetLatestBlockhash(ctx context.Context, commitment rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error)
	SendTransaction(ctx context.Context, tx *solana.Transaction) (solana.Signature, error)
}

// RealRPCClient is the real implementation of the RPCClientInterface.
type RealRPCClient struct {
	Client *rpc.Client
}

// GetLatestBlockhash fetches the latest blockhash.
func (r *RealRPCClient) GetLatestBlockhash(ctx context.Context, commitment rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error) {
	return r.Client.GetLatestBlockhash(ctx, commitment)
}

// SendTransaction sends a Solana transaction.
func (r *RealRPCClient) SendTransaction(ctx context.Context, tx *solana.Transaction) (solana.Signature, error) {
	return r.Client.SendTransaction(ctx, tx)
}

// MockRPCClient is a mock implementation of the RPCClientInterface for testing.
type MockRPCClient struct {
	MockGetLatestBlockhash func(ctx context.Context, commitment rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error)
	MockSendTransaction    func(ctx context.Context, tx *solana.Transaction) (solana.Signature, error)
}

// GetLatestBlockhash mocks the GetLatestBlockhash method.
func (m *MockRPCClient) GetLatestBlockhash(ctx context.Context, commitment rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error) {
	if m.MockGetLatestBlockhash != nil {
		return m.MockGetLatestBlockhash(ctx, commitment)
	}
	return nil, nil
}

// SendTransaction mocks the SendTransaction method.
func (m *MockRPCClient) SendTransaction(ctx context.Context, tx *solana.Transaction) (solana.Signature, error) {
	if m.MockSendTransaction != nil {
		return m.MockSendTransaction(ctx, tx)
	}
	// Return a valid solana.Signature instead of an empty string.
	return solana.Signature{}, nil
}

var mockRPCClient *MockRPCClient

// SetMockRPCClient allows test code to set a mock RPC client globally.
func SetMockRPCClient(mockClient *MockRPCClient) {
	mockRPCClient = mockClient
}

// GetRPCClient returns the appropriate RPC client based on the testing environment.
func GetRPCClient(cfg *config.Config) RPCClientInterface {
	if cfg.Testing {
		if mockRPCClient == nil {
			panic("MockRPCClient is not set. Call SetMockRPCClient before using GetRPCClient in testing mode.")
		}
		return mockRPCClient
	}
	return &RealRPCClient{Client: rpc.New(cfg.RPCConnection)}
}
