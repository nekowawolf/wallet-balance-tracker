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
	RPC_URL_MEGAETH  = "https://carrot.megaeth.com/rpc"
	CHAIN_ID_MEGAETH = 6342
	MAX_RECIPIENTS_MEGAETH   = 50
	DELAY_SECONDS_MEGAETH    = 2
)

var (
	cyan1    = color.New(color.FgCyan).SprintFunc()
	yellow1  = color.New(color.FgYellow).SprintFunc()
	green1   = color.New(color.FgGreen).SprintFunc()
	red1     = color.New(color.FgRed).SprintFunc()
	magenta1 = color.New(color.FgMagenta).SprintFunc()
	blue1    = color.New(color.FgBlue).SprintFunc()
)

const erc20ABIMegaETH = `[
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

func loadMegaETHClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(RPC_URL_MEGAETH)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MegaETH RPC: %v", err)
	}
	return client, nil
}

func MegaETH() {
	fmt.Println("\n" + cyan1("Check token:"))
	fmt.Println(green1("1. Native (ETH)"))
	fmt.Println(green1("2. Input token address manually"))
	fmt.Print(cyan1("Enter your choice: "))

	reader := bufio.NewReader(os.Stdin)
	tokenChoice, _ := reader.ReadString('\n')
	tokenChoice = strings.TrimSpace(tokenChoice)

	switch tokenChoice {
	case "1":
		fmt.Println("\n" + yellow1("Checking native balances for all configured1 addresses..."))
		CheckMegaETHNativeBalances()
	case "2":
		fmt.Print("\n" + cyan1("Enter token contract address to check: "))
		tokenAddress, _ := reader.ReadString('\n')
		tokenAddress = strings.TrimSpace(tokenAddress)
		if tokenAddress != "" {
			CheckMegaETHTokenBalances(tokenAddress)
		} else {
			fmt.Println(red1("Token address cannot be empty."))
			os.Exit(1)
		}
	default:
		fmt.Println(red1("Invalid choice"))
		os.Exit(1)
	}
}

func CheckMegaETHNativeBalances() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(red1("Error loading .env file: %v"), err)
	}

	client, err := loadMegaETHClient()
	if err != nil {
		log.Fatal(red1(err))
	}
	defer client.Close()

	address := getaddress()
	if len(address) == 0 {
		log.Fatal(red1("No valid address found in .env"))
	}

	for i, addr := range address {
		balance, err := client.BalanceAt(context.Background(), addr, nil)
		if err != nil {
			log.Printf(red1("Error checking balance for address %s: %v"), addr.Hex(), err)
			continue
		}

		fmt.Printf("\n%s #%d: %s\n", cyan1("[Wallet]"), i+1, addr.Hex())
		fmt.Printf("%s: %s ETH\n", magenta1("Balance"), green1(fmt.Sprintf("%.4f", weiToDecimalMegaETH(balance, 18))))
		fmt.Println("\n▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔")

		if i < len(address)-1 {
			time.Sleep(DELAY_SECONDS * time.Second)
		}
	}

	fmt.Println(green1("\n✅ CHECKED ADDRESS SUCCESS"))
	fmt.Println("\nFollow X : 0xNekowawolf\n")
}

func CheckMegaETHTokenBalances(tokenAddress string) {
	if !common.IsHexAddress(tokenAddress) {
		log.Fatal(red1("Invalid token contract address format"))
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf(red1("Error loading .env file: %v"), err)
	}

	client, err := loadMegaETHClient()
	if err != nil {
		log.Fatal(red1(err))
	}
	defer client.Close()

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABIMegaETH))
	if err != nil {
		log.Fatalf(red1("Failed to parse ABI: %v"), err)
	}

	tokenName, err := getTokenNameMegaETH(client, common.HexToAddress(tokenAddress), parsedABI)
	if err != nil {
		log.Printf(yellow1("Warning: Could not get token name (%v), using contract address instead"), err)
		tokenName = tokenAddress
	}

	address := getaddress()
	if len(address) == 0 {
		log.Fatal(red1("No valid address found in .env"))
	}

	for i, addr := range address {
		balance, err := getTokenBalanceMegaETH(client, common.HexToAddress(tokenAddress), addr, parsedABI)
		if err != nil {
			log.Printf(red1("Error checking token balance for address %s: %v"), addr.Hex(), err)
			continue
		}

		fmt.Printf("\n%s #%d: %s\n", cyan1("[Wallet]"), i+1, addr.Hex())
		fmt.Printf("%s: %s %s\n", magenta1("Balance"), green1(fmt.Sprintf("%.4f", weiToDecimalMegaETH(balance, 18))), tokenName)
		fmt.Println("\n▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔")

		if i < len(address)-1 {
			time.Sleep(DELAY_SECONDS * time.Second)
		}
	}

	fmt.Println(green1("\n✅ CHECKED ADDRESS SUCCESS"))
	fmt.Println("\nFollow X : 0xNekowawolf\n")
}

func getTokenNameMegaETH(client *ethclient.Client, tokenAddress common.Address, parsedABI abi.ABI) (string, error) {
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

func getTokenBalanceMegaETH(client *ethclient.Client, tokenAddress, walletAddress common.Address, parsedABI abi.ABI) (*big.Int, error) {
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

func weiToDecimalMegaETH(wei *big.Int, decimals int) float64 {
    weiFloat := new(big.Float).SetInt(wei)
    divisor := new(big.Float).SetFloat64(math.Pow10(decimals))
    result := new(big.Float).Quo(weiFloat, divisor)
    floatVal, _ := result.Float64()
    return floatVal
}