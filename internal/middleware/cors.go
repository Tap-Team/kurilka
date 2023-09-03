package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func CORS(h http.Handler) http.Handler {
	return cors.AllowAll().Handler(h)
}
