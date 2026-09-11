package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"shipment/internal/application"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "shipment/gen"
	infraKafka "shipment/internal/infra/kafka"
	infraPostgres "shipment/internal/infra/postgres"
	grpcTransport "shipment/internal/transport/grpc"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/shipments?sslmode=disable")
	grpcPort := getEnv("GRPC_PORT", "50051")
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	kafkaTopic := getEnv("KAFKA_TOPIC", "shipment-events")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
	log.Println("connected to database")

	shipmentRepo := infraPostgres.NewShipmentRepo(db)
	eventRepo := infraPostgres.NewShipmentEventRepo(db)
	eventPublisher := infraKafka.NewShipmentEventPublisher(kafkaBrokers, kafkaTopic)
	defer eventPublisher.Close()

	shipmentService := application.NewShipmentService(shipmentRepo, eventRepo, eventPublisher)

	handler := grpcTransport.NewShipmentHandler(shipmentService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShipmentServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
