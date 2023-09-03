package middleware

import (
	"errors"
	"mime"
	"net/http"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
)

func JSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				httphelpers.Error(w, errors.New("Malformed Content-Type header"))
				return
			}

			if mt != "application/json" {
				httphelpers.Error(w, errors.New("Content-Type header must be application/json"))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
