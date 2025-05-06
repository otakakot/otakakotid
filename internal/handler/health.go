package handler

import (
	"net/http"

	"github.com/otakakot/otakakotid/internal/database"
)

func Health(rw http.ResponseWriter, req *http.Request) {
	if _, err := database.D1.ExecContext(req.Context(), "SELECT 1"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.Write([]byte("OK"))
}
