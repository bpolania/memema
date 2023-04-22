package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const openaiAPI = "https://api.openai.com/v1/"
const chatGPTAPI = "https://api.openai.com/v1/engines/davinci-codex/completions"

const erc20ContractTemplate = "path/to/your/erc20/contract/template.sol"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/deployMemeCoin", deployMemeCoinHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func deployMemeCoinHandler(w http.ResponseWriter, r *http.Request) {
	news := getNews()
	coinName := generateMemeCoinName(news)
	whitePaper := generateWhitePaper(coinName)

	tokenAddress, err := deployERC20Token(coinName)
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

func deployERC20Token(coinName string) (common.Address, error) {
	// Load contract code and replace placeholder with coin name
	contractCode, err := ioutil.ReadFile(erc20ContractTemplate)
	if err != nil {
		return common.Address{}, err
	}
	contractCodeStr := strings.Replace(string(contractCode), "<TOKEN_NAME>", coinName, -1)

	// Convert the contract code string to []byte
	contractCodeBytes := []byte(contractCodeStr)

	// Connect to Ethereum client
	client, err := ethclient.Dial(os.Getenv("INFURA_ARBITRUM_RPC_URL"))
	if err != nil {
		return common.Address{}, err
	}

	// other code...

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	address, _, _, err := erc20.DeployYourContract(auth, client, contractCodeBytes, coinName)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}
