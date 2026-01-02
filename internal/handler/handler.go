package handler

import (
	"cmp"
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otakakot/otakakotid/pkg/api"
	"github.com/otakakot/otakakotid/pkg/schema"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	db *pgxpool.Pool
}

// New creates a new API handler instance.
func New(db *pgxpool.Pool) api.ServerInterface {
	return api.NewStrictHandler(&Handler{
		db: db,
	}, []api.StrictMiddlewareFunc{
		Middleware(),
	})
}

// Health implements api.StrictServerInterface.
func (hdl *Handler) Health(
	ctx context.Context,
	request api.HealthRequestObject,
) (api.HealthResponseObject, error) {
	if err := schema.New(hdl.db).Health(ctx); err != nil {
		return nil, err
	}

	return api.Health200Response{}, nil
}

// OpenIDConfiguration implements api.StrictServerInterface.
func (hdl *Handler) OpenIDConfiguration(
	ctx context.Context,
	request api.OpenIDConfigurationRequestObject,
) (api.OpenIDConfigurationResponseObject, error) {
	issuer := cmp.Or(os.Getenv("ISSUER"), "http://localhost:8080")

	return api.OpenIDConfiguration200JSONResponse{
		Issuer:                           issuer,
		TokenEndpoint:                    issuer + "/token",
		UserinfoEndpoint:                 issuer + "/userinfo",
		AuthorizationEndpoint:            issuer + "/authorize",
		JwksUri:                          issuer + "/certs",
		ResponseTypesSupported:           []string{"code"},
		SubjectTypesSupported:            []string{"public"},
		IdTokenSigningAlgValuesSupported: []string{"RS256"},
	}, nil
}
