package middlewares

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	utils "github.com/Yash-Kansagara/GoAuth/internal/Utils"
)

/*
rate limits requests based on client id (ip, name, token etc)
*/
func RateLimiterMiddleware(next http.Handler) http.Handler {
	reqPerSec, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_RATE"))
	burst, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BURST"))
	rateLimiter := utils.NewRateLimiter(reqPerSec, burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			if rateLimiter.Allow(host) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests) // code 429
			}
		} else {
			// log it but allow request anyway, can be restricted
			log.Println("RateLimiter: Error getting host from remote address", r.RemoteAddr, err)
			next.ServeHTTP(w, r)
		}
	})
}
