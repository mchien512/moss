package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	entryApp "moss/go/internal/app/entry"
	pb "moss/go/internal/genproto/entry"
	"moss/go/internal/interceptors"
	"moss/go/internal/repository/db"
	entryRepo "moss/go/internal/repository/entry"
	entryService "moss/go/internal/service/entry"
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
	pb.RegisterEntryServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
