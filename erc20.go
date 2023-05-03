package main

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// YourContract is a generic contract struct
type YourContract struct {
	address common.Address
	client  *ethclient.Client
}

// DeployYourContract deploys a new ERC20 token contract
func DeployYourContract(auth *bind.TransactOpts, client *ethclient.Client, contractCode []byte, coinName string, abi abi.ABI) (common.Address, *types.Transaction, *YourContract, error) {
	// Deploy the contract
	address, tx, _, err := bind.DeployContract(auth, abi, contractCode, client)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	// Create a new instance of the YourContract struct
	contractInstance := &YourContract{address: address, client: client}

	return address, tx, contractInstance, nil
}
