package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/mohamedfawas/rmshop-admin-service/internal/config"
	"github.com/mohamedfawas/rmshop-admin-service/internal/repository"
	"github.com/mohamedfawas/rmshop-admin-service/internal/service"
	adminv1 "github.com/mohamedfawas/rmshop-proto/gen/v1/admin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository and service
	adminRepo := repository.NewAdminRepository(db)
	adminService := service.NewAdminService(adminRepo)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	adminv1.RegisterAdminServiceServer(grpcServer, adminService)
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Handle shutdown gracefully
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Received shutdown signal")
		grpcServer.GracefulStop()
	}()

	log.Printf("Starting admin service on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
