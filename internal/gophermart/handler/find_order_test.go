package handler

import (
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestFindAllOrderHandler(t *testing.T) {
	srv := startTestServer(t)
	defer srv.server.Close()

	tests := []struct {
		name            string
		userLogin       string
		mockErrorReturn interface{}
		ordersReturn    []model.OrderEntity
		wantCode        int
	}{
		{
			name:      "Successful searching of all orders",
			userLogin: "ilya",
			ordersReturn: []model.OrderEntity{
				model.OrderEntity{
					OrderNumber: 79927398713,
					Status:      "PROCESSED",
					Accrual:     Float32Ptr(200.21),
					UserLogin:   "ilya",
					CreatedAt:   time.Now(),
				},
			},
			mockErrorReturn: nil,
			wantCode:        http.StatusOK,
		},
		{
			name:            "Successful searching of all orders (no content)",
			userLogin:       "ilya",
			ordersReturn:    nil,
			mockErrorReturn: postgres.ErrOrdersNotExist,
			wantCode:        http.StatusNoContent,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			{
				srv.mockOrderRepository.EXPECT().FindAllOrders(test.userLogin).Return(test.ordersReturn, test.mockErrorReturn)
				req := createGetRequest(
					srv.GetFullPath(GetOrderHandlerPath),
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

func Float32Ptr(i float32) *float32 {
	return &i
}
