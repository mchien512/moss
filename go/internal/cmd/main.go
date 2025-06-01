package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	entryApp "moss/go/internal/app/entry"
	linkApp "moss/go/internal/app/link"
	pbentry "moss/go/internal/genproto/entry"
	pblink "moss/go/internal/genproto/link"
	"moss/go/internal/interceptors"
	"moss/go/internal/repository/db"
	entryRepo "moss/go/internal/repository/entry"
	linkRepo "moss/go/internal/repository/link"
	entryService "moss/go/internal/service/entry"
	"moss/go/internal/service/link"
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
	service := entryService.NewService(app)

	linkRepo := linkRepo.NewRepository(dbConn)
	linkApp := linkApp.NewApp(linkRepo)
	linkService := link.NewService(linkApp)

	// Start listening on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server with interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryServerInterceptor()),
	)

	// Register the EntryService
	pbentry.RegisterEntryServiceServer(grpcServer, service)
	// Register the LinkService
	pblink.RegisterLinkServiceServer(grpcServer, linkService)

	reflection.Register(grpcServer)

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
