package main

import (
	"fmt"
	"math/big"
)

// 解释复杂的 big.Float 操作
func explainBigFloatOperations() {
	fmt.Println("=== 解释 big.Float 复杂操作 ===\n")

	// 模拟一个余额 (以 Wei 为单位)
	// 1 ETH = 1,000,000,000,000,000,000 Wei (18个0)
	balanceInWei := big.NewInt(159396525300000000) // 约 0.159 ETH

	fmt.Printf("原始余额 (Wei): %s\n", balanceInWei.String())
	fmt.Println()

	// ==================== 复杂写法 (你代码中的) ====================
	fmt.Println("1. 复杂写法 (链式调用):")
	fmt.Println("   new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))")
	fmt.Println()

	// 这行代码等价于以下步骤：
	complexResult := new(big.Float).Quo(new(big.Float).SetInt(balanceInWei), big.NewFloat(1e18))
	fmt.Printf("   结果: %s ETH\n", complexResult.String())
	fmt.Println()

	// ==================== 分步解释 ====================
	fmt.Println("2. 分步解释:")
	fmt.Println()

	// 步骤1: 创建第一个 big.Float
	step1 := new(big.Float)
	fmt.Printf("   步骤1: new(big.Float) 创建空的浮点数\n")
	fmt.Printf("   结果: %s\n", step1.String())
	fmt.Println()

	// 步骤2: 将 big.Int 转为 big.Float
	step2 := new(big.Float).SetInt(balanceInWei)
	fmt.Printf("   步骤2: SetInt(%s) 设置整数值\n", balanceInWei.String())
	fmt.Printf("   结果: %s\n", step2.String())
	fmt.Println()

	// 步骤3: 创建除数
	step3 := big.NewFloat(1e18) // 1e18 = 1000000000000000000
	fmt.Printf("   步骤3: big.NewFloat(1e18) 创建除数\n")
	fmt.Printf("   1e18 = %s\n", step3.String())
	fmt.Println()

	// 步骤4: 执行除法
	step4 := new(big.Float).Quo(step2, step3)
	fmt.Printf("   步骤4: Quo(被除数, 除数) 执行除法\n")
	fmt.Printf("   %s ÷ %s = %s\n", step2.String(), step3.String(), step4.String())
	fmt.Println()

	// ==================== 简化写法 ====================
	fmt.Println("3. 更清晰的写法:")
	fmt.Println()

	// 方法1: 分步骤写
	balanceFloat := new(big.Float).SetInt(balanceInWei)
	ethDivisor := big.NewFloat(1e18)
	result1 := new(big.Float).Quo(balanceFloat, ethDivisor)
	fmt.Printf("   方法1 (分步): %s ETH\n", result1.String())

	// 方法2: 创建辅助函数
	result2 := weiToEth(balanceInWei)
	fmt.Printf("   方法2 (函数): %s ETH\n", result2.String())

	// 方法3: 使用字符串格式化
	result3 := weiToEthString(balanceInWei)
	fmt.Printf("   方法3 (字符串): %s ETH\n", result3)
	fmt.Println()

	// ==================== 为什么这么复杂？ ====================
	fmt.Println("4. 为什么需要这么复杂的操作？")
	fmt.Println()
	fmt.Println("   以太坊的数字很大，普通的 float64 不够精确:")
	fmt.Printf("   - 1 ETH = %d Wei (19位数字!)\n", int64(1e18))
	fmt.Println("   - float64 只有约16位精度")
	fmt.Println("   - big.Float 可以处理任意精度的小数")
	fmt.Println()

	// ==================== 单位转换表 ====================
	fmt.Println("5. 以太坊单位转换:")
	fmt.Println("   1 ETH = 1,000,000,000,000,000,000 Wei")
	fmt.Println("   1 ETH = 1,000,000,000 Gwei (Gas price 常用)")
	fmt.Println("   1 Gwei = 1,000,000,000 Wei")
	fmt.Println()

	// 演示不同单位
	exampleWei := big.NewInt(1500000000000000000) // 1.5 ETH
	fmt.Printf("   示例: %s Wei\n", exampleWei.String())
	fmt.Printf("   = %s ETH\n", weiToEth(exampleWei).String())
	fmt.Printf("   = %s Gwei\n", weiToGwei(exampleWei).String())
}

// 辅助函数：Wei 转 ETH
func weiToEth(wei *big.Int) *big.Float {
	eth := new(big.Float).SetInt(wei)
	return eth.Quo(eth, big.NewFloat(1e18))
}

// 辅助函数：Wei 转 Gwei
func weiToGwei(wei *big.Int) *big.Float {
	gwei := new(big.Float).SetInt(wei)
	return gwei.Quo(gwei, big.NewFloat(1e9))
}

// 辅助函数：Wei 转 ETH (字符串格式，保留6位小数)
func weiToEthString(wei *big.Int) string {
	eth := weiToEth(wei)
	return fmt.Sprintf("%.6f", eth)
}

// func main() {
// 	explainBigFloatOperations()
// }
