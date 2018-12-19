package twilio

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func NewMockTwilioVerificationHander(shouldFail bool) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/protected/json/phones/verification/start", func(w http.ResponseWriter, r *http.Request) {
		if shouldFail {
			data := `{
				"message": "No pending verifications for +49 179-449-1095 found.",
				"success": false,
				"errors": {
					"message": "No pending verifications for +49 179-449-1095 found."
				},
				"error_code": "60023"
			}`

			http.Error(w, data, http.StatusInternalServerError)
			return
		}
		data := `{
				"carrier": "Telefonica (O2 Germany GmbH & Co. OHG)",
				"is_cellphone": true,
				"message": "Text message sent to +49 179-449-1095.",
				"seconds_to_expire": 599,
				"uuid": "ac5ed8e0-e5af-0136-a779-0a5b7c2a32fe",
				"success": true
			}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))

	})

	mux.HandleFunc("/protected/json/phones/verification/check", func(w http.ResponseWriter, r *http.Request) {
		if shouldFail {
			data := `{
				"message": "No pending verifications for +49 179-449-1095 found.",
				"success": false,
				"errors": {
					"message": "No pending verifications for +49 179-449-1095 found."
				},
				"error_code": "60023"
			}`

			http.Error(w, data, http.StatusInternalServerError)
			return
		}

		data := `{
			"message": "Verification code is correct.",
			"success": true
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

	return mux
}

func TestTwilioVerifierSendCode_ShouldPass(t *testing.T) {
	srv := httptest.NewServer(NewMockTwilioVerificationHander(false))
	defer srv.Close()

	cl := &http.Client{Timeout: 10 * time.Second}
	locale := "en"
	apiKey := "50m3h@rd2gu355t3xt0rh@5h"
	host := srv.URL

	verifier := NewTwilioVerifier(cl, host, locale, apiKey)
	err := verifier.SendCode("49", "179-449-1095")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTwilioVerifierSendCode_ShouldFail(t *testing.T) {
	srv := httptest.NewServer(NewMockTwilioVerificationHander(true))
	defer srv.Close()

	cl := &http.Client{Timeout: 10 * time.Second}
	locale := "en"
	apiKey := "50m3h@rd2gu355t3xt0rh@5h"
	host := srv.URL

	verifier := NewTwilioVerifier(cl, host, locale, apiKey)
	err := verifier.SendCode("49", "179-449-1095")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTwilioVerifierVerifyCode_ShouldPass(t *testing.T) {
	srv := httptest.NewServer(NewMockTwilioVerificationHander(false))
	defer srv.Close()

	cl := &http.Client{Timeout: 10 * time.Second}
	locale := "en"
	apiKey := "50m3h@rd2gu355t3xt0rh@5h"
	host := srv.URL

	verifier := NewTwilioVerifier(cl, host, locale, apiKey)
	err := verifier.VerifyCode("49", "179-449-1095", "4591")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTwilioVerifierVerifyCode_ShouldFail(t *testing.T) {
	srv := httptest.NewServer(NewMockTwilioVerificationHander(true))
	defer srv.Close()

	cl := &http.Client{Timeout: 10 * time.Second}
	locale := "en"
	apiKey := "50m3h@rd2gu355t3xt0rh@5h"
	host := srv.URL

	verifier := NewTwilioVerifier(cl, host, locale, apiKey)
	err := verifier.VerifyCode("49", "179-449-1095", "4591")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
