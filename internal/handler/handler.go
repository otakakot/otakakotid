package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otakakot/otakakotid/pkg/api"
	"github.com/otakakot/otakakotid/pkg/schema"
)

var _ api.ServerInterface = (*Handler)(nil)

type Handler struct {
	db *pgxpool.Pool
}

// New creates a new Handler.
func New(db *pgxpool.Pool) *Handler {
	return &Handler{
		db: db,
	}
}

// Health implements api.ServerInterface.
func (hdl *Handler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := schema.New(hdl.db).Health(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
