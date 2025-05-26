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
)

var session = map[string]*jwt.NumericDate{}

type Claims struct {
	jwt.RegisteredClaims
	Login     string
	SessionID string
}

func BuildToken(key string, login string) (string, error) {
	sessionID := uuid.New().String()
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		Login:     login,
		SessionID: sessionID,
	})
	token, err := tokenWithClaims.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	createUserSession(sessionID)
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

func CheckUserSession(sessionID string) bool {
	_, isSessionExists := session[sessionID]
	return isSessionExists
}

func createUserSession(sessionID string) {
	session[sessionID] = jwt.NewNumericDate(time.Now().Add(TokenExp))
}
