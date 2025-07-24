package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// 辅助函数：将 Wei 转换为 ETH (更易读)
func weiToEth(wei *big.Int) string {
	eth := new(big.Float).SetInt(wei)
	eth.Quo(eth, big.NewFloat(1e18))
	return fmt.Sprintf("%.6f", eth)
}

// 辅助函数：将 Wei 转换为 Gwei (Gas 价格常用)
func weiToGwei(wei *big.Int) string {
	gwei := new(big.Float).SetInt(wei)
	gwei.Quo(gwei, big.NewFloat(1e9))
	return fmt.Sprintf("%.2f", gwei)
}

func task01() {
	ctx := context.Background()

	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// 从环境变量获取配置
	sepoliaRPC := os.Getenv("SEPOLIA_RPC")
	if sepoliaRPC == "" {
		sepoliaRPC = "https://eth-sepolia.g.alchemy.com/v2/5kxZJaABVsl6R8LWJEcDvkapc6nwG8ik" // 默认值
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("PRIVATE_KEY environment variable is required")
	}

	recipientAddr := os.Getenv("RECIPIENT_ADDR")
	if recipientAddr == "" {
		log.Fatal("RECIPIENT_ADDR environment variable is required")
	}

	// connect to Sepolia network
	client, err := ethclient.DialContext(ctx, sepoliaRPC)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()
	fmt.Println("Connected to Sepolia network")

	// 首先检查连接是否正常
	latestBlock, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get latest block (connection issue): %v", err)
	}
	fmt.Printf("Latest Block Number: %d\n", latestBlock.Number().Uint64())

	// query block by number
	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatalf("Failed to retrieve block: %v", err)
	}
	fmt.Printf("Block Number: %d\n", block.Number().Uint64())
	fmt.Printf("Block Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Block Time: %d\n", block.Time())
	fmt.Printf("Block Transactions: %d\n", len(block.Transactions()))

	// prepare and send a transaction
	fmt.Println("\n=== Preparing Transaction ===")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("From Address: %s\n", fromAddress.Hex())
	fmt.Printf("To Address: %s\n", recipientAddr)

	// 检查账户余额
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Printf("Account Balance: %s ETH\n", weiToEth(balance))

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}
	fmt.Printf("Nonce: %d\n", nonce)
	value := big.NewInt(1e15) // 0.001 ETH
	gasLimit := uint64(21000) // standard gas limit for ETH transfer
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	fmt.Printf("Transfer Amount: %s ETH\n", weiToEth(value))
	fmt.Printf("Gas Price: %s Gwei\n", weiToGwei(gasPrice))
	fmt.Printf("Gas Limit: %d\n", gasLimit)

	// 计算总费用 (包括gas费)
	totalCost := new(big.Int).Add(value, new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit))))
	fmt.Printf("Total Cost (including gas): %s ETH\n", weiToEth(totalCost))

	// 检查余额是否足够
	if balance.Cmp(totalCost) < 0 {
		log.Fatalf("Insufficient balance! Need %s ETH but only have %s ETH",
			weiToEth(totalCost), weiToEth(balance))
	}

	toAddress := common.HexToAddress(recipientAddr)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Println("\n=== Transaction Sent Successfully ===")
	fmt.Printf("Transaction Hash: %s\n", signedTx.Hash().Hex())
	fmt.Printf("View on Etherscan: https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())
	fmt.Printf("From: %s\n", fromAddress.Hex())
	fmt.Printf("To: %s\n", toAddress.Hex())
	fmt.Printf("Amount: %s ETH\n", weiToEth(value))
	fmt.Printf("Gas Price: %s Gwei\n", weiToGwei(gasPrice))
	fmt.Println("\nNote: It may take 15-30 seconds for the transaction to be confirmed on the network.")
	fmt.Println("Check the Etherscan link above to monitor the transaction status.")
}
