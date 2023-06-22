package jwt

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ole-larsen/uploader/service/settings"
)

type JwtToken struct {
	Token string `json:"token"`
}

func SignJwt(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func CreateJWTToken(email string, token string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"token": token,
	}
	return SignJwt(claims, settings.Settings.Secret)
}

func VerifyJwt(token string, secret string) (map[string]interface{}, error) {
	jwToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !jwToken.Valid {
		return nil, fmt.Errorf("invalid authorization token")
	}
	return jwToken.Claims.(jwt.MapClaims), nil
}

func GetBearerToken(header string) (string, error) {
	if header == "" {
		return "", fmt.Errorf("an authorization header is required")
	}
	token := strings.Split(header, " ")
	if len(token) != 2 {
		return "", fmt.Errorf("malformed bearer token")
	}
	return token[1], nil
}
