package solana

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gagliardetto/solana-go"
)

// RawTransaction represents a raw Solana transaction.
type RawTransaction struct {
	Blockhash    string
	PayerPubKey  []byte
	Instructions []Instruction
}

// SignedTransaction represents a signed Solana transaction.
type SignedTransaction struct {
	SerializedTx []byte
	Signatures   [][]byte
}

// CreateTransaction creates a new Solana transaction with instructions.
func CreateTransaction(recentBlockhash solana.Hash, instructions []*solana.Instruction) (*solana.Transaction, error) {
	// Convert []*solana.Instruction to []solana.Instruction
	var insts []solana.Instruction
	for _, inst := range instructions {
		insts = append(insts, *inst)
	}

	tx, err := solana.NewTransaction(
		insts,
		recentBlockhash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return tx, nil
}

// GeneratePrivateKey creates a random private key (Ed25519 format).
func GeneratePrivateKey() []byte {
	privateKey := make([]byte, 64)
	rand.Read(privateKey)
	return privateKey
}

// DerivePublicKey derives a public key from the given private key (Ed25519 format).
func DerivePublicKey(privateKey []byte) []byte {
	return privateKey[32:] // Ed25519: public key is in the last 32 bytes.
}

// GetRecentBlockhash fetches the recent blockhash from the RPC node.
func GetRecentBlockhash(ctx context.Context, rpcURL string) (string, error) {
	response := struct {
		Result struct {
			Blockhash string `json:"blockhash"`
		} `json:"result"`
	}{}

	err := RPCRequest(ctx, rpcURL, "getRecentBlockhash", []interface{}{}, &response)
	if err != nil {
		return "", err
	}

	return response.Result.Blockhash, nil
}

// RPCRequest sends a JSON-RPC request to the Solana RPC endpoint.
func RPCRequest(ctx context.Context, rpcURL, method string, params interface{}, result interface{}) error {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("RPC error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected RPC response")
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
