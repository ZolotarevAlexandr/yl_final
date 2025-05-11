package orchestrator

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	pb "github.com/ZolotarevAlexandr/yl_final/grpc"
)

func RunOrchestrator() {
	mux := http.NewServeMux()

	// Public (no JWT)
	mux.Handle("/api/v1/register",
		LoggingMiddleware(ErrorHandlingMiddleware(http.HandlerFunc(RegisterHandler))))
	mux.Handle("/api/v1/login",
		LoggingMiddleware(ErrorHandlingMiddleware(http.HandlerFunc(LoginHandler))))

	// Auth-protected chain
	authChain := func(h http.HandlerFunc) http.Handler {
		return JWTAuthMiddleware(
			LoggingMiddleware(
				ErrorHandlingMiddleware(
					http.HandlerFunc(h),
				),
			),
		)
	}

	// Healthcheck
	mux.Handle("/api/v1/ping", LoggingMiddleware(http.HandlerFunc(handlePing)))

	// Calculator endpoints
	mux.Handle("/api/v1/calculate", authChain(handleCalculate))
	mux.Handle("/api/v1/expressions", authChain(handleListExpressions))
	mux.Handle("/api/v1/expressions/", authChain(handleGetExpression))

	addr := fmt.Sprintf(":%s", Port)
	log.Printf("HTTP listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

// RunGRPCServer starts the TaskService gRPC endpoint on the given port.
// Agents will call GetTask() and SubmitTask() here.
func RunGRPCServer(port string) {
	addr := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("gRPC listen failed: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &taskServer{})
	log.Printf("gRPC listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve failed: %v", err)
	}
}
