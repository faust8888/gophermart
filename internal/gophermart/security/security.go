package security

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	TokenExp               = time.Hour * 3
	AuthorizationTokenName = "Authorization"
	AuthSecretKey          = "secret"
)

var session = map[string]*jwt.NumericDate{}

type Claims struct {
	jwt.RegisteredClaims
	Login     string
	SessionId string
}

func BuildToken(key string, login string) (string, error) {
	sessionId := uuid.New().String()
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		Login:     login,
		SessionId: sessionId,
	})
	token, err := tokenWithClaims.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	createUserSession(sessionId)
	return token, nil
}

func GetToken(req *http.Request) string {
	tokenCookie, err := req.Cookie(AuthorizationTokenName)
	if tokenCookie != nil && err == nil {
		return tokenCookie.Value
	}
	return ""
}

func GetClaims(token string, encodedKey string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(encodedKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("get user id: claims parsing: %w", err)
	}
	return claims, nil
}

func CheckUserSession(sessionId string) bool {
	_, isSessionExists := session[sessionId]
	return isSessionExists
}

func createUserSession(sessionId string) {
	session[sessionId] = jwt.NewNumericDate(time.Now().Add(TokenExp))
}
