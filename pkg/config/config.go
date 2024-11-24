package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all the configuration items for the application.
type Config struct {
	RPCConnection        string `yaml:"rpc_connection"`          // Solana RPC endpoint
	WSConnection         string `yaml:"ws_connection"`           // Solana WebSocket endpoint
	RedisPassword        string `yaml:"redis_password"`          // Redis password for caching
	RedisHostname        string `yaml:"redis_hostname"`          // Redis host
	RedisPort            string `yaml:"redis_port"`              // Redis port
	PrivateKey           string `yaml:"private_key"`             // Private key for Solana wallet
	JitoBundleEndpoint   string `yaml:"jito_bundle_endpoint"`    // Jito block engine bundle endpoint
	JitoTxEndpoint       string `yaml:"jito_tx_endpoint"`        // Jito block engine transaction endpoint
	BloXrouteEndpoint    string `yaml:"bloxroute_endpoint"`      // BloXroute endpoint
	BloXrouteAuthHeader  string `yaml:"bloxroute_auth_header"`   // BloXroute authorization header
	BloXrouteTipAddress  string `yaml:"bloxroute_tip_address"`   // BloXroute tip address
	HeliusAPIKey         string `yaml:"helius_api_key"`          // Helius API key
	HeliusRPC            string `yaml:"helius_rpc"`              // Helius RPC endpoint
	HeliusStakedRPC      string `yaml:"helius_staked_rpc"`       // Helius staked RPC endpoint
	MainnetEndpoint      string `yaml:"mainnet_endpoint"`        // Mainnet endpoint
	LamportsPerSol       string `yaml:"lamports_per_sol"`        // Lamports per SOL const
	WSOLAddress          string `yaml:"wsol_address"`            // Wrapped SOL token address
	USDCAddress          string `yaml:"usdc_address"`            // USDC token address
	RaydiumAMMProgramID  string `yaml:"raydium_amm_program_id"`  // Raydium AMM Program ID
	RaydiumCLMMProgramID string `yaml:"raydium_clmm_program_id"` // Raydium CLMM Program ID
}

// LoadConfig reads the configuration from a YAML file.
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	cfg := &Config{}
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	return cfg, nil
}

// Validate checks if all required configuration fields are set.
func (cfg *Config) Validate() error {
	var missingFields []string

	if cfg.RPCConnection == "" {
		missingFields = append(missingFields, "rpc_connection")
	}
	if cfg.WSConnection == "" {
		missingFields = append(missingFields, "ws_connection")
	}
	if cfg.PrivateKey == "" {
		missingFields = append(missingFields, "private_key")
	}
	if cfg.BloXrouteAuthHeader == "" {
		missingFields = append(missingFields, "bloxroute_auth_header")
	}
	if cfg.HeliusAPIKey == "" {
		missingFields = append(missingFields, "helius_api_key")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required configuration fields: %v", missingFields)
	}

	return nil
}
