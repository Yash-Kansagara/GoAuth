package server

import "net/http"

// mux to handle all routes
func GetRootMux() *http.ServeMux {
	mux := http.NewServeMux()

	// root handler
	mux.HandleFunc("/{$}", handleRoot)

	// add other routes
	RegisterSignupHandler(mux)

	return mux
}
