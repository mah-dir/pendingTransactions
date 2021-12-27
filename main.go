package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/joho/godotenv"
)

func main() {

	// Parse three arguments from command line
	// from, to, amount

	from := flag.String("from", "", "from address")
	to := flag.String("to", "", "to address")

	flag.Parse()

	if *from != "" {
		if !common.IsHexAddress(*from) {
			log.Fatal("Invalid from address")
		}
	}

	if *to != "" {
		if !common.IsHexAddress(*to) {
			log.Fatal("Invalid to address")
		}
	}

	godotenv.Load()

	wss := os.Getenv("WSS")
	https := os.Getenv("HTTPS")

	rpcClient, err := rpc.Dial(wss)

	if err != nil {
		log.Fatal("Failed to connect to rpc. ", err)
	}

	ethclient, err := ethclient.Dial(https)

	if err != nil {
		log.Fatal("Failed to connect to ethclient. ", err)
	}

	gethClient := gethclient.New(rpcClient)

	logs := make(chan common.Hash)
	_, err = gethClient.SubscribePendingTransactions(context.Background(), logs)

	if err != nil {
		log.Fatal("failed to subscribe", err)
	}

	for {
		select {
		case h := <-logs:

			tx, _, err := ethclient.TransactionByHash(context.Background(), h)

			if err != nil {
				// fmt.Println("https://etherscan.io/tx/" + h.String())
				// log.Fatal("failed to get transaction ", err)
				continue
			}

			msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), big.NewInt(0))

			if err != nil {
				log.Fatal("failed to get from address ", err)
			}

			// printTxDetails(msg.From().String(), tx.To().String(), tx.Value().String())

			if *from != "" && common.HexToAddress(*from) == msg.From() {
				printTxDetails(msg.From().String(), tx.To().String(), tx.Value().String())
				fmt.Println("https://etherscan.io/tx/" + h.String())
			} else if *to != "" && common.HexToAddress(*to) == *tx.To() {
				printTxDetails(msg.From().String(), tx.To().String(), tx.Value().String())
				fmt.Println("https://etherscan.io/tx/" + h.String())
			}

		}
	}

}

func printTxDetails(from string, to string, amount string) {
	fmt.Printf("FROM: %30s TO: %30s AMOUNT: %s \n", from, to, amount)
}
