package server

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Yash-Kansagara/GoAuth/internal/middlewares"
	"golang.org/x/net/http2"
)

// start http server
func Start() {

	apiPort := os.Getenv("API_PORT")
	TLSEnable, _ := strconv.ParseBool(os.Getenv("TLS_ENABLE"))

	// setup http server
	mux := GetRootMux()
	handler := ApplyMiddlewares(mux)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", apiPort),
		Handler: handler,
	}
	defer srv.Close()

	// check TLS config
	if TLSEnable {
		srv.TLSConfig = getTLSConfig()
	}

	// configure server
	err := http2.ConfigureServer(srv, &http2.Server{})
	if err != nil {
		log.Fatal("Failed Confuguring http2 server", err)
	}

	certFile := os.Getenv("TLS_CERT_FILEPATH")
	keyFile := os.Getenv("TLS_KEY_FILEPATH")
	// start HTTP(S) server
	if TLSEnable {
		err = srv.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			log.Fatal("Failed starting HTTPS server", err)
		}
	} else {
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatal("Failed starting HTTP server", err)
		}
	}

}

func ApplyMiddlewares(mux *http.ServeMux) http.Handler {
	// applied/runs from bottom
	handler := middlewares.CompressionMiddleware(mux)
	handler = middlewares.SecurityHeaders(handler)
	handler = middlewares.RateLimiterMiddleware(handler)
	// can move config somewhere else
	handler = middlewares.HPPMiddleware(handler, &middlewares.HppMiddlewareConfig{
		VerifyBodyForm:      true,
		VerifyQuery:         true,
		AllowedBodyFormKeys: []string{"name", "id"},
		AllowedQueryKeys:    []string{"sortBy", "sortOrder", "id", "name"},
	})
	handler = middlewares.Cors(handler)

	return handler
}

func getTLSConfig() *tls.Config {
	conf := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	return conf
}

func logReqDetails(req *http.Request) {
	log.Println(req.Method, req.URL.Path)
}

func handleRoot(res http.ResponseWriter, req *http.Request) {
	logReqDetails(req)
	io.WriteString(res, "Works\n")
}
