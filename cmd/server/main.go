package main

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/hattonuri/soa-hmw-2/internal/adapter/grpc_server"
	"github.com/hattonuri/soa-hmw-2/internal/config"
	"github.com/hattonuri/soa-hmw-2/internal/generated/proto/base"
	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	s := grpc.NewServer()
	cfg := &config.Server{}
	if err := config.InitServer(cfg); err != nil {
		log.Fatalf("error parsing from env: %v", err)
	}
	srv := grpc_server.NewServer(cfg)
	base.RegisterMafiaServiceServer(s, srv)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Start server")
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
