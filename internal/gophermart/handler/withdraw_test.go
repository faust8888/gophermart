package handler

import (
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWithdrawHandler(t *testing.T) {
	srv := startTestServer(t)
	defer srv.server.Close()

	ctx := gomock.Any()

	tests := []struct {
		name            string
		orderNumber     int64
		withdrawSum     float32
		userLogin       string
		mockErrorReturn interface{}
		wantCode        int
	}{
		{
			name:            "Successful withdraw",
			userLogin:       "faust88888",
			orderNumber:     4539319503436467,
			withdrawSum:     320,
			mockErrorReturn: nil,
			wantCode:        http.StatusOK,
		},
		{
			name:            "Unsuccessful withdraw (not enough balance)",
			userLogin:       "faust88888",
			orderNumber:     4539319503436467,
			withdrawSum:     320,
			mockErrorReturn: postgres.ErrNotEnoughBalance,
			wantCode:        http.StatusPaymentRequired,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			{
				srv.mockWithdrawRepository.EXPECT().Withdraw(ctx, test.userLogin, test.orderNumber, test.withdrawSum).Return(test.mockErrorReturn)

				body := fmt.Sprintf(`{"order": "%d", "sum": %f}`, test.orderNumber, test.withdrawSum)
				req := createPostRequest(
					srv.GetFullPath(WithdrawHandlerPath),
					body,
				)
				req.SetCookie(
					&http.Cookie{
						Name:  security.AuthorizationTokenName,
						Value: registrationAndAuthentication(t, srv, test.userLogin),
					},
				)

				resp, _ := req.Send()

				assert.Equal(t, test.wantCode, resp.StatusCode())
			}
		})
	}
}
