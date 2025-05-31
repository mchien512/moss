package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"lumo/internal/app/entry"
	pb "lumo/internal/genproto/entry"
	"lumo/internal/interceptors"
	"lumo/internal/repository/db"
	entryRepo "lumo/internal/repository/entry"
	entryService "lumo/internal/service/entry"
	"net"
)

func main() {
	// Initialize database
	config := &db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "lumo_db",
		SSLMode:  "disable",
	}

	dbConn, err := db.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize layers
	repo := entryRepo.NewRepository(dbConn)
	app := entry.NewApp(repo)
	service := entryService.NewService(app)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server with interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryServerInterceptor()),
	)

	// Register service
	pb.RegisterEntryServiceServer(s, service)

	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
