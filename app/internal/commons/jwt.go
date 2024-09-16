package commons

import (
	"context"
	"net/http"
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
func (jwtConf *ConfigJWT) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			tokenStr := r.Header["Authorization"][0]
			claims, err := jwtConf.ExtractClaims(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Token required", http.StatusUnauthorized)
			return
		}
	})
}
