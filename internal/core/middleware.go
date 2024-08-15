package core

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
)

func ContentTypeSetter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func RequestIdHeader(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqId := middleware.GetReqID(r.Context())
		w.Header().Set("Request_Id", reqId)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
