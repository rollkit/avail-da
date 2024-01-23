package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rollkit/avail-da"
	"github.com/rollkit/go-da/proxy"
)

func main() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Error reading config file:", err)
	}
	// Parse the configuration data into a Config struct
	var config avail.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("Error parsing config file:", err)
	}
	ctx := context.Background()
	da := avail.NewAvailDA(config, ctx)
	srv := proxy.NewServer(da, grpc.Creds(insecure.NewCredentials()))
	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalln("failed to create network listener:", err)
	}
	log.Println("serving avail-da over gRPC on:", lis.Addr())
	err = srv.Serve(lis)
	if !errors.Is(err, grpc.ErrServerStopped) {
		log.Fatalln("gRPC server stopped with error:", err)
	}
}
