package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/Duncanian/iam-gateway/server"
	pb "github.com/Duncanian/iam-gateway/server/protobuf"
	"github.com/Duncanian/iam-gateway/server/resolvers"
	"google.golang.org/grpc"
)

func grpcConnection() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}

	defer conn.Close()

	client := pb.NewSimpleServerClient(conn)

	err, _ = client.CreateUser(context.Background(), &pb.GoogleIdToken{Token: "tevgvybububvvg"})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Println("Created user!")
}

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	grpcConnection()

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(server.NewExecutableSchema(server.Config{Resolvers: &resolvers.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
