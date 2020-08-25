package gate

import (
	"github.com/noah-blockchain/autodeleg/internal/env"
	"github.com/noah-blockchain/autodeleg/internal/helpers"
	"github.com/noah-blockchain/noah-node-go-api"
	"github.com/noah-blockchain/noah-node-go-api/errors"
	"github.com/sirupsen/logrus"
	"math/big"
)

type NoahGate struct {
	api    *noah_node_go_api.NoahNodeApi
	Logger *logrus.Logger
}

//New instance of Noah Gate
func New(logger *logrus.Logger) *NoahGate {
	return &NoahGate{
		api:    noah_node_go_api.New(env.GetEnv(env.NoahApiNodeEnv, "")),
		Logger: logger,
	}
}

//Return nonce for address
func (mg *NoahGate) GetNonce(address string) (*string, error) {
	response, err := mg.api.GetAddress(address)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"address": address,
		}).Warn(err)
		return nil, err
	}
	if response.Error != nil {
		err = errors.NewNodeError(response.Error.Message, response.Error.Code)
		mg.Logger.WithFields(logrus.Fields{
			"address": address,
		}).Warn(err)
		return nil, err
	}
	return &response.Result.TransactionCount, nil
}

//Return Balance for address
func (mg *NoahGate) GetBalance(address string) (*big.Int, error) {
	response, err := mg.api.GetAddress(address)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"address": address,
		}).Warn(err)
		return nil, err
	}
	if response.Error != nil {
		err = errors.NewNodeError(response.Error.Message, response.Error.Code)
		mg.Logger.WithFields(logrus.Fields{
			"address": address,
		}).Warn(err)
		return nil, err
	}
	qNoahBalance := helpers.StringToBigInt(response.Result.Balance["NOAH"])
	return qNoahBalance, nil
}
