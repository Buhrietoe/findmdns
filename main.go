package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Buhrietoe/findmdns/api"
	"github.com/Buhrietoe/findmdns/discovery"
)

func main() {
	// Check if 'serve' subcommand is used
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		// Run HTTP server mode
		runHTTPServer()
	} else {
		// Run JSON output mode (default behavior)
		runJSONOutput()
	}
}

func runJSONOutput() {
	// Create discovery manager
	manager := discovery.NewManager()

	// Perform scan
	log.Println("Performing scan...")
	err := manager.Scan()
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	// Get devices and output as JSON
	devices := manager.GetDevices()
	
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(devices)
	if err != nil {
		log.Fatalf("Failed to output JSON: %v", err)
	}
}

func runHTTPServer() {
	// Create discovery manager
	manager := discovery.NewManager()

	// Perform initial scan
	log.Println("Performing initial scan...")
	err := manager.Scan()
	if err != nil {
		log.Printf("Initial scan failed: %v", err)
	}

	// Create API handler
	handler := api.NewHandler(manager)

	// Set up routes
	router := api.SetupRoutes(handler)

	// Add static file serving
	staticPath := filepath.Join(os.Getenv("PWD"), "static")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	// For all other routes, serve the UI files (this enables SPA routing)
	uiPath := filepath.Join(os.Getenv("PWD"), "ui")
	
	// Create a custom handler for non-API routes
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If it's an API route, it should already be handled by the router
		if r.URL.Path != "/" && (r.URL.Path[:5] == "/api/" || r.URL.Path[:7] == "/static/") {
			http.NotFound(w, r)
			return
		}
		
		// For all other paths, serve index.html to enable SPA routing
		// This allows JavaScript to handle the routing properly
		http.ServeFile(w, r, filepath.Join(uiPath, "index.html"))
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HTTP server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}