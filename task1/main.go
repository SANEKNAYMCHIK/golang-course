package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type GitHubRepo struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	CreatedAt       time.Time `json:"created_at"`
}

func getCLIValues() (string, string, error) {
	repo := flag.String("repo", "", "name of repository, that you want check")
	owner := flag.String("owner", "", "owner's name of this repository")
	flag.Parse()
	if *repo == "" || *owner == "" {
		return "", "", errors.New("repo and owner values are required")
	}
	return *repo, *owner, nil
}

func createRepoInfo(repo GitHubRepo) string {
	res := fmt.Sprintf(
		"Name: %s\nDescription: %s\nStargazersCount: %d\nForksCount: %d\nCreatedAt: %s\n",
		repo.Name, repo.Description, repo.StargazersCount, repo.ForksCount, repo.CreatedAt,
	)
	return res
}

func main() {
	repo, owner, err := getCLIValues()
	if err != nil {
		log.Fatalf("some problems with required parameters: %s", err)
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("internal problems of server: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("something went wrong: %s", resp.Status)
	}
	var repoInfo GitHubRepo
	err = json.NewDecoder(resp.Body).Decode(&repoInfo)
	if err != nil {
		log.Fatalf("incorrect decoding to Go struct: %s", err)
	}
	result := createRepoInfo(repoInfo)
	fmt.Print(result)
}
