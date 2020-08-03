package main

import (
	"fmt"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "github.com/nats-io/stan.go"
	"github.com/noah-blockchain/autodeleg/internal/api"
	_ "github.com/noah-blockchain/autodeleg/internal/api"
	"github.com/noah-blockchain/autodeleg/internal/env"
	noah_node_go_api "github.com/noah-blockchain/noah-node-go-api"
)

func main() {
	nodeApi := noah_node_go_api.New(env.GetEnv(env.NoahApiNodeEnv, ""))
	fmt.Println("Starting auto-delegator service with port")
	api.Run(nodeApi)
}
