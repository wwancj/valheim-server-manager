package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"valheim-server-manager/app"
	"valheim-server-manager/app/server"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create app instance
	application := app.NewApp()
	application.Startup(nil) // No Wails context needed

	// Create WebSocket hub and wire it as the event emitter
	hub := server.NewHub()
	go hub.Run()
	application.SetEventEmitter(hub)

	// Create handlers
	handlers := server.NewHandlers(application, hub)

	// Create router with API routes
	router := server.NewRouter(handlers)

	// Serve static frontend files (React SPA)
	distFS, _ := fs.Sub(assets, "frontend/dist")
	fileServer := http.FileServer(http.FS(distFS))
	router.ServeSPA(fileServer)

	// Add middleware
	handler := server.CORSMiddleware(server.LoggingMiddleware(router))

	// Start server
	addr := ":13256"
	log.Printf("Valheim Server Manager starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
