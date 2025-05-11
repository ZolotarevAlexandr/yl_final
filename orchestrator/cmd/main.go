package main

import (
	"sync"

	"github.com/ZolotarevAlexandr/yl_final/db"
	"github.com/ZolotarevAlexandr/yl_final/orchestrator/orchestrator"
)

func main() {
	db.Init("calc.db")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); orchestrator.RunOrchestrator() }()
	go func() { defer wg.Done(); orchestrator.RunGRPCServer("50051") }()
	wg.Wait()
}
