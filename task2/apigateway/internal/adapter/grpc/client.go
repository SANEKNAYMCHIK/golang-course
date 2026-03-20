package adapter

import (
	"context"

	"github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/domain"
	pb "github.com/SANEKNAYMCHIK/distrib-system/pkg/proto"
)

type GRPCClient struct {
	client pb.RepoServiceClient
}

func NewGRPCClient(c pb.RepoServiceClient) *GRPCClient {
	return &GRPCClient{
		client: c,
	}
}

func (gc *GRPCClient) GetRepoInfo(ctx context.Context, owner, repo string) (domain.Repo, error) {
	resp, err := gc.client.GetRepoInfo(ctx, &pb.RepoRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return domain.Repo{}, err
	}
	return domain.Repo{
		Name:        resp.Name,
		Description: resp.Description,
		Stargazers:  resp.Stargazers,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
	}, nil
}
