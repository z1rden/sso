package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"sso/internal/sso/auth"
	"sso/internal/sso/core"
	"sso/pkg/config"
	"sso/pkg/databases/postgres"
	"sso/pkg/logger"
	pb_sso "sso/pkg/sso"
)

func main() {
	f, err := os.OpenFile("sso.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	l := logger.New(f)
	ctx := context.Background()
	cfg := config.MustLoad()

	server := grpc.NewServer()
	sso := core.NewService(auth.NewService(postgres.NewDatabase(ctx, cfg, l), l), l)
	pb_sso.RegisterSSOServer(server, sso)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		l.Error("failed to listen: %v", err)
		os.Exit(1)
	}

	go func() {
		conn, err := grpc.NewClient(":50052",
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			l.Error("failed to connect: %v", err)
			os.Exit(1)
		}

		client := pb_sso.NewSSOClient(conn)
		res, err := client.Login(ctx, &pb_sso.LoginRequest{
			Email:    "admin@example.com",
			Password: "admin",
			AppId:    1,
		})
		if err != nil {
			l.Error("failed to register: %v", err)
		}
		l.Info("Registered user %d", res.Token)
	}()

	if err := server.Serve(lis); err != nil {
		l.Error("failed to serve: %v", err)
		os.Exit(1)
	}

}
