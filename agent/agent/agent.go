// agent/worker.go
package agent

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ZolotarevAlexandr/yl_final/calculator/calculator"
)

// worker is a goroutine that continuously fetches tasks via gRPC,
// computes them, and submits the results back via gRPC.
func worker(workerID int, client *Client) {
	for {
		task, err := client.FetchTask()
		if err != nil {
			// Check if this is a "no task available" error
			if status.Code(err) == codes.NotFound {
				// No tasks available, wait and retry
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Other errors (connection issues, etc.)
			log.Printf("Worker %d: error fetching task: %v", workerID, err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Simulate long computation time
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

		// Compute the operation using calculator package
		result, err := calculator.EvaluateOperation(task.Operation, task.Arg1, task.Arg2)
		if err != nil {
			log.Printf("Worker %d: error computing task %s: %v", workerID, task.Id, err)
			continue
		}

		// Send the result back to the orchestrator via gRPC
		err = client.SendResult(task.Id, result)
		if err != nil {
			log.Printf("Worker %d: error posting result for task %s: %v", workerID, task.Id, err)
			continue
		}

		log.Printf("Worker %d: completed task %s with result %f", workerID, task.Id, result)
	}
}

// RunAgent starts a pool of worker goroutines that process tasks
func RunAgent(client *Client) {
	// Start worker pool
	for i := 0; i < ComputingPower; i++ {
		go worker(i, client)
	}

	fmt.Printf("Started agent with %d workers connected to gRPC server\n", ComputingPower)

	// Infinite wait
	select {}
}
