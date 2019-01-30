package v1

import (
	"github.com/go-chi/chi"
	"github.com/hoodcops/xcore/pkg/twilio"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// InitRoutes sets up all the endpoints exposed under this
// version of the API
func InitRoutes(
	db *sqlx.DB,
	verifier *twilio.TwilioVerifier,
	secret string,
	logger *zap.Logger,
) *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/v1/users", mobileUsersRoutes(db, verifier, secret, logger))
	router.Mount("/v1/profiles", userProfilesRoutes(db, logger))
	router.Mount("/v1/contacts", userContactsRoutes(db, logger))

	return router
}
