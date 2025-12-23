package middlewares

import (
	"net/http"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do not allow browsers to prefetch DNS
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		// deny page inside <iframe>
		w.Header().Set("X-Frame-Options", "DENY")

		// deprecated https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-XSS-Protection
		// w.Header().Set("X-XSS-Protection", "1,mode=block")

		// don't try to determine content type, trust what is sent
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// allow only HTTPS
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// dont allow content from different source
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		next.ServeHTTP(w, r)
	})
}
