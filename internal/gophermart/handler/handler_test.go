package handler

import (
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/mocks"
	"github.com/faust8888/gophermart/internal/gophermart/route"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
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

	h := New(cfg, withUserRepositoryMock(mockUserRepository))
	router := route.New(h)

	return &testServer{
		server:             httptest.NewServer(router),
		mockUserRepository: mockUserRepository,
	}
}

func createPostRequest(url string, body interface{}, headers ...requestHeader) *resty.Request {
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = url
	req.Body = body
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

func withUserRepositoryMock(mockRepository service.UserRepository) Option {
	return func(h *Handler) {
		h.registerService = service.NewRegisterService(mockRepository)
		h.loginService = service.NewLoginService(mockRepository)
	}
}

type testServer struct {
	server             *httptest.Server
	mockUserRepository *mocks.MockUserRepository
}

func (s *testServer) GetFullPath(handlerURI string) string {
	return s.server.URL + handlerURI
}

type requestHeader struct {
	HeaderName  string
	HeaderValue string
}
