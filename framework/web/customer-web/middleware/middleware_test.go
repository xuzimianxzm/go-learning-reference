package middleware

import (
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	var r = NewRouter()
	r.Use(logger)
	r.Use(timeout)
	r.Use(ratelimit)
	r.Add("/", http.HandlerFunc(hello))
}

func logger(http.Handler) http.Handler {
	return http.HandlerFunc(nil)
}

func timeout(http.Handler) http.Handler {
	return http.HandlerFunc(nil)
}

func ratelimit(http.Handler) http.Handler {
	return http.HandlerFunc(nil)
}

func hello(wr http.ResponseWriter, r *http.Request) {
	wr.Write([]byte("hello"))
}
