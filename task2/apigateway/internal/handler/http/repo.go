package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SANEKNAYMCHIK/distrib-system/apigateway/internal/usecase"
	"github.com/go-chi/chi"
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

// GetRepoInfo godoc
// @Summary      Get repository information
// @Description  Returns information about a GitHub repository by owner and repo name.
// @Tags         repos
// @Accept       json
// @Produce      json
// @Param        owner   path      string  true  "Repository owner (user or organization)"
// @Param        repo    path      string  true  "Repository name"
// @Success      200  {object}  domain.Repo
// @Failure      301  {string}  string  "Moved permanently"
// @Failure      403  {string}  string  "Forbidden"
// @Failure      404  {string}  string  "Not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /repos/{owner}/{repo} [get]
func (rh *RepoHandler) GetRepoInfo(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	data, err := rh.usecase.GetRepo(r.Context(), owner, repo)

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendGRPCError(w, st)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendGRPCError(w http.ResponseWriter, st *status.Status) {
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
}
