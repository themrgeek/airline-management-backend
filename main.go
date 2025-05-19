package main

import (
	"log"
	"net/http"

	"github.com/themrgeek/airline-management-backend/config"
	"github.com/themrgeek/airline-management-backend/routes"
)

func main() {
	// Load configuration
	config.LoadEnv()
	config.ConnectDB()

	// Initialize router
	router := routes.InitializeRouter()

	// Configure HTTPS server
	srv := &http.Server{
		Addr:    ":443",
		Handler: router,
	}

	// Redirect HTTP to HTTPS
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS)); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Start HTTPS server
	if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Fatalf("HTTPS server error: %v", err)
	}
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}
