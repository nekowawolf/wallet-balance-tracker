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

	switch choice {
	case "1":
		chain.Monad()
	default:
		fmt.Println("Invalid choice. Please select a valid option.")
		os.Exit(1)
	}
}