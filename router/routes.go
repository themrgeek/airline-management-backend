package router

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/themrgeek/airline-management-backend/controller"
)

// SetupRoutes configures and returns the router with all routes and middleware
func SetupRoutes(db *sql.DB) http.Handler {
	r := mux.NewRouter()

	// Setup middleware
	r.Use(loggingMiddleware)
	r = setupCORS(r)

	// Setup routes
	setupBasicRoutes(r, db)

	// Wrap with standard logging
	return handlers.LoggingHandler(log.Writer(), r)
}

// loggingMiddleware logs request details
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"Method: %s Path: %s IP: %s Duration: %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// setupCORS adds CORS middleware
func setupCORS(r *mux.Router) *mux.Router {
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.Use(handlers.CORS(headers, methods, origins))
	return r
}

// setupBasicRoutes adds all routes to the router
func setupBasicRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controller.RegisterUser(w, r, db)
	}).Methods(http.MethodPost)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.LoginUser(w, r, db)
	}).Methods(http.MethodPost)
}

// homeHandler handles the root endpoint
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Welcome to the Airline Management System API"}`))
}

// StartServer starts the HTTPS server with the provided configuration
func StartServer(handler http.Handler, certFile, keyFile string) error {
	server := &http.Server{
		Addr:         ":443",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting HTTPS server on port 443")
	return server.ListenAndServeTLS(certFile, keyFile)
}
