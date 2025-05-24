package handler

import (
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestNewCreateOrderHandler(t *testing.T) {
	srv := startTestServer(t)
	defer srv.server.Close()

	tests := []struct {
		name            string
		orderNumber     int64
		userLogin       string
		mockErrorReturn interface{}
		wantCode        int
	}{
		{
			name:            "Successful creation of order",
			userLogin:       "faust88888",
			orderNumber:     4539319503436467,
			mockErrorReturn: nil,
			wantCode:        http.StatusAccepted,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			{
				srv.mockOrderRepository.EXPECT().BeginTransaction()
				srv.mockOrderRepository.EXPECT().CreateOrder(test.userLogin, test.orderNumber).Return(test.mockErrorReturn)
				srv.mockBalanceRepository.EXPECT().CreateDefaultBalance(test.userLogin)
				srv.mockOrderRepository.EXPECT().CommitTransaction(gomock.Any())

				body := strconv.FormatInt(test.orderNumber, 10)
				req := createPostRequest(
					srv.GetFullPath(CreateOrderHandlerPath),
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

func TestNewCreateOrder_WithIncorrectOrderNumber(t *testing.T) {
	srv := startTestServer(t)
	defer srv.server.Close()
	incorrectOrderNumber := "4539319503436461"

	req := createPostRequest(
		srv.GetFullPath(CreateOrderHandlerPath),
		incorrectOrderNumber)
	req.SetCookie(
		&http.Cookie{
			Name:  security.AuthorizationTokenName,
			Value: registrationAndAuthentication(t, srv, "userLogin"),
		},
	)

	resp, _ := req.Send()

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode())
}
