package usecase

import (
	"context"

	"github.com/SANEKNAYMCHIK/distrib-system/pkg/domain"
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
	return ruc.provider.GetRepoInfo(ctx, owner, repo)
}
