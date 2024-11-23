package main

import (
	"context"
	"log"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/raydium"
)

func main() {
	log.Println("Starting Corvus Bot...")

	// Load configuration
	cfg, err := config.LoadConfig("pkg/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Raydium client
	raydiumClient, err := raydium.NewRaydiumClient(cfg.RPCConnection, cfg.RaydiumAMMProgramID)
	if err != nil {
		log.Fatalf("Failed to initialize Raydium client: %v", err)
	}

	// Fetch AMM Pool data
	ctx := context.Background()
	poolAddress := "BVdxeDmh4YWgzR3v4tnetwCnU12zPdcXwooAspiiUGaB" // Replace with an actual pool address
	poolData, err := raydiumClient.FetchAMMPool(ctx, poolAddress)
	if err != nil {
		log.Fatalf("Failed to fetch AMM Pool data: %v", err)
	}
	log.Printf("Fetched AMM Pool data: %+v", poolData)

	// Example: Parse a transaction
	//txSignature := "fnzGQb3gegHU5gNt15kfwMoEJDuH8HDurj8k3yR8scEXmC2qDgQ4gcp5Y57nAXot658Nj5xzMmG9c6Kmmkz2xCo" // Replace with an actual transaction signature
	//err = raydium.ParseTransaction(ctx, cfg.RPCConnection, txSignature)
	//if err != nil {
	//	log.Fatalf("Failed to parse transaction: %v", err)
	//}

	log.Println("Corvus Bot completed successfully!")
}
