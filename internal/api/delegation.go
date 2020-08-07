package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/noah-blockchain/autodeleg/internal/env"
	"github.com/noah-blockchain/autodeleg/internal/helpers"
	"github.com/noah-blockchain/noah-go-node/core/transaction"
	noah_node_go_api "github.com/noah-blockchain/noah-node-go-api"
	"net/http"
	"strconv"
	"time"
)

type Delegations struct {
	Txs []string `json:"transactions"`
}

type inner struct {
	Hash string `json:"hash"`
}
type transactionResponse struct {
	Data inner `json:"data"`
}

func Delegate(c *gin.Context) {
	var err error
	var url = fmt.Sprintf("%s/api/v1/transaction/push", env.GetEnv(env.NoahGateApi, ""))
	gate, ok := c.MustGet("gate").(*noah_node_go_api.NoahNodeApi)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Type cast error",
			},
		})
		return
	}
	var dlg Delegations
	if err = c.ShouldBindJSON(&dlg); err != nil {
		//	gate.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//start delegation
	go func() {
		txCount := 0
		txLen := len(dlg.Txs)
		for txCount < txLen {
			tx := dlg.Txs[txCount]
			decodeString, err := hex.DecodeString(tx)
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return
			}
			decodedTx, err := transaction.TxDecoder.DecodeFromBytes(decodeString)
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return
			}
			sender, err := decodedTx.Sender()
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return
			}
			response, err := gate.GetAddress(sender.String())
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return
			}
			if response.Error != nil {
				fmt.Println("ERROR: ", response.Error)
				return
			}

			nonce := decodedTx.Nonce
			qNoahBalance := helpers.StringToBigInt(response.Result.Balance["NOAH"])
			accNonce, err := strconv.ParseUint(response.Result.TransactionCount, 10, 64)
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return
			}
			if nonce-1 != accNonce {
				fmt.Println("ERROR: ", err.Error())
				return
			}
			txData, ok := decodedTx.GetDecodedData().(*transaction.DelegateData)
			if !ok {
				fmt.Println("ERROR type casting ")
				return
			}
			amount := txData.Value
			cmp := amount.Cmp(qNoahBalance)
			if cmp == -1 || cmp == 0 {
				var payload = make(map[string]string)
				payload["transaction"] = tx
				response, err := helpers.HttpPost(url, payload)
				if err != nil {
					fmt.Println("ERROR: ", err.Error())
					return
				}
				var dat map[string]interface{}
				if err = json.Unmarshal(response, &dat); err != nil {
					fmt.Println("ERROR: ", err)
				}
				if _, exists := dat["data"]; exists {
					var txresponse transactionResponse
					err = json.Unmarshal(response, &txresponse)
					if err != nil {
						fmt.Println("ERROR: ", err)
					}
					fmt.Println("TX HASH: ", txresponse.Data.Hash)
					txCount++
					// SLEEP!
					time.Sleep(time.Second * 10) // пауза 10сек, Nonce чтобы в блокчейна +1
				} else {
					fmt.Println("ERROR: ", dat)
				}
			}
		}
		fmt.Println("Delegation success")
	}()
	c.JSON(http.StatusOK, "Delegation is started")
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":    "Noah Auto-delegator API",
		"version": "0.0.1",
	})
}
