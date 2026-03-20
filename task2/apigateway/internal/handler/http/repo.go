package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	data, err := rh.usecase.GetRepo(r.Context(), owner, repo)

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.NotFound:
			http.Error(w, st.Message(), http.StatusNotFound)
		case codes.Unknown:
			http.Error(w, st.Message(), http.StatusMovedPermanently)
		case codes.PermissionDenied:
			http.Error(w, st.Message(), http.StatusForbidden)
		default:
			http.Error(w, st.Message(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
