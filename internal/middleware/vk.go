package middleware

import "net/http"

func VK(next http.Handler) http.Handler {
	return next
}
