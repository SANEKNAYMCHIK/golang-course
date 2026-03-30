package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/SANEKNAYMCHIK/distrib-system/apigateway/docs"
	adapter "github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/adapter/grpc"
	handler "github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/handler/http"
	"github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/usecase"
	pb "github.com/SANEKNAYMCHIK/distrib-system/pkg/proto"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const gRPC_env = "GRPC_CONN"

// @title           GitHub Repo API
// @version         1.0
// @description     API gateway for GitHub repositories.
// @host      localhost:8080
// @BasePath  /
// @schemes   http
func main() {
	grpcConn := os.Getenv(gRPC_env)
	if grpcConn == "" {
		grpcConn = "localhost:50051"
	}
	conn, err := grpc.NewClient(grpcConn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewRepoServiceClient(conn)

	gAdapter := adapter.NewGRPCClient(grpcClient)
	useCase := usecase.NewRepoUseCase(gAdapter)
	repoHandler := handler.NewRepoHandler(useCase)

	r := chi.NewRouter()

	r.Get("/repos/{owner}/{repo}", repoHandler.GetRepoInfo)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/repos", repoHandler.GetRepoInfo)
	// mux.HandleFunc("/swagger", httpSwagger.WrapHandler)

	log.Println("server started ")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
