package handler

import (
	"context"
	"net/http"

	"github.com/otakakot/otakakotid/pkg/api"
)

func Middleware() api.StrictMiddlewareFunc {
	return func(next api.StrictHandlerFunc, operationID string) api.StrictHandlerFunc {
		return func(
			ctx context.Context,
			w http.ResponseWriter,
			r *http.Request,
			request any,
		) (any, error) {
			return next(ctx, w, r, request)
		}
	}
}
