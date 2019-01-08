package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hoodcops/xcore/pkg/db"
	"github.com/hoodcops/xcore/pkg/twilio"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func startSignIn(verifier *twilio.TwilioVerifier, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			CountryCode string `json:"countryCode"`
			PhoneNumber string `json:"phoneNumber"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			respondAsBadRequest(w, NewInvalidPayloadResponse(err))
			return
		}

		errRes := NewErrorResponse("Missing values for required parameters")
		if len(payload.CountryCode) == 0 {
			errRes.AddError(NewMissingParamError("countryCode"))
		}

		if len(payload.PhoneNumber) == 0 {
			errRes.AddError(NewMissingParamError("phoneNumber"))
		}

		if errRes.HasErrors() {
			respondAsBadRequest(w, errRes)
			return
		}

		err = verifier.SendCode(payload.CountryCode, payload.PhoneNumber)
		if err != nil {
			logger.Error("failed sending verification code",
				zap.String("phone_number", payload.CountryCode+payload.PhoneNumber),
				zap.Error(err),
			)

			respondAsInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		respondWithData(w, OkResponse{Data: payload, Info: "Verification code sent successfully"})
	}
}

func verifyCode(verifier *twilio.TwilioVerifier, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			CountryCode      string `json:"countryCode"`
			PhoneNumber      string `json:"phoneNumber"`
			VerificationCode string `json:"verificationCode"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			respondAsBadRequest(w, NewInvalidPayloadResponse(err))
			return
		}

		errRes := NewErrorResponse("Missing values for required parameters")
		if len(payload.CountryCode) == 0 {
			errRes.AddError(NewMissingParamError("countryCode"))
		}

		if len(payload.PhoneNumber) == 0 {
			errRes.AddError(NewMissingParamError("phoneNumber"))
		}

		if len(payload.VerificationCode) == 0 {
			errRes.AddError(NewMissingParamError("verificationCode"))
		}

		if errRes.HasErrors() {
			respondAsBadRequest(w, errRes)
			return
		}

		err = verifier.VerifyCode(payload.CountryCode, payload.PhoneNumber, payload.VerificationCode)
		if err != nil {
			logger.Error("failed verifying phone number",
				zap.String("phone_number", payload.CountryCode+payload.PhoneNumber),
				zap.String("verification_code", payload.VerificationCode),
				zap.Error(err),
			)

			respondAsInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		respondWithData(w, OkResponse{Data: payload, Info: "Phone number verified successfully"})
	}
}

func createUser(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			PhoneNumber string `json:"phoneNumber"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			respondAsBadRequest(w, NewInvalidPayloadResponse(err))
			return
		}

	}
}

func getAllMobileUsers(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := db.NewMobileUsersRepo(dbConn)
		users, err := repo.GetAll()
		if err != nil {
			logger.Error("failed fetching all mobile users from db", zap.Error(err))
			respondAsInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		respondWithData(w, OkResponse{Data: users})
	}
}

func mobileUsersRoutes(dbConn *sqlx.DB, verifier *twilio.TwilioVerifier, secretKey string, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	// router.Post("/signin/start", ValidateJWT(startSignIn(verifier, logger), secretKey))
	router.Get("/", getAllMobileUsers(dbConn, logger))
	router.Post("/signin/start", startSignIn(verifier, logger))
	router.Post("/signin/verify", verifyCode(verifier, logger))
	// router.Post("/signin/finish", createUser(dbConn, logger))

	return router
}
