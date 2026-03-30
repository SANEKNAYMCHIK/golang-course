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
	template := `
	Name: %s
	Description: %s
	StargazersCount: %d
	ForksCount: %d
	CreatedAt: %s
	`
	return fmt.Sprintf(template, repo.Name, repo.Description, repo.StargazersCount, repo.ForksCount, repo.CreatedAt)
}

func (g *GitHubReposClient) GetRepoInfo(owner, repo string) (domain.Repo, error) {
	url := fmt.Sprintf("%s/%s/%s", g.baseURL, owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("internal problems of server: %s", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNotFound:
		return domain.Repo{}, status.Error(codes.NotFound, resp.Status)
	case http.StatusMovedPermanently:
		return domain.Repo{}, status.Error(codes.Unknown, resp.Status)
	case http.StatusForbidden:
		return domain.Repo{}, status.Error(codes.PermissionDenied, resp.Status)
	case http.StatusInternalServerError:
		return domain.Repo{}, status.Error(codes.DataLoss, resp.Status)
	}
	var repoInfo gitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		log.Printf("incorrect decoding to Go struct: %s\n", err)
		return domain.Repo{}, fmt.Errorf("need to fix decoding to Go struct")
	}
	return domain.Repo{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		Stargazers:  repoInfo.StargazersCount,
		Forks:       repoInfo.ForksCount,
		CreatedAt:   repoInfo.CreatedAt,
	}, nil
}
