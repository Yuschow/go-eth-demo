package main

import (
	"context"
	"fmt"
	"time"
)

// 这个文件展示了不同类型的 Context 的使用
func contextExamples() {
	fmt.Println("=== Context 示例 ===\n")

	// 1. 基本的 Background Context
	fmt.Println("1. context.Background() - 基础上下文")
	basicCtx := context.Background()
	fmt.Printf("   类型: %T\n", basicCtx)
	fmt.Println("   用途: 作为根上下文，永远不会被取消")
	fmt.Println()

	// 2. 带超时的 Context
	fmt.Println("2. context.WithTimeout() - 带超时的上下文")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 记得释放资源

	fmt.Println("   这个上下文会在5秒后自动超时")
	deadline, ok := timeoutCtx.Deadline()
	if ok {
		fmt.Printf("   截止时间: %v\n", deadline)
	}
	fmt.Println()

	// 3. 可以取消的 Context
	fmt.Println("3. context.WithCancel() - 可取消的上下文")
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	fmt.Println("   这个上下文可以手动取消")
	fmt.Println("   调用 cancelFunc() 就会取消")

	// 模拟一个可能被取消的操作
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("   >>> 2秒后自动取消上下文")
		cancelFunc() // 取消上下文
	}()

	// 等待上下文被取消
	<-cancelCtx.Done()
	fmt.Printf("   上下文被取消了，原因: %v\n", cancelCtx.Err())
	fmt.Println()

	// 4. 在实际网络操作中的使用
	fmt.Println("4. 在网络操作中使用 Context")
	fmt.Println("   在你的以太坊代码中:")
	fmt.Println("   - client.BlockByNumber(ctx, blockNumber)")
	fmt.Println("   - client.BalanceAt(ctx, address, nil)")
	fmt.Println("   - client.SendTransaction(ctx, signedTx)")
	fmt.Println()
	fmt.Println("   如果网络很慢或出现问题，Context 可以:")
	fmt.Println("   ✓ 设置超时时间避免程序卡死")
	fmt.Println("   ✓ 允许用户取消长时间运行的操作")
	fmt.Println("   ✓ 跟踪和管理所有相关的网络请求")
}

// 演示带超时的网络操作
func networkOperationWithTimeout() {
	fmt.Println("\n=== 带超时的网络操作示例 ===")

	// 创建一个3秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("开始网络操作 (最多等待3秒)...")

	// 模拟一个可能很慢的网络操作
	done := make(chan bool)
	go func() {
		// 模拟网络延迟
		time.Sleep(2 * time.Second) // 这次只需要2秒，会成功
		done <- true
	}()

	// 等待操作完成或超时
	select {
	case <-done:
		fmt.Println("✓ 网络操作成功完成!")
	case <-ctx.Done():
		fmt.Printf("✗ 操作超时了: %v\n", ctx.Err())
	}
}

// 如果你想运行这些示例，取消下面的注释
func main() {
	contextExamples()
	networkOperationWithTimeout()
}
