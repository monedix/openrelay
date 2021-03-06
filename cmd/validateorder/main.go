package main

import (
	"os"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/config"
	"log"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func hexToBytes(address string) [20]byte {
	slice, err := hex.DecodeString(address)
	if err != nil {
		return [20]byte{}
	}
	output := [20]byte{}
	copy(output[:], slice[:])
	return output
}

func main() {
	rpcURL := os.Args[1]
	orderFile := os.Args[2]
	feeToken, err := config.NewRpcFeeToken(rpcURL)
	if err != nil {
		log.Fatalf("Error creating RpcOrderValidator: '%v'", err.Error())
	}
	tokenProxy, err := config.NewRpcTokenProxy(rpcURL)
	if err != nil {
		log.Fatalf("Error creating RpcOrderValidator: '%v'", err.Error())
	}
	fundChecker, err := funds.NewRpcOrderValidator(rpcURL, feeToken, tokenProxy, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
	filledLookup, err := funds.NewRPCFilledLookup(rpcURL, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
	newOrder := types.Order{}
	if orderData, err := ioutil.ReadFile(orderFile); err == nil {
		if err := json.Unmarshal(orderData, &newOrder); err != nil {
			log.Fatalf(err.Error())
		}
	}
	newOrder.TakerAssetAmountFilled, err = filledLookup.GetAmountFilled(&newOrder)
	if err != nil {
		log.Fatalf(err.Error())
	}
	newOrder.Cancelled, err = filledLookup.GetCancelled(&newOrder)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(fundChecker.ValidateOrder(&newOrder))
}
