package usecase

import "github.com/SANEKNAYMCHIK/distrib-system/collector/internal/domain"

type RepoProvider interface {
	GetRepoInfo(owner, repo string) (domain.Repo, error)
}

type RepoUseCase struct {
	provider RepoProvider
}

func NewRepoUseCase(provider RepoProvider) *RepoUseCase {
	return &RepoUseCase{
		provider: provider,
	}
}

func (ruc *RepoUseCase) GetRepo(owner, repo string) (domain.Repo, error) {
	return ruc.provider.GetRepoInfo(owner, repo)
}
