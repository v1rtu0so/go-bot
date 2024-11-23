package main

import (
	"fmt"
	"os"
	"strconv"

	"corvus_bot/internal/trading"
	"corvus_bot/pkg/config"

	"github.com/spf13/cobra"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("pkg/config/config.yaml") // Corrected to reference the module's structure
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Root command
	var rootCmd = &cobra.Command{
		Use:   "corvus_bot",
		Short: "A Solana trading bot",
	}

	// Buy command
	var buyCmd = &cobra.Command{
		Use:   "buy [dex] [address] [amount] [fee] [slippage]",
		Short: "Execute a buy operation on a specific DEX",
		Args:  cobra.ExactArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			dex := args[0]
			address := args[1]
			amount, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				fmt.Printf("Invalid amount: %v\n", err)
				os.Exit(1)
			}
			fee, err := strconv.ParseFloat(args[3], 64)
			if err != nil {
				fmt.Printf("Invalid fee: %v\n", err)
				os.Exit(1)
			}
			slippage, err := strconv.ParseFloat(args[4], 64)
			if err != nil {
				fmt.Printf("Invalid slippage: %v\n", err)
				os.Exit(1)
			}

			err = trading.Buy(cfg, dex, address, amount, fee, slippage)
			if err != nil {
				fmt.Printf("Error executing buy: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Buy operation completed successfully!")
		},
	}

	// Sell command
	var sellCmd = &cobra.Command{
		Use:   "sell [dex] [address] [amount] [fee] [slippage]",
		Short: "Execute a sell operation on a specific DEX",
		Args:  cobra.ExactArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			dex := args[0]
			address := args[1]
			amount, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				fmt.Printf("Invalid amount: %v\n", err)
				os.Exit(1)
			}
			fee, err := strconv.ParseFloat(args[3], 64)
			if err != nil {
				fmt.Printf("Invalid fee: %v\n", err)
				os.Exit(1)
			}
			slippage, err := strconv.ParseFloat(args[4], 64)
			if err != nil {
				fmt.Printf("Invalid slippage: %v\n", err)
				os.Exit(1)
			}

			err = trading.Sell(cfg, dex, address, amount, fee, slippage)
			if err != nil {
				fmt.Printf("Error executing sell: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Sell operation completed successfully!")
		},
	}

	// Add subcommands
	rootCmd.AddCommand(buyCmd)
	rootCmd.AddCommand(sellCmd)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
