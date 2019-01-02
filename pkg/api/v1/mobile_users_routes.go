package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
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

// func createUser(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var payload = struct {
// 			PhoneNumber string `json:"phoneNumber"`
// 		}{}

// 		err := json.NewDecoder(r.Body).Decode(&payload)
// 		if err != nil {
// 			logger.Error("failed decoding request payload", zap.Error(err))
// 			http.Error(w, "failed decoding request payload", http.StatusBadRequest)
// 			return
// 		}

// 		repo := db.NewMobileUsersRepo(dbConn)
// 		user, err := repo.GetByPhoneNumber(payload.PhoneNumber)
// 		if err != nil {
// 			logger.Error("failed checking if user exists", zap.Error(err))
// 			http.Error(w, "failed checking if user exists", http.StatusInternalServerError)
// 			return
// 		}

// 		if user != nil {
// 			logger.Info("user does not already exist. creating new user")

// 			user = &db.MobileUser{
// 				Msisdn:      payload.PhoneNumber,
// 				CreatedAt:   time.Now(),
// 				LastLoginAt: time.Now(),
// 			}

// 			user, err = repo.Create(user)
// 			if err != nil {
// 				logger.Error("failed saving new mobile user", zap.Error(err))
// 				http.Error(w, "failed saving new mobile user", http.StatusBadRequest)
// 				return
// 			}

// 			render.JSON(w, r, Response{Data: user, Info: "Account created successfully"})
// 		} else {

// 			render.JSON(w, r, Response{Data: user, Info: "Account created successfully"})
// 		}

// 	}
// }

func mobileUsersRoutes(dbConn *sqlx.DB, verifier *twilio.TwilioVerifier, secretKey string, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	// router.Post("/signin/start", ValidateJWT(startSignIn(verifier, logger), secretKey))
	router.Post("/signin/start", startSignIn(verifier, logger))
	router.Post("/signin/verify", verifyCode(verifier, logger))
	// router.Post("/signin/finish", createUser(dbConn, logger))

	return router
}
