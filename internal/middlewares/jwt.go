package middlewares

import (
	"net/http"
	"os"

	utils "github.com/Yash-Kansagara/GoAuth/internal/Utils"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// read cookie with jwt token
		token, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// parse jwt token string
		_, err = jwt.Parse(token.Value, keyFunc, parseOptions()...)
		if utils.WriteIfError(w, err, "Authorization Error", http.StatusUnauthorized) {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func keyFunc(t *jwt.Token) (any, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func parseOptions() []jwt.ParserOption {
	return []jwt.ParserOption{
		// allow only valid signing methods
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	}
}
