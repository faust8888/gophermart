package handler

import (
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	srv := startTestServer(t)
	defer srv.server.Close()

	tests := []struct {
		name            string
		login           string
		password        string
		mockErrorReturn interface{}
		wantCode        int
	}{
		{
			name:            "Successful login",
			login:           "faust8888",
			password:        "qwerty123",
			mockErrorReturn: nil,
			wantCode:        http.StatusOK,
		},
		{
			name:            "Failed login (user not found)",
			login:           "faust1111",
			password:        "qwerty123",
			mockErrorReturn: postgres.ErrUserNotExist,
			wantCode:        http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			{
				body := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, test.login, test.password)
				srv.mockUserRepository.EXPECT().CheckUser(test.login, test.password).Return(test.mockErrorReturn)

				resp, _ := createPostRequest(srv.GetFullPath(UserLoginHandlerPath), body).Send()

				assert.Equal(t, test.wantCode, resp.StatusCode())
				token := getTokenFromResponse(resp)

				if test.wantCode == http.StatusOK {
					assert.NotEmpty(t, token)

					claims, err := security.GetClaims(token, security.AuthSecretKey)
					assert.NoError(t, err)

					assert.True(t, security.CheckUserSession(claims.SessionID))
				} else {
					assert.Empty(t, token)
				}
			}
		})
	}
}

func TestLogin_WithEmptyLogin(t *testing.T) {
	srv := startTestServer(t)
	defer srv.server.Close()

	emptyLogin, password := "", "qwerty123"
	body := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, emptyLogin, password)
	srv.mockUserRepository.EXPECT().CreateUser(emptyLogin, password).AnyTimes().Return(nil)

	resp, _ := createPostRequest(srv.GetFullPath(UserLoginHandlerPath), body).Send()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
}
