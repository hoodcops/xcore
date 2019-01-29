package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hoodcops/xcore/pkg/db"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func createUserProfile(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile := new(db.UserProfile)

		err := json.NewDecoder(r.Body).Decode(profile)
		if err != nil {
			respondAsBadRequest(w, NewInvalidPayloadResponse(err))
			return
		}

		repo := db.NewUserProfilesRepo(dbConn)
		profile, err = repo.Create(profile)
		if err != nil {
			logger.Error("failed creating user profile", zap.Error(err))
			respondAsInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		respondWithData(w, OkResponse{Data: profile})
	}
}

func getAllUserProfiles(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := db.NewUserProfilesRepo(dbConn)
		users, err := repo.GetAll()
		if err != nil {
			logger.Error("failed fetching all user profiles from db", zap.Error(err))
			respondAsInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		respondWithData(w, OkResponse{Data: users})
	}
}

func userProfilesRoutes(dbConn *sqlx.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", createUserProfile(dbConn, logger))
	router.Get("/", getAllUserProfiles(dbConn, logger))

	return router
}
