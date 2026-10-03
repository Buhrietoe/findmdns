package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/Buhrietoe/findmdns/internal/api"
	"github.com/Buhrietoe/findmdns/internal/discovery"
)

//go:embed ui/index.html static/*
var assets embed.FS

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printUsage()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "serve" {
		serveFlags := flag.NewFlagSet("serve", flag.ExitOnError)
		portFlag := serveFlags.String("port", "", "HTTP server port (overrides PORT env var)")
		serveFlags.Parse(os.Args[2:])
		runHTTPServer(*portFlag)
	} else {
		runJSONOutput()
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `findmdns - mDNS service discovery tool

Usage:
  findmdns [flags]              Scan and output devices as JSON to stdout
  findmdns serve [--port PORT]  Start HTTP server with web UI

Global Flags:
  -h, --help    Show usage information

Serve Flags:
  -port string  HTTP server port (overrides PORT env var, default 8080)
`)
}

func runJSONOutput() {
	manager := discovery.NewManager()

	log.Println("Performing scan...")
	err := manager.Scan()
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	devices := manager.GetDevices()

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(devices)
	if err != nil {
		log.Fatalf("Failed to output JSON: %v", err)
	}
}

func runHTTPServer(portFlag string) {
	manager := discovery.NewManager()

	log.Println("Performing initial scan...")
	err := manager.Scan()
	if err != nil {
		log.Printf("Initial scan failed: %v", err)
	}

	handler := api.NewHandler(manager)
	router := api.SetupRoutes(handler)

	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		log.Fatalf("Failed to create static filesystem: %v", err)
	}
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	uiIndex, err := assets.ReadFile("ui/index.html")
	if err != nil {
		log.Fatalf("Failed to read index.html: %v", err)
	}
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 5 && (r.URL.Path[:5] == "/api/" || r.URL.Path[:8] == "/static/") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(uiIndex)
	})

	port := portFlag
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HTTP server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
