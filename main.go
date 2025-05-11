package main

import (
	"fmt"
	"log"

	"github.com/ZolotarevAlexandr/yl_final/agent/agent"
	"github.com/ZolotarevAlexandr/yl_final/db"
	"github.com/ZolotarevAlexandr/yl_final/orchestrator/orchestrator"
)

func main() {
	db.Init("calc.db")

	go orchestrator.RunOrchestrator()
	go orchestrator.RunGRPCServer("50051")

	addr := fmt.Sprintf("localhost:%s", agent.GRPCPort)
	client, err := agent.NewClient(addr)
	if err != nil {
		log.Fatalf("unable to connect to orchestrator gRPC at %s: %v", addr, err)
	}
	defer client.Close()

	go agent.RunAgent(client)
	select {}
}
