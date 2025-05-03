package security

import (
	"fmt"
	appHandler "github.com/faust8888/gophermart/internal/gophermart/handler"
	"net/http"
)

func NewMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != appHandler.UserRegisterHandlerPath {
			fmt.Printf("Request Path: %s\n", path)
		}
		handler.ServeHTTP(w, r)
	})
}
