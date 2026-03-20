package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/usecase"
)

type RepoHandler struct {
	*http.Client
	usecase *usecase.RepoUseCase
}

func NewRepoHandler(uc *usecase.RepoUseCase) *RepoHandler {
	return &RepoHandler{
		usecase: uc,
	}
}

func (rh *RepoHandler) GetRepoInfo(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")
	log.Println(owner, repo)

	data, err := rh.usecase.GetRepo(r.Context(), owner, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	json.NewEncoder(w).Encode(data)
}
