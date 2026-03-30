package usecase

import (
	"context"
	"fmt"

	"github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/domain"
)

type RepoProvider interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (domain.Repo, error)
}

type RepoUseCase struct {
	provider RepoProvider
}

func NewRepoUseCase(provider RepoProvider) *RepoUseCase {
	return &RepoUseCase{
		provider: provider,
	}
}

func (ruc *RepoUseCase) GetRepo(ctx context.Context, owner, repo string) (domain.Repo, error) {
	if owner == "" || repo == "" {
		return domain.Repo{}, fmt.Errorf("owner and repo parameters are required")
	}
	return ruc.provider.GetRepoInfo(ctx, owner, repo)
}
