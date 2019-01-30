package v1

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
			logger.Error("failed parsing request body", zap.Error(err))
			renderBadRequest(w, NewInvalidPayloadResponse(err))
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
			renderBadRequest(w, errRes)
			return
		}

		err = verifier.SendCode(payload.CountryCode, payload.PhoneNumber)
		if err != nil {
			logger.Error("failed sending verification code",
				zap.String("phone_number", payload.CountryCode+payload.PhoneNumber),
				zap.Error(err),
			)

			renderInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		renderData(w, OkResponse{Data: payload, Info: "Verification code sent successfully"})
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
			renderBadRequest(w, NewInvalidPayloadResponse(err))
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
			renderBadRequest(w, errRes)
			return
		}

		err = verifier.VerifyCode(payload.CountryCode, payload.PhoneNumber, payload.VerificationCode)
		if err != nil {
			logger.Error("failed verifying phone number",
				zap.String("phone_number", payload.CountryCode+payload.PhoneNumber),
				zap.String("verification_code", payload.VerificationCode),
				zap.Error(err),
			)

			renderInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		renderData(w, OkResponse{Data: payload, Info: "Phone number verified successfully"})
	}
}

func createUser(dbConn *sqlx.DB, secretKey string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			PhoneNumber string `json:"phoneNumber"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			renderBadRequest(w, NewInvalidPayloadResponse(err))
			return
		}

		errRes := NewErrorResponse("Missing values for required parameters")
		if len(payload.PhoneNumber) == 0 {
			errRes.AddError(NewMissingParamError("phoneNumber"))
		}

		if errRes.HasErrors() {
			renderBadRequest(w, errRes)
			return
		}

		repo := db.NewMobileUsersRepo(dbConn)
		user, err := repo.GetByPhoneNumber(payload.PhoneNumber)
		if err != nil {
			logger.Error("failed checking if user already exists in db", zap.String("phoneNumber", payload.PhoneNumber), zap.Error(err))
			renderInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		// user already exists in database
		if user != nil {
			token, err := generateToken(user.Msisdn)
			if err != nil {
				logger.Error("failed generating JWT for user", zap.String("phoneNumber", user.Msisdn), zap.Error(err))
				renderInternalServerError(w, NewInternalServerErrorResponse(err))
				return
			}

			renderData(w, struct {
				Data      interface{} `json:"data"`
				AuthToken string      `json:"authToken"`
				Info      string      `json:"info"`
			}{
				Data:      user,
				AuthToken: token,
				Info:      "Welcome back!",
			})
			return
		}

		user = &db.MobileUser{
			Msisdn: payload.PhoneNumber,
		}

		user, err = repo.Create(user)
		if err != nil {
			logger.Error("failed saving user into db", zap.String("phoneNumber", payload.PhoneNumber), zap.Error(err))
			renderInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		token, err := generateToken(user.Msisdn)
		if err != nil {
			logger.Error("failed generating JWT for user", zap.String("phoneNumber", user.Msisdn), zap.Error(err))
			renderInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		renderData(w, struct {
			Data      interface{} `json:"data"`
			AuthToken string      `json:"authToken"`
			Info      string      `json:"info"`
		}{
			Data:      user,
			AuthToken: token,
			Info:      "Welcome to Hoodcops!",
		})
	}
}

func getAllMobileUsers(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := db.NewMobileUsersRepo(dbConn)
		users, err := repo.GetAll()
		if err != nil {
			logger.Error("failed fetching all mobile users from db", zap.Error(err))
			renderInternalServerError(w, NewInternalServerErrorResponse(err))
			return
		}

		renderData(w, OkResponse{Data: users})
	}
}

func generateToken(msisdn string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "edward pie"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(`sdfsd`))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func mobileUsersRoutes(dbConn *sqlx.DB, verifier *twilio.TwilioVerifier, secretKey string, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	// router.Post("/signin/start", ValidateJWT(startSignIn(verifier, logger), secretKey))
	router.Get("/", getAllMobileUsers(dbConn, logger))
	router.Post("/", createUser(dbConn, secretKey, logger))
	router.Post("/{userId}/profile", createUserProfile(dbConn, logger))
	router.Post("/signin/start", startSignIn(verifier, logger))
	router.Post("/signin/verify", verifyCode(verifier, logger))

	return router
}
