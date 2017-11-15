package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type conf struct {
	Endpoint       string `json:"endpoint"`
	AccountPass    string `json:"accountPass"`
	AccountAddress string `json:"accountAddress"`
}

var parsedConfig = conf{}

func main() {

	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal("unable to read config file, exiting...")
	}
	if err := json.Unmarshal(file, &parsedConfig); err != nil {
		log.Fatal("unable to marshal config file, exiting...")
	}

	md5ToSave := os.Args[1:]
	if len(md5ToSave) < 1 {
		log.Fatal("please pass in an MD5 to save")
	}

	ethClient, err := ethclient.Dial(parsedConfig.Endpoint)
	if err != nil {
		log.Println("error dialing eth endpoint: ", err)
	}

	// create a context. TODO: review
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	networkID, err := ethClient.NetworkID(ctx)
	if err != nil {
		log.Println("error getting network ID: ", err)
	}

	log.Println("network id: ", networkID)

	// keystore stuff
	myKeyStore := keystore.NewKeyStore(".", 1, 2)
	myAccounts := myKeyStore.Accounts()
	log.Println("accounts: ", myAccounts[0])

	// transaction stuff
	goodGasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("error getting suggested gas price: ", err)
	}
	log.Println("suggested gas price is: ", goodGasPrice)

	nonce, err := ethClient.NonceAt(ctx, myAccounts[0].Address, nil)
	if err != nil {
		log.Println("error getting latest nonce: ", err)
	}

	toAddress := common.HexToAddress(parsedConfig.AccountAddress)
	log.Println("sending to: ", toAddress.String())

	amount := big.NewInt(1)
	gasLimit := big.NewInt(40000)
	gasPrice := big.NewInt(0)
	gasPrice.Div(goodGasPrice, big.NewInt(2))
	dataToSend := []byte(strings.Join(md5ToSave, ""))
	dummyTrans := types.NewTransaction(nonce, toAddress, amount, gasLimit, goodGasPrice, dataToSend)

	// sign transaction
	signedTx, err := myKeyStore.SignTxWithPassphrase(myAccounts[0], parsedConfig.AccountPass, dummyTrans, networkID)
	if err != nil {
		log.Println("error signing tx: ", err)
	}
	log.Printf("%v wei + %v × %v gas", signedTx.Value(), signedTx.Gas(), signedTx.GasPrice())

	// get tx hash
	txHash := signedTx.Hash()

	// send transaction
	err = ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("error sending transaction: ", err)
		return
	}

	log.Println("tx hash: ", txHash.String())
	txPending := true
	for txPending {
		_, txPending, err = ethClient.TransactionByHash(ctx, txHash)
		if err != nil {
			log.Println("error getting tx status: ", err)
		}
		if txPending {
			log.Println("tx is pending...")
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	txRcpt, err := ethClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Println("error getting tx receipt: ", err)
	}
	log.Println("receipt: ", txRcpt.String())

}
