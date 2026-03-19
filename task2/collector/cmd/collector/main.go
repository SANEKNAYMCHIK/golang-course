package main

import (
	"log"
	"net"

	adapter "github.com/SANEKNAYMCHIK/distrib-system/collector/internal/adapter/github"
	handler "github.com/SANEKNAYMCHIK/distrib-system/collector/internal/handler/grpc"
	"github.com/SANEKNAYMCHIK/distrib-system/collector/internal/usecase"
	pb "github.com/SANEKNAYMCHIK/distrib-system/pkg/proto"
	"google.golang.org/grpc"
)

func main() {
	githubClient := adapter.NewGitHubReposClient()
	useCase := usecase.NewRepoUseCase(githubClient)
	handler := handler.NewRepoHandler(useCase)

	server := grpc.NewServer()
	pb.RegisterRepoServiceServer(server, handler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("grpc server started")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
