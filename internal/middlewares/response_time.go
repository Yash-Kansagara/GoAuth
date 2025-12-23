package middlewares

import (
	"log"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("ResponseTimeMiddleware->")
		start := time.Now()
		next.ServeHTTP(w, r)
		diff := time.Since(start)
		log.Println("RT:", diff.Milliseconds())

	})
}
