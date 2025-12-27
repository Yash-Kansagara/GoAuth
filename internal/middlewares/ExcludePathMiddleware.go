package middlewares

import (
	"net/http"
	"strings"
)

// middleware = middleware which will be wrapped with a path check
// paths = prefix sring which will be matched as string prefix
// returns middleware
func ExcludePath(middleware func(next http.Handler) http.Handler, paths ...string) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqPath := r.URL.Path
			for _, path := range paths {
				if strings.HasPrefix(reqPath, path) {
					// if path matches, skipp calling middleware and call the next middleware
					next.ServeHTTP(w, r)
					return
				}
			}
			// if no path matches execute intented middleware
			middleware(next).ServeHTTP(w, r)
		})

	}
}
