package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mpgallage/xmcrud/handlers"
	"github.com/mpgallage/xmcrud/models"
	"net/http"
	"strings"
)

func JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("there was an err")
					}
					return handlers.JwtKey, nil
				})
				if err != nil {
					_ = json.NewEncoder(w).Encode(models.Exception{Message: err.Error()})
					return
				}
				if token.Valid {
					next.ServeHTTP(w, r)
				} else {
					_ = json.NewEncoder(w).Encode(models.Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			_ = json.NewEncoder(w).Encode(models.Exception{Message: "An authorization header is required"})
		}
	})
}
