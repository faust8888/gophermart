package handler

import (
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/mocks"
	"github.com/faust8888/gophermart/internal/gophermart/route"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startTestServer(t *testing.T) *testServer {
	ctrl := gomock.NewController(t)

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	mockUserRepository := mocks.NewMockUserRepository(ctrl)
	mockOrderRepository := mocks.NewMockOrderRepository(ctrl)
	mockBalanceRepository := mocks.NewMockBalanceRepository(ctrl)
	mockWithdrawRepository := mocks.NewMockWithdrawRepository(ctrl)

	h := New(cfg, withRepositoryMocks(
		mockOrderRepository,
		mockBalanceRepository,
		mockUserRepository,
		mockWithdrawRepository,
		cfg))
	router := route.New(h)

	return &testServer{
		server:                 httptest.NewServer(router),
		mockUserRepository:     mockUserRepository,
		mockOrderRepository:    mockOrderRepository,
		mockBalanceRepository:  mockBalanceRepository,
		mockWithdrawRepository: mockWithdrawRepository,
	}
}

func createPostRequest(url string, body interface{}, headers ...RequestHeader) *resty.Request {
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = url
	req.Body = body
	for _, header := range headers {
		req.SetHeader(header.HeaderName, header.HeaderValue)
	}
	return req
}

func createGetRequest(url string, headers ...RequestHeader) *resty.Request {
	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = url
	for _, header := range headers {
		req.SetHeader(header.HeaderName, header.HeaderValue)
	}
	return req
}

func getTokenFromResponse(res *resty.Response) string {
	for _, cookie := range res.Cookies() {
		if cookie.Name == security.AuthorizationTokenName {
			return cookie.Value
		}
	}
	return ""
}

func registrationAndAuthentication(t *testing.T, srv *testServer, login string) string {
	ctx := gomock.Any()
	registerBody := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, login, "password")
	srv.mockUserRepository.EXPECT().CreateUser(ctx, login, "password").Return(nil)
	resp, _ := createPostRequest(srv.GetFullPath(UserRegisterHandlerPath), registerBody).Send()
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	loginBody := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, login, "password")
	srv.mockUserRepository.EXPECT().CheckUser(ctx, login, "password").Return(nil)
	resp, _ = createPostRequest(srv.GetFullPath(UserLoginHandlerPath), loginBody).Send()
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	return getTokenFromResponse(resp)
}

func withRepositoryMocks(
	orderMockRepository service.OrderRepository,
	balanceMockRepository service.BalanceRepository,
	userMockRepository service.UserRepository,
	withdrawMockRepository service.WithdrawRepository,
	cfg *config.Config) Option {
	balanceService := service.NewBalanceService(balanceMockRepository)
	return func(h *Handler) {
		h.findOrderService = service.NewOrderService(orderMockRepository, balanceService, cfg)
		h.orderService = service.NewOrderService(orderMockRepository, balanceService, cfg)
		h.registerService = service.NewRegisterService(userMockRepository)
		h.loginService = service.NewLoginService(userMockRepository)
		h.withdrawService = service.NewWithdrawService(withdrawMockRepository, orderMockRepository)
		h.balanceService = balanceService
	}
}

type testServer struct {
	server                 *httptest.Server
	mockUserRepository     *mocks.MockUserRepository
	mockOrderRepository    *mocks.MockOrderRepository
	mockBalanceRepository  *mocks.MockBalanceRepository
	mockWithdrawRepository *mocks.MockWithdrawRepository
}

func (s *testServer) GetFullPath(handlerURI string) string {
	return s.server.URL + handlerURI
}

type RequestHeader struct {
	HeaderName  string
	HeaderValue string
}
