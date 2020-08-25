package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/noah-blockchain/autodeleg/internal/env"
	"github.com/noah-blockchain/autodeleg/internal/gate"
	"github.com/noah-blockchain/autodeleg/internal/helpers"
	"github.com/noah-blockchain/noah-go-node/core/transaction"
	"github.com/sirupsen/logrus"
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
type TxDelegateResponse struct {
	Data inner `json:"data"`
}

func Delegate(c *gin.Context) {
	var err error
	var url = fmt.Sprintf("%s/api/v1/transaction/push", env.GetEnv(env.NoahGateApi, ""))
	gate, ok := c.MustGet("gate").(*gate.NoahGate)
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
				gate.Logger.WithFields(logrus.Fields{
					"transaction": tx,
				}).Errorf("Transaction decode error: ", err)
				return
			}
			decodedTx, err := transaction.TxDecoder.DecodeFromBytes(decodeString)
			if err != nil {
				gate.Logger.WithFields(logrus.Fields{
					"transaction": tx,
				}).Errorf("Transaction decode error: ", err)
				return
			}
			sender, err := decodedTx.Sender()
			if err != nil {
				gate.Logger.WithFields(logrus.Fields{
					"transaction": tx,
				}).Errorf("Transaction decode error: ", err)
				return
			}
			address := sender.String()
			nonce := decodedTx.Nonce
			qNoahBalance, err := gate.GetBalance(address)
			if err != nil {
				return
			}
			resultNonce, err := gate.GetNonce(address)
			if err != nil {
				return
			}
			addrNonce, err := strconv.ParseUint(*resultNonce, 10, 64)
			if err != nil {
				gate.Logger.WithFields(logrus.Fields{
					"address": address,
				}).Warn(err)
				return
			}
			if nonce-1 != addrNonce {
				gate.Logger.WithFields(logrus.Fields{
					"expected": nonce - 1,
					"got":      addrNonce,
				}).Info("nonce differ stop delegation")
				return
			}
			decodedTxData, ok := decodedTx.GetDecodedData().(*transaction.DelegateData)
			if !ok {
				gate.Logger.WithFields(logrus.Fields{
					"transaction": tx,
				}).Errorf("Transaction decode error: ", err)
				return
			}
			amount := decodedTxData.Value
			cmp := amount.Cmp(qNoahBalance)
			if cmp == -1 || cmp == 0 {
				var payload = make(map[string]string)
				payload["transaction"] = tx
				delResp, err := helpers.HttpPost(url, payload)
				if err != nil {
					gate.Logger.WithFields(logrus.Fields{}).Errorf("Transaction delegate error: ", err)
					return
				}
				var body map[string]interface{}
				if err = json.Unmarshal(delResp, &body); err != nil {
					gate.Logger.Error(err)
					return
				}
				if _, exists := body["data"]; exists {
					var resp TxDelegateResponse
					err = json.Unmarshal(delResp, &resp)
					if err != nil {
						gate.Logger.Error(err)
						return
					}
					gate.Logger.WithFields(logrus.Fields{
						"hash": resp.Data.Hash,
					}).Info("Tx success")
					txCount++
					// SLEEP!
					time.Sleep(time.Second * 10) // пауза 10сек, Nonce чтобы в блокчейна +1
				} else {
					gate.Logger.WithFields(logrus.Fields{
						"error": body,
					}).Warn("GATE ERROR")
				}
			}
		}
	}()
	c.JSON(http.StatusOK, gin.H{"message": "Delegation started!"})
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":    "Noah Auto-delegator API",
		"version": "0.0.1",
	})
}
