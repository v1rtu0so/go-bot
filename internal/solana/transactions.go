package solana

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// SignTransaction signs the transaction with the provided private keys.
func SignTransaction(tx *Transaction, privateKeys [][]byte) error {
	serializedTx, err := tx.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize transaction: %w", err)
	}

	for _, key := range privateKeys {
		signature := ed25519.Sign(key, serializedTx)
		tx.Signatures = append(tx.Signatures, signature)
	}

	return nil
}

// SendTransaction sends a signed transaction to the Solana network.
func SendTransaction(rpcURL string, tx *Transaction) (string, error) {
	serializedTx, err := tx.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	encodedTx := base64.StdEncoding.EncodeToString(serializedTx)

	params := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "sendTransaction",
		"params":  []interface{}{encodedTx},
	}

	data, _ := json.Marshal(params)
	resp, err := http.Post(rpcURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Result, nil
}

// BuildRawTransaction creates a raw transaction with instructions.
func BuildRawTransaction(instructions []Instruction, blockhash string, payerPubKey []byte) RawTransaction {
	return RawTransaction{
		Blockhash:    blockhash,
		PayerPubKey:  payerPubKey,
		Instructions: instructions,
	}
}

// SendRawTransaction sends the signed transaction to the Solana RPC endpoint.
func SendRawTransaction(ctx context.Context, rpcURL string, signedTx SignedTransaction) (string, error) {
	encodedTx := base64.StdEncoding.EncodeToString(signedTx.SerializedTx)
	params := []interface{}{encodedTx}
	var response struct {
		Result string `json:"result"`
	}
	err := RPCRequest(ctx, rpcURL, "sendTransaction", params, &response)
	if err != nil {
		return "", err
	}
	return response.Result, nil
}
