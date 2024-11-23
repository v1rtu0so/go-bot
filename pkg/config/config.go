package config

import (
	"log"
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
	BloxrouteEndpoint    string `yaml:"bloxroute_endpoint"`      // BloXroute endpoint
	BloxrouteAuthHeader  string `yaml:"bloxroute_auth_header"`   // BloXroute authorization header
	BloxrouteTipAddress  string `yaml:"bloxroute_tip_address"`   // BloXroute tip address
	HeliusAPIKey         string `yaml:"helius_api_key"`          // Helius API key
	HeliusRPC            string `yaml:"helius_rpc"`              // Helius RPC endpoint
	HeliusStakedRPC      string `yaml:"helius_staked_rpc"`       // Helius staked RPC endpoint
	MainnetEndpoint      string `yaml:"mainnet_endpoint"`        // Mainnet endpoint
	WSOLAddress          string `yaml:"wsol_address"`            // Wrapped SOL token address
	USDCAddress          string `yaml:"usdc_address"`            // USDC token address
	RaydiumAMMProgramID  string `yaml:"raydium_amm_program_id"`  // Raydium AMM Program ID
	RaydiumCLMMProgramID string `yaml:"raydium_clmm_program_id"` // Raydium CLMM Program ID
}

// LoadConfig reads the configuration from a YAML file.
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	cfg := &Config{}
	if err := decoder.Decode(cfg); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
		return nil, err
	}

	// Log a warning for empty sensitive fields.
	if cfg.PrivateKey == "" {
		log.Println("WARNING: Private key is not set in the configuration.")
	}
	if cfg.BloxrouteAuthHeader == "" {
		log.Println("WARNING: BloXroute auth header is not set in the configuration.")
	}
	if cfg.HeliusAPIKey == "" {
		log.Println("WARNING: Helius API key is not set in the configuration.")
	}

	return cfg, nil
}
