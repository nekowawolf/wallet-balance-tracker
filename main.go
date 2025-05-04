package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nekowawolf/wallet-balance-tracker/chain"
)

func main() {
	fmt.Println("\nSelect chain:")
	fmt.Println("1. Monad")
	fmt.Print("Enter your choice: ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice != "1" {
		fmt.Println("Invalid choice. Currently only Monad is supported.")
		os.Exit(1)
	}

	fmt.Println("\nCheck token:")
	fmt.Println("1. Native (MON)")
	fmt.Println("2. Input token address manually")
	fmt.Print("Enter your choice: ")

	tokenChoice, _ := reader.ReadString('\n')
	tokenChoice = strings.TrimSpace(tokenChoice)

	switch tokenChoice {
	case "1":
		fmt.Println("\nChecking native balances for all configured addresses...")
		chain.CheckMonadNativeBalances()
	case "2":
		fmt.Print("\nEnter token contract address to check: ")
		tokenAddress, _ := reader.ReadString('\n')
		tokenAddress = strings.TrimSpace(tokenAddress)
		if tokenAddress != "" {
			chain.CheckCustomTokenBalances(tokenAddress)
		}
	default:
		fmt.Println("Invalid choice")
		os.Exit(1)
	}
}