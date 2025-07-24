package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/local/go-eth-demo/go-eth-demo/counter"
)

func task02() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
	// 从环境变量获取配置
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth-sepolia.g.alchemy.com/v2/5kxZJaABVsl6R8LWJEcDvkapc6nwG8ik" // 默认值
	}
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("PRIVATE_KEY environment variable is required")
	}
	recipientAddr := os.Getenv("RECIPIENT_ADDR")
	if recipientAddr == "" {
		log.Fatal("RECIPIENT_ADDR environment variable is required")
	}
	contractAddr := os.Getenv("CONTRACT_ADDR")
	if contractAddr == "" {
		log.Fatal("CONTRACT_ADDR environment variable is required")
	}
	// 连接到以太坊客户端
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()
	log.Println("Connected to Sepolia successfully")
	// 加载私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	log.Println("Private key loaded successfully")
	// 获取网络 ID
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}
	log.Printf("Connected to Sepolia network: %s", chainID.String())
	log.Println("Recipient address:", recipientAddr)
	log.Println("Contract address:", contractAddr)
	// 创建授权的交易发送者
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	log.Println("Authorized transactor created successfully")
	// 创建合约实例
	address := common.HexToAddress(contractAddr)
	contract, err := counter.NewCounter(address, client)
	if err != nil {
		log.Fatalf("Failed to create contract instance: %v", err)
	}
	log.Println("Contract instance created successfully")

	// 先查询当前值（交易前）
	countBefore, err := contract.GetCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("Failed to get counter value before increment: %v", err)
	}
	log.Printf("Counter value BEFORE increment: %d", countBefore)
	// 发送交易以递增计数器
	tx, err := contract.Increment(auth)
	if err != nil {
		log.Fatalf("Failed to increment counter: %v", err)
	}
	log.Printf("Counter increment transaction sent: %s", tx.Hash().Hex())
	log.Println("Waiting for transaction to be confirmed...")

	// 等待交易确认
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for transaction confirmation: %v", err)
	}

	if receipt.Status == 1 {
		log.Printf("Transaction confirmed successfully in block: %d", receipt.BlockNumber.Uint64())
		log.Printf("Gas used: %d", receipt.GasUsed)
	} else {
		log.Fatalf("Transaction failed with status: %d", receipt.Status)
	}

	// 等待一点时间让状态同步
	log.Println("Waiting for state synchronization...")
	time.Sleep(2 * time.Second)

	// 现在查询计数器值（交易已确认）
	count, err := contract.GetCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("Failed to get counter value: %v", err)
	}
	log.Printf("Current counter value after confirmation: %d", count)

	// 验证是否真的递增了
	if count.Cmp(countBefore) > 0 {
		log.Printf("✅ SUCCESS: Counter incremented from %d to %d", countBefore, count)
	} else {
		log.Printf("❌ WARNING: Counter did not increment! Before: %d, After: %d", countBefore, count)
		log.Printf("Check transaction details on Etherscan: https://sepolia.etherscan.io/tx/%s", tx.Hash().Hex())

		// 再次查询，使用最新区块
		log.Println("Retrying query with latest block...")
		time.Sleep(1 * time.Second)
		countRetry, err := contract.GetCount(&bind.CallOpts{
			Context:     ctx,
			BlockNumber: nil, // 使用最新区块
		})
		if err != nil {
			log.Printf("Retry query failed: %v", err)
		} else {
			log.Printf("Retry result: %d", countRetry)
		}
	}
}
