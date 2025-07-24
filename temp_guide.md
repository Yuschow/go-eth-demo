package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// 演示多笔交易的正确处理方式
func demonstrateMultipleTransactions() {
	fmt.Println("=== 多笔交易的正确处理方式 ===\n")

	// 模拟参数
	ctx := context.Background()
	// client := ethclient.DialContext(ctx, rpcURL) // 实际使用时的连接

	// 假设的地址和参数
	fromAddress := common.HexToAddress("0xAADc6710b5DC4802eE6ACCa12a3b34dfc3C4FD93")
	toAddress := common.HexToAddress("0xF92F0E5AdB38f15a1e8514ea49de3f6028b8ff7d")
	value := big.NewInt(1e15) // 0.001 ETH
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(20e9) // 20 Gwei

	fmt.Println("1. ❌ 错误方式 (会导致 nonce 冲突):")
	fmt.Println("```go")
	fmt.Println("for i := 0; i < 3; i++ {")
	fmt.Println("    nonce, _ := client.PendingNonceAt(ctx, fromAddress)  // 每次获取")
	fmt.Println("    // 问题：可能都得到相同的 nonce!")
	fmt.Println("    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)")
	fmt.Println("    // 发送交易...")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("2. ✅ 正确方式1 (手动管理 nonce):")
	fmt.Println("```go")
	fmt.Printf("// 只获取一次起始 nonce\n")
	fmt.Printf("nonce, _ := client.PendingNonceAt(ctx, fromAddress)\n")
	fmt.Printf("fmt.Printf(\"起始 nonce: %%d\\n\", nonce)\n")
	fmt.Println()
	fmt.Printf("// 发送多笔交易，手动递增 nonce\n")
	fmt.Printf("for i := 0; i < 3; i++ {\n")
	fmt.Printf("    currentNonce := nonce + uint64(i)\n")
	fmt.Printf("    fmt.Printf(\"交易%%d 使用 nonce: %%d\\n\", i+1, currentNonce)\n")
	fmt.Printf("    \n")
	fmt.Printf("    tx := types.NewTransaction(currentNonce, toAddress, value, gasLimit, gasPrice, nil)\n")
	fmt.Printf("    // 签名并发送...\n")
	fmt.Printf("}\n")
	fmt.Println("```")
	fmt.Println()

	// 模拟执行
	startNonce := uint64(5) // 假设当前 nonce 是 5
	fmt.Printf("模拟执行结果:\n")
	fmt.Printf("起始 nonce: %d\n", startNonce)
	for i := 0; i < 3; i++ {
		currentNonce := startNonce + uint64(i)
		fmt.Printf("交易%d 使用 nonce: %d\n", i+1, currentNonce)
	}
	fmt.Println()

	fmt.Println("3. ✅ 正确方式2 (批量构建再发送):")
	fmt.Println("```go")
	fmt.Println("// 先构建所有交易")
	fmt.Println("nonce, _ := client.PendingNonceAt(ctx, fromAddress)")
	fmt.Println("var transactions []*types.Transaction")
	fmt.Println()
	fmt.Println("for i := 0; i < 3; i++ {")
	fmt.Println("    tx := types.NewTransaction(nonce+uint64(i), toAddress, value, gasLimit, gasPrice, nil)")
	fmt.Println("    signedTx, _ := types.SignTx(tx, signer, privateKey)")
	fmt.Println("    transactions = append(transactions, signedTx)")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("// 然后批量发送")
	fmt.Println("for i, tx := range transactions {")
	fmt.Println("    err := client.SendTransaction(ctx, tx)")
	fmt.Println("    fmt.Printf(\"交易%d 发送: %s\\n\", i+1, tx.Hash().Hex())")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()
}

// 展示高级场景：处理失败重试
func demonstrateAdvancedNonceManagement() {
	fmt.Println("=== 高级场景：处理交易失败和重试 ===\n")

	fmt.Println("4. 🔄 处理交易失败的情况:")
	fmt.Println("```go")
	fmt.Println("type NonceManager struct {")
	fmt.Println("    client   *ethclient.Client")
	fmt.Println("    address  common.Address")
	fmt.Println("    nextNonce uint64")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("func (nm *NonceManager) GetNextNonce() uint64 {")
	fmt.Println("    // 如果是第一次，从网络获取")
	fmt.Println("    if nm.nextNonce == 0 {")
	fmt.Println("        nonce, _ := nm.client.PendingNonceAt(context.Background(), nm.address)")
	fmt.Println("        nm.nextNonce = nonce")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    // 返回并递增")
	fmt.Println("    result := nm.nextNonce")
	fmt.Println("    nm.nextNonce++")
	fmt.Println("    return result")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("func (nm *NonceManager) ResetFromNetwork() {")
	fmt.Println("    // 如果有交易失败，重新从网络同步")
	fmt.Println("    nonce, _ := nm.client.PendingNonceAt(context.Background(), nm.address)")
	fmt.Println("    nm.nextNonce = nonce")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("5. ⚡ 实际使用:")
	fmt.Println("```go")
	fmt.Println("nonceManager := &NonceManager{client: client, address: fromAddress}")
	fmt.Println()
	fmt.Println("for i := 0; i < 5; i++ {")
	fmt.Println("    nonce := nonceManager.GetNextNonce()")
	fmt.Println("    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)")
	fmt.Println("    ")
	fmt.Println("    // 如果发送失败，可以重置 nonce 管理器")
	fmt.Println("    if err := client.SendTransaction(ctx, signedTx); err != nil {")
	fmt.Println("        nonceManager.ResetFromNetwork()  // 重新同步")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println("```")
}

// 警告和最佳实践
func showWarningsAndBestPractices() {
	fmt.Println("=== ⚠️  重要警告和最佳实践 ===\n")

	fmt.Println("🚨 常见陷阱:")
	fmt.Println("1. 并发发送多笔交易时每次都调用 PendingNonceAt")
	fmt.Println("2. 忽略交易失败后的 nonce 状态")
	fmt.Println("3. 不处理网络延迟导致的 nonce 不同步")
	fmt.Println()

	fmt.Println("✅ 最佳实践:")
	fmt.Println("1. 批量交易时只获取一次 nonce，然后手动递增")
	fmt.Println("2. 实现 nonce 管理器来处理复杂场景")
	fmt.Println("3. 监控交易状态，失败时重新同步 nonce")
	fmt.Println("4. 在高频交易场景下考虑使用交易池")
	fmt.Println("5. 测试网络先验证逻辑，再上主网")
	fmt.Println()

	fmt.Println("💡 简单记忆:")
	fmt.Println("- 单笔交易：每次获取 PendingNonceAt ✅")
	fmt.Println("- 多笔交易：获取一次，手动 +1 ✅")
	fmt.Println("- 高频交易：使用 nonce 管理器 ✅")
}

func main() {
	demonstrateMultipleTransactions()
	demonstrateAdvancedNonceManagement()
	showWarningsAndBestPractices()
}
