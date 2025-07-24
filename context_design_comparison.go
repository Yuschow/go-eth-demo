package main

import (
	"context"
	"fmt"
)

// 这个文件对比 Go 的显式 Context 设计与其他语言的隐式设计

// ==================== Go 的方式 (显式) ====================
func goStyleExample() {
	fmt.Println("=== Go 语言的显式 Context 方式 ===")
	fmt.Println()

	ctx := context.Background()

	// 每个函数都明确知道自己需要 context
	fmt.Println("✓ 优点:")
	fmt.Println("  - 一眼就能看出哪些函数可能有网络操作")
	fmt.Println("  - 可以清楚地控制超时和取消")
	fmt.Println("  - 函数签名就是文档")
	fmt.Println()

	// 模拟 Go 风格的 API 调用
	result1 := fetchUserData(ctx, "user123")
	result2 := updateDatabase(ctx, result1)
	sendNotification(ctx, result2)

	fmt.Println("  - 每一步都在 context 的控制之下")
	fmt.Println()
}

// Go 风格：明确传递 context
func fetchUserData(ctx context.Context, userID string) string {
	// 这里可以检查 ctx 是否被取消
	select {
	case <-ctx.Done():
		return "操作被取消"
	default:
		// 正常执行
		return fmt.Sprintf("用户数据: %s", userID)
	}
}

func updateDatabase(ctx context.Context, data string) string {
	// 每个函数都能响应取消信号
	select {
	case <-ctx.Done():
		return "数据库更新被取消"
	default:
		return "数据库已更新: " + data
	}
}

func sendNotification(ctx context.Context, message string) {
	select {
	case <-ctx.Done():
		fmt.Println("通知发送被取消")
	default:
		fmt.Println("通知已发送: " + message)
	}
}

// ==================== 其他语言的方式 (隐式) ====================
func otherLanguagesStyleExample() {
	fmt.Println("=== 其他语言的隐式方式 (模拟) ===")
	fmt.Println()

	fmt.Println("❓ 问题:")
	fmt.Println("  - 看不出哪些函数可能阻塞")
	fmt.Println("  - 超时控制通常是全局的或隐藏的")
	fmt.Println("  - 难以精确控制每个操作")
	fmt.Println()

	// 模拟其他语言的 API 调用 (没有显式 context)
	fmt.Println("// 其他语言可能这样写:")
	fmt.Println("result1 := fetchUserDataHidden('user123')")
	fmt.Println("result2 := updateDatabaseHidden(result1)")
	fmt.Println("sendNotificationHidden(result2)")
	fmt.Println()
	fmt.Println("问题: 你不知道这些函数:")
	fmt.Println("  - 需要多长时间")
	fmt.Println("  - 是否可以取消")
	fmt.Println("  - 是否有网络调用")
	fmt.Println()
}

// ==================== 实际对比 ====================
func realWorldComparison() {
	fmt.Println("=== 实际开发中的对比 ===")
	fmt.Println()

	// JavaScript/Python 风格 (隐式)
	fmt.Println("🌐 JavaScript/Python 常见写法:")
	fmt.Println(`
// JavaScript
async function transferMoney() {
    const balance = await getBalance();        // 不知道超时时间
    const result = await sendTransaction();    // 不知道能否取消
    return result;
}

// Python
def transfer_money():
    balance = get_balance()          # 可能卡死
    result = send_transaction()      # 无法控制
    return result
`)

	// Go 风格 (显式)
	fmt.Println("🔧 Go 语言写法:")
	fmt.Println(`
// Go
func TransferMoney(ctx context.Context) error {
    balance, err := getBalance(ctx)        // 明确可以超时/取消
    if err != nil {
        return err
    }
    
    result, err := sendTransaction(ctx)    // 明确在控制之下
    if err != nil {
        return err
    }
    
    return nil
}
`)
	fmt.Println()
}

// ==================== 为什么 Go 选择这种设计？ ====================
func whyGoChoseThisDesign() {
	fmt.Println("=== 为什么 Go 选择显式 Context？ ===")
	fmt.Println()

	fmt.Println("🎯 设计理念:")
	fmt.Println("1. **可预测性** - 看函数签名就知道行为")
	fmt.Println("2. **可控性** - 每一层都可以控制超时和取消")
	fmt.Println("3. **可组合性** - 容易组合不同的 context")
	fmt.Println("4. **可测试性** - 容易模拟超时和取消场景")
	fmt.Println()

	fmt.Println("🔍 具体好处:")
	fmt.Println("• 代码审查时一眼看出潜在的阻塞点")
	fmt.Println("• 新手也能快速理解哪些操作可能耗时")
	fmt.Println("• 容易添加监控和日志")
	fmt.Println("• 测试时可以轻松模拟各种异常情况")
	fmt.Println()

	fmt.Println("🤝 权衡:")
	fmt.Println("✓ 代码更冗长，但更清晰")
	fmt.Println("✓ 需要多传一个参数，但控制更精确")
	fmt.Println("✓ 学习成本稍高，但避免了隐藏的陷阱")
}

// 运行这个示例：go run context_design_comparison.go
// func main() {
// 	goStyleExample()
// 	otherLanguagesStyleExample()
// 	realWorldComparison()
// 	whyGoChoseThisDesign()
// }
