package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hoodcops/xcore/pkg/db"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func createContacts(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			UserID   int               `json:"userId"`
			Contacts []*db.UserContact `json:"contacts"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			logger.Error("failed parsing request body", zap.Error(err))
			renderBadRequest(w, NewInvalidPayloadResponse(err))
			return
		}

		for _, contact := range payload.Contacts {
			contact.UserID = payload.UserID
		}

		repo := db.NewUserContactsRepo(dbConn)
		savedContacts, err := repo.CreateContacts(payload.Contacts)
		if err != nil {
			logger.Debug("failed saving user contacts", zap.Error(err))
			renderInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		renderData(w, OkResponse{Data: savedContacts})
	}
}

func getAllContacts(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := db.NewUserContactsRepo(dbConn)
		contacts, err := repo.GetAll()
		if err != nil {
			logger.Debug("failed fetching all contacts from database", zap.Error(err))
			renderBadRequest(w, NewInternalServerErrorResponse(err))
			return
		}

		renderData(w, OkResponse{Data: contacts})
	}
}

func userContactsRoutes(dbConn *sqlx.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", createContacts(dbConn, logger))
	router.Get("/", getAllContacts(dbConn, logger))

	return router
}
