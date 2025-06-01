package main

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	entryApp "moss/go/internal/app/entry"
	linkApp "moss/go/internal/app/link"
	entryconnect "moss/go/internal/genproto/protobuf/entry/entryconnect"
	linkconnect "moss/go/internal/genproto/protobuf/link/linkconnect"
	"moss/go/internal/repository/db"
	entryRepo "moss/go/internal/repository/entry"
	linkRepo "moss/go/internal/repository/link"
	entryService "moss/go/internal/service/entry"
	linkService "moss/go/internal/service/link"
)

func main() {
	// Initialize database
	config := &db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "moss_db",
		SSLMode:  "disable",
	}

	dbConn, err := db.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize layers
	repo := entryRepo.NewRepository(dbConn)
	app := entryApp.NewApp(repo)
	entrySvc := entryService.NewService(app)

	linkRepo := linkRepo.NewRepository(dbConn)
	linkApp := linkApp.NewApp(linkRepo)
	linkSvc := linkService.NewService(linkApp)

	// Create Connect adapters for your services
	entryServicePath, entryConnectSvc := entryconnect.NewEntryServiceHandler(
		entrySvc,
		connect.WithInterceptors(
		// Add your interceptors here
		),
	)

	linkServicePath, linkConnectSvc := linkconnect.NewLinkServiceHandler(
		linkSvc,
		connect.WithInterceptors(
		// Add your interceptors here
		),
	)

	// Set up CORS middleware
	corsMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Connect-Protocol-Version")
			w.Header().Set("Access-Control-Max-Age", "3600")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle(entryServicePath, entryConnectSvc)
	mux.Handle(linkServicePath, linkConnectSvc)

	// Use h2c to support HTTP/2 without TLS
	handler := corsMiddleware(mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}

	log.Printf("Connect server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
