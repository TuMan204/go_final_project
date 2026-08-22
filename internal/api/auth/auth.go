package auth

import (
	"crypto/sha256"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(pass string) (string, error) {
	claims := jwt.MapClaims{
		"payload": sha256.Sum256([]byte(pass)),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString([]byte(pass))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func Auth(nextHandler http.HandlerFunc, pass string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(pass) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			var isValid bool
			passJWT, err := GenerateJWT(pass)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if passJWT == jwt {
				isValid = true
			}
			if !isValid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		nextHandler(w, r)
	})
}
