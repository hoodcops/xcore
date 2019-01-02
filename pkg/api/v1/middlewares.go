package v1

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// ValidateJWT is a middleware that validates JWT tokens passed in request headers
// If a token is not present, an Unauthorized access response is sent. However, if
// a token is present but invalid for some reasons, an Unauthorized access response
// is sent with explanations into how the issue can be fixed.
func ValidateJWT(handler http.Handler, secret string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid JWT")
				}

				return secret, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				handler.ServeHTTP(w, r)
			}
		} else {
			fmt.Fprintf(w, "Unauthorized access")
		}
	})
}
