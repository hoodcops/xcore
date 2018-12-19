package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hoodcops/xcore/pkg/twilio"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func startSignIn(db *sqlx.DB, verifier *twilio.TwilioVerifier, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			CountryCode string `json:"countryCode"`
			PhoneNumber string `json:"phoneNumber"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			logger.Error("failed decoding request payload", zap.Error(err))
			http.Error(w, "failed decoding request payload", http.StatusBadRequest)
			return
		}

		err = verifier.SendCode(payload.CountryCode, payload.PhoneNumber)
		if err != nil {
			logger.Error("failed sending verification code",
				zap.String("phone_number", payload.CountryCode+payload.PhoneNumber),
				zap.Error(err),
			)

			http.Error(w, "failed sending verification code", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, Response{Data: payload, Info: "Verification code sent successfully"})
	}
}

func verifyCode(db *sqlx.DB, verifier *twilio.TwilioVerifier, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			CountryCode      string `json:"countryCode"`
			PhoneNumber      string `json:"phoneNumber"`
			VerificationCode string `json:"verificationCode"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			logger.Error("failed decoding request payload", zap.Error(err))
			http.Error(w, "failed decoding request payload", http.StatusBadRequest)
			return
		}

		err = verifier.VerifyCode(payload.CountryCode, payload.PhoneNumber, payload.VerificationCode)
		if err != nil {
			logger.Error("failed verifying phone number",
				zap.String("phone_number", payload.CountryCode+payload.PhoneNumber),
				zap.String("verification_code", payload.VerificationCode),
				zap.Error(err),
			)

			http.Error(w, "failed verifying phone number", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, Response{Data: payload, Info: "Phone number verified successfully"})
	}
}

func mobileUsersRoutes(db *sqlx.DB, verifier *twilio.TwilioVerifier, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/signin/start", startSignIn(db, verifier, logger))
	router.Post("/signin/verify", verifyCode(db, verifier, logger))

	return router
}
