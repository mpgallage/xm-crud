package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mpgallage/xmcrud/models"
	"net/http"
	"os"
	"time"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func CreateToken(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
		"exp":      time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	})
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		fmt.Println(err)
	}
	_ = json.NewEncoder(w).Encode(models.JwtToken{Token: tokenString})
}
