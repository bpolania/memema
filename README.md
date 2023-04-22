# Meme Machine 

This project is a Golang server that interacts with the OpenAI API and the Arbitrum network to generate and deploy meme coins. It retrieves the latest news, generates a meme coin name, creates a white paper, and deploys an ERC20 token on the Arbitrum network.

## Prerequisites

- Go 1.16 or later
- An OpenAI API key
- An Ethereum private key with some Ether for gas fees
- Infura or another Ethereum provider with Arbitrum support

## Installation

1. Clone this repository:
```
git clone https://github.com/yourusername/meme-coin-generator.git
```
2. Change into the project directory:
```
cd meme-coin-generator
```
3. Install the required Go packages:
```
go get -u github.com/ethereum/go-ethereum
go get -u github.com/joho/godotenv
```

## Configuration

1. Create a `.env` file in the project's root directory with the following contents:
```
OPENAI_API_KEY=<your_openai_api_key>
ETH_PRIVATE_KEY=<your_ethereum_private_key>
INFURA_ARBITRUM_RPC_URL=<your_infura_arbitrum_rpc_url>
```

Replace `<your_openai_api_key>`, `<your_ethereum_private_key>`, and `<your_infura_arbitrum_rpc_url>` with your actual API keys and URLs.

## Usage

1. Start the server:
```
go run main.go
```

2. Send an HTTP request to the `/deployMemeCoin` endpoint using a tool like `curl`:
```
curl http://localhost:8080/deployMemeCoin
```

This request will trigger the meme coin creation and deployment process. The server will return a JSON response containing the meme coin name and token contract address.

## abi

Create another project to compile to contract and get the abi
```
solc --abi --bin --overwrite -o build path/to/your/erc20/contract/template.sol
```
then 
```
abigen --abi=build/template.abi --bin=build/template.bin --pkg=erc20 --out=erc20/erc20.go
```


