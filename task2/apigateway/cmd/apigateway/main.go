package main

import (
	"log"
	"net/http"

	adapter "github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/adapter/grpc"
	handler "github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/handler/http"
	"github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/usecase"
	pb "github.com/SANEKNAYMCHIK/distrib-system/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewRepoServiceClient(conn)

	gAdapter := adapter.NewGRPCClient(grpcClient)
	useCase := usecase.NewRepoUseCase(gAdapter)
	repoHandler := handler.NewRepoHandler(useCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/get_repo", repoHandler.GetRepoInfo)

	log.Println("server started ")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
