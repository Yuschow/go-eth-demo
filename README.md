# go-eth-demo

A simple Go application to interact with the Ethereum blockchain using the Sepolia testnet.

## Features

- Connect to Ethereum Sepolia testnet
- Query block information by block number
- Send ETH transactions

## Setup

1. Clone the repository:
```bash
git clone https://github.com/Yuschow/go-eth-demo.git
cd go-eth-demo
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure environment variables:
```bash
cp .env.example .env
```

Edit the `.env` file and set the following variables:
- `SEPOLIA_RPC`: RPC endpoint for Sepolia testnet (optional, has default value)
- `PRIVATE_KEY`: Your Ethereum private key (without 0x prefix) - **Required**
- `RECIPIENT_ADDR`: The recipient address for transactions - **Required**

**⚠️ Security Warning**: Never commit your `.env` file or expose your private key!

## Usage

现在程序会自动加载 `.env` 文件中的环境变量。只需要：

1. 确保你已经创建并配置了 `.env` 文件：
```bash
cp .env.example .env
# 编辑 .env 文件，设置你的实际值
```

2. 直接运行程序：
```bash
go run go-eth-demo/main.go
```

**替代方法**：你也可以直接设置环境变量：
```bash
export PRIVATE_KEY="your_private_key_here"
export RECIPIENT_ADDR="0xrecipient_address_here"
go run go-eth-demo/main.go
```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `SEPOLIA_RPC` | Sepolia testnet RPC endpoint | No | Alchemy default endpoint |
| `PRIVATE_KEY` | Your Ethereum private key (without 0x) | Yes | - |
| `RECIPIENT_ADDR` | Transaction recipient address | Yes | - |