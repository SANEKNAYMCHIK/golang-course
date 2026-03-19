package handler

import (
	"context"

	"github.com/SANEKNAYMCHIK/distrib-system/collector/internal/usecase"
	pb "github.com/SANEKNAYMCHIK/distrib-system/pkg/proto"
)

type RepoHandler struct {
	pb.UnimplementedRepoServiceServer
	usecase *usecase.RepoUseCase
}

func NewRepoHandler(uc *usecase.RepoUseCase) *RepoHandler {
	return &RepoHandler{
		usecase: uc,
	}
}

func (rh *RepoHandler) GetRepoInfo(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {
	repo, err := rh.usecase.GetRepo(req.Owner, req.Repo)
	if err != nil {
		return nil, err
	}
	return &pb.RepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stargazers:  repo.Stargazers,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}, nil
}
