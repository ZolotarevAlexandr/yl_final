package main

import (
	"fmt"
	"log"

	"github.com/ZolotarevAlexandr/yl_final/agent/agent"
)

func main() {
	// build the address from env‐driven port
	addr := fmt.Sprintf("localhost:%s", agent.GRPCPort)

	// instantiate the gRPC client
	client, err := agent.NewClient(addr)
	if err != nil {
		log.Fatalf("unable to connect to orchestrator gRPC at %s: %v", addr, err)
	}
	defer client.Close()

	// hand off to RunAgent, which now uses this client
	agent.RunAgent(client)
}
