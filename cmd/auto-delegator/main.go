package main

import (
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "github.com/nats-io/stan.go"
	"github.com/noah-blockchain/autodeleg/internal/api"
	_ "github.com/noah-blockchain/autodeleg/internal/api"
	"github.com/noah-blockchain/autodeleg/internal/env"
	"github.com/noah-blockchain/autodeleg/internal/gate"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	//Init Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	if env.GetEnvAsBool(env.DebugModeEnv, true) {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.WarnLevel)
	}

	gateService := gate.New(logger)
	api.Run(gateService)
}
