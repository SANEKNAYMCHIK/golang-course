package adapter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SANEKNAYMCHIK/distrib-system/collector/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GitHubReposClient struct {
	baseURL string
	client  *http.Client
}

func NewGitHubReposClient() *GitHubReposClient {
	return &GitHubReposClient{
		baseURL: "https://api.github.com/repos",
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

type gitHubRepo struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	StargazersCount int32  `json:"stargazers_count"`
	ForksCount      int32  `json:"forks_count"`
	CreatedAt       string `json:"created_at"`
}

func (repo gitHubRepo) String() string {
	res := fmt.Sprintf(
		"Name: %s\nDescription: %s\nStargazersCount: %d\nForksCount: %d\nCreatedAt: %s\n",
		repo.Name, repo.Description, repo.StargazersCount, repo.ForksCount, repo.CreatedAt,
	)
	return res
}

func (g *GitHubReposClient) GetRepoInfo(owner, repo string) (domain.Repo, error) {
	url := fmt.Sprintf("%s/%s/%s", g.baseURL, owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("internal problems of server: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return domain.Repo{}, status.Error(codes.NotFound, resp.Status)
		case http.StatusMovedPermanently:
			return domain.Repo{}, status.Error(codes.Unknown, resp.Status)
		case http.StatusForbidden:
			return domain.Repo{}, status.Error(codes.PermissionDenied, resp.Status)
		default:
			return domain.Repo{}, status.Error(codes.DataLoss, resp.Status)
		}
	}
	var repoInfo gitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		log.Fatalf("incorrect decoding to Go struct: %s", err)
	}
	return domain.Repo{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		Stargazers:  repoInfo.StargazersCount,
		Forks:       repoInfo.ForksCount,
		CreatedAt:   repoInfo.CreatedAt,
	}, nil
}
