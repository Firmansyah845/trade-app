package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"awesomeProjectCr/internal/handler"

	"github.com/go-chi/httprate"
	"github.com/rs/zerolog/log"
)

type Func func(handler http.Handler) http.Handler

func Recover() Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					response := handler.ErrorResponse{
						Message: "please contact admin",
						Code:    http.StatusInternalServerError,
					}
					log.Error().
						Interface("panic", err).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("stack", string(debug.Stack())).
						Msg("panic recovered")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(response)
					return
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimiter() Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httprate.Limit(1, time.Minute,
				httprate.WithKeyFuncs(
					func(r *http.Request) (string, error) {
						token := r.Header.Get("Authorization")
						return token, nil
					},
				),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					response := handler.ErrorResponse{
						Message: "too many request",
						Code:    http.StatusTooManyRequests,
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusTooManyRequests)
					_ = json.NewEncoder(w).Encode(response)
				}),
			)
			next.ServeHTTP(w, r)
		})
	}
}
