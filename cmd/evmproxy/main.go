package main

import (
	"context"
	"crypto/ecdsa"
	"evmproxy/pkg/evmproxy"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	cli "github.com/jawher/mow.cli"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	app := cli.App("evmproxy", "Deploy smart contract wallet proxy to EVM-compatible network")

	var (
		addrNetwork = app.String(cli.StringOpt{
			Name:   "addr",
			Value:  "http://localhost:8545",
			EnvVar: "ADDR",
			Desc:   "Ethereum network dial address",
		})
		privateKey = app.String(cli.StringOpt{
			Name:      "private-key",
			EnvVar:    "PRIVATE_KEY",
			HideValue: true,
			Desc:      "Private key of publisher",
		})
		wallet = app.String(cli.StringOpt{
			Name:   "wallet",
			EnvVar: "WALLET",
			Desc:   "Hex address of the proxy backend wallet",
		})
	)

	// Specify the action to execute when the app is invoked correctly
	app.Action = func() {
		client, err := ethclient.Dial(*addrNetwork)
		if err != nil {
			log.Fatal("dial: ", err)
		}
		if privateKey == nil || len(*privateKey) == 0 {
			log.Fatal("missing PRIVATE_KEY")
		}
		if wallet == nil || len(*wallet) == 0 {
			log.Fatal("missing WALLET")
		}
		privateKey, err := crypto.HexToECDSA(*privateKey)
		if err != nil {
			log.Fatal("private key: ", err)
		}
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("invalid private key")
		}

		toAddress := common.HexToAddress(*wallet)
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal("nonce error: ", err)
		}

		chainID, err := client.ChainID(context.Background())
		if err != nil {
			log.Fatal("chain error: ", err)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			log.Fatal("transaction: ", err)
		}
		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(0)      // in wei
		auth.GasLimit = uint64(3000000) // in units
		auth.GasPrice = big.NewInt(1000000)

		bytecode := evmproxy.GetBytecode(toAddress)
		version := "1.0"
		parsed, _ := abi.JSON(strings.NewReader("{}"))
		address, tx, _, err := bind.DeployContract(auth, parsed, bytecode, client, version)
		if err != nil {
			log.Fatal("deploy error: ", err)
		}

		fmt.Println("Address : ", address.Hex())
		fmt.Println("Proxy to: ", toAddress.Hex())
		fmt.Println("Tx hash : ", tx.Hash().Hex())
	}

	// Invoke the app passing in os.Args
	app.Run(os.Args)
}
