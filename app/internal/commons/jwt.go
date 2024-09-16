package commons

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ConfigJWT struct {
	SecretJWT       string
	ExpiresDuration int
}

func (jwtConf *ConfigJWT) GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * time.Duration(jwtConf.ExpiresDuration)).Unix(),
	})
	return token.SignedString([]byte(jwtConf.SecretJWT))
}

func (jwtConf *ConfigJWT) ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConf.SecretJWT), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// create middleware to check token for request without framework gin
func (jwtConf *ConfigJWT) JWTMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			tokenStr := r.Header["Authorization"][0]
			tokenStr = strings.Split(tokenStr, "Bearer ")[1]
			claims, err := jwtConf.ExtractClaims(tokenStr)
			if err != nil {
				ErrorResponse(w, http.StatusUnauthorized, err)
				return
			}

			// convert claims to context
			ctx := context.WithValue(r.Context(), "user", claims["email"])
			next(w, r.WithContext(ctx))
		} else {
			ErrorResponse(w, http.StatusUnauthorized, errors.New("Token required"))
			return
		}
	})
}
