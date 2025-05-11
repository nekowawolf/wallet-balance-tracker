package chain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"math"
	"bufio"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

const (
	RPC_URL_MONAD  = "https://testnet-rpc.monad.xyz"
	CHAIN_ID_MONAD = 10143
	MAX_RECIPIENTS = 50
	DELAY_SECONDS  = 2
)

var (
	cyan    = color.New(color.FgCyan).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
)

const erc20ABI = `[
	{
		"constant": true,
		"inputs": [{"name": "_owner", "type": "address"}],
		"name": "balanceOf",
		"outputs": [{"name": "balance", "type": "uint256"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [{"name": "", "type": "string"}],
		"type": "function"
	}
]`

func loadMonadClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(RPC_URL_MONAD)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Monad RPC: %v", err)
	}
	return client, nil
}

func Monad() {
	fmt.Println("\n" + cyan("Check token:"))
	fmt.Println(green("1. Native (MON)"))
	fmt.Println(green("2. Input token address manually"))
	fmt.Print(cyan("Enter your choice: "))

	reader := bufio.NewReader(os.Stdin)
	tokenChoice, _ := reader.ReadString('\n')
	tokenChoice = strings.TrimSpace(tokenChoice)

	switch tokenChoice {
	case "1":
		fmt.Println("\n" + yellow("Checking native balances for all configured addresses..."))
		CheckMonadNativeBalances()
	case "2":
		fmt.Print("\n" + cyan("Enter token contract address to check: "))
		tokenAddress, _ := reader.ReadString('\n')
		tokenAddress = strings.TrimSpace(tokenAddress)
		if tokenAddress != "" {
			CheckCustomTokenBalances(tokenAddress)
		} else {
			fmt.Println(red("Token address cannot be empty."))
			os.Exit(1)
		}
	default:
		fmt.Println(red("Invalid choice"))
		os.Exit(1)
	}
}

func getaddress() []common.Address {
	var address []common.Address
	for i := 1; i <= MAX_RECIPIENTS; i++ {
		envVar := fmt.Sprintf("ADDRESS%d", i)
		if addr := os.Getenv(envVar); addr != "" {
			cleanAddr := strings.Trim(addr, `" `)
			if common.IsHexAddress(cleanAddr) {
				address = append(address, common.HexToAddress(cleanAddr))
			}
		}
	}
	return address
}

func CheckMonadNativeBalances() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(red("Error loading .env file: %v"), err)
	}

	client, err := loadMonadClient()
	if err != nil {
		log.Fatal(red(err))
	}
	defer client.Close()

	address := getaddress()
	if len(address) == 0 {
		log.Fatal(red("No valid address found in .env"))
	}

	for i, addr := range address {
		balance, err := client.BalanceAt(context.Background(), addr, nil)
		if err != nil {
			log.Printf(red("Error checking balance for address %s: %v"), addr.Hex(), err)
			continue
		}

		fmt.Printf("\n%s #%d: %s\n", cyan("[Wallet]"), i+1, addr.Hex())
		fmt.Printf("%s: %s MON\n", magenta("Balance"), green(fmt.Sprintf("%.4f", weiToDecimal(balance, 18))))
		fmt.Println("\n▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔")

		if i < len(address)-1 {
			time.Sleep(DELAY_SECONDS * time.Second)
		}
	}

	fmt.Println(green("\n✅ CHECKED ADDRESS SUCCESS"))
	fmt.Println("\nFollow X : 0xNekowawolf\n")
}

func CheckCustomTokenBalances(tokenAddress string) {
	if !common.IsHexAddress(tokenAddress) {
		log.Fatal(red("Invalid token contract address format"))
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf(red("Error loading .env file: %v"), err)
	}

	client, err := loadMonadClient()
	if err != nil {
		log.Fatal(red(err))
	}
	defer client.Close()

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		log.Fatalf(red("Failed to parse ABI: %v"), err)
	}

	tokenName, err := getTokenName(client, common.HexToAddress(tokenAddress), parsedABI)
	if err != nil {
		log.Printf(yellow("Warning: Could not get token name (%v), using contract address instead"), err)
		tokenName = tokenAddress
	}

	address := getaddress()
	if len(address) == 0 {
		log.Fatal(red("No valid address found in .env"))
	}

	for i, addr := range address {
		balance, err := getTokenBalance(client, common.HexToAddress(tokenAddress), addr, parsedABI)
		if err != nil {
			log.Printf(red("Error checking token balance for address %s: %v"), addr.Hex(), err)
			continue
		}

		fmt.Printf("\n%s #%d: %s\n", cyan("[Wallet]"), i+1, addr.Hex())
		fmt.Printf("%s: %s %s\n", magenta("Balance"), green(fmt.Sprintf("%.4f", weiToDecimal(balance, 18))), tokenName)
		fmt.Println("\n▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔")

		if i < len(address)-1 {
			time.Sleep(DELAY_SECONDS * time.Second)
		}
	}

	fmt.Println(green("\n✅ CHECKED ADDRESS SUCCESS"))
	fmt.Println("\nFollow X : 0xNekowawolf\n")
}

func getTokenName(client *ethclient.Client, tokenAddress common.Address, parsedABI abi.ABI) (string, error) {
	data, err := parsedABI.Pack("name")
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	var name string
	err = parsedABI.UnpackIntoInterface(&name, "name", result)
	if err != nil {
		return "", err
	}

	return name, nil
}

func getTokenBalance(client *ethclient.Client, tokenAddress, walletAddress common.Address, parsedABI abi.ABI) (*big.Int, error) {
	data, err := parsedABI.Pack("balanceOf", walletAddress)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var balance *big.Int
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func weiToDecimal(wei *big.Int, decimals int) float64 {
    weiFloat := new(big.Float).SetInt(wei)
    divisor := new(big.Float).SetFloat64(math.Pow10(decimals))
    result := new(big.Float).Quo(weiFloat, divisor)
    floatVal, _ := result.Float64()
    return floatVal
}