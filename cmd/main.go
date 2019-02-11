package main

import (
	"fmt"
	"github.com/abilitypoint/pkg/config"
	"github.com/abilitypoint/pkg/dbclient"
	"github.com/abilitypoint/pkg/services"
)

// PORT specifies port for server to listen on
var PORT = "4000"
var appName = "abilitypoint"
var cfg config.Config

func init() {
	cfg.LoadEnvironmentVariables()
}

func main() {
	fmt.Printf("Starting %v\n", appName)
	initializeNeo4jClient()
	service.StartWebServer(PORT)
}

func initializeNeo4jClient() {
	service.DBClient = &dbclient.Neo4jClient{}
	service.DBClient.OpenNeo4jClient(cfg.ReturnEnvironmentVariables())
}
