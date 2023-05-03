package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const openaiAPI = "https://api.openai.com/v1/"
const chatGPTAPI = "https://api.openai.com/v1/engines/davinci-codex/completions"

const erc20ContractTemplate = "./contracts/template.sol"
const erc20ContractAbi = "./contracts/template.json"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/deployMemeCoin", deployMemeCoinHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func randomThreeLetters(input string) string {
	rand.Seed(time.Now().UnixNano())

	// Remove digits and spaces from the input string
	filtered := strings.Map(func(r rune) rune {
		if '0' <= r && r <= '9' || r == ' ' {
			return -1
		}
		return r
	}, input)

	// If the filtered string has fewer than 3 characters, add random letters
	for len(filtered) < 3 {
		filtered += string(rune(rand.Intn('z'-'a'+1) + 'a'))
	}

	// Pick 3 random characters from the filtered string
	var result string
	for i := 0; i < 3; i++ {
		index := rand.Intn(len(filtered))
		result += string(filtered[index])
	}

	// Convert the result string to uppercase
	return strings.ToUpper(result)
}

func deployMemeCoinHandler(w http.ResponseWriter, r *http.Request) {
	news := getNews()
	coinName := generateMemeCoinName(news)
	coinSymbol := randomThreeLetters(coinName)
	whitePaper := generateWhitePaper(coinName)

	tokenAddress, err := deployERC20Token(coinName, coinSymbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("./whitepapers/%s_whitepaper.txt", coinName), []byte(whitePaper), 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"coinName":     coinName,
		"tokenAddress": tokenAddress.Hex(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getNews() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", openaiAPI+"news", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var newsData map[string]interface{}
	json.Unmarshal(body, &newsData)

	news := newsData["data"].(string)
	return news
}

func generateMemeCoinName(news string) string {
	client := &http.Client{}
	data := map[string]interface{}{
		"prompt":     fmt.Sprintf("Generate a meme coin name based on this news: %s", news),
		"max_tokens": 5,
	}
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", chatGPTAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response map[string]interface{}
	json.Unmarshal(body, &response)
	choices := response["choices"].([]interface{})
	coinName := strings.TrimSpace(choices[0].(map[string]interface{})["text"].(string))

	return coinName
}

func generateWhitePaper(coinName string) string {
	client := &http.Client{}
	data := map[string]interface{}{
		"prompt":     fmt.Sprintf("Generate a white paper for a meme coin called %s.", coinName),
		"max_tokens": 300,
	}
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", chatGPTAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response map[string]interface{}
	json.Unmarshal(body, &response)
	choices := response["choices"].([]interface{})
	whitePaper := strings.TrimSpace(choices[0].(map[string]interface{})["text"].(string))

	return whitePaper
}

func deployERC20Token(coinName string, coinSymbol string) (common.Address, error) {
	// Load contract code and replace placeholder with coin name
	contractCode, err := ioutil.ReadFile(erc20ContractTemplate)
	contractCodeStr := strings.Replace(string(contractCode), "<TOKEN_NAME>", coinName, -1)
	contractCodeStr = strings.Replace(string(contractCode), "<TOKEN_SYMBOL>", coinSymbol, -1)
	abiFile, err := ioutil.ReadFile(erc20ContractAbi)
	var contractABI abi.ABI
	err = json.Unmarshal(abiFile, &contractABI)

	// Convert the contract code string to []byte
	contractCodeBytes := []byte(contractCodeStr)

	// Connect to Ethereum client
	client, err := ethclient.Dial(os.Getenv("INFURA_ARBITRUM_RPC_URL"))
	if err != nil {
		return common.Address{}, err
	}

	// other code...
	privateKeyHex := os.Getenv("ETH_PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	address, _, _, err := DeployYourContract(auth, client, contractCodeBytes, coinName, contractABI)
	if err != nil {
		return common.Address{}, err
	}
	if !ok {
		return common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}

	return address, nil
}
