package main

import (
	"context"
	"log"
	"net"

	"github.com/Sna-ken/videoweb/config"
	handler "github.com/Sna-ken/videoweb/internal/auth"
	"github.com/Sna-ken/videoweb/internal/auth/repository"
	authservice "github.com/Sna-ken/videoweb/internal/auth/service"
	authserver "github.com/Sna-ken/videoweb/kitex_gen/auth/authservice"
	"github.com/Sna-ken/videoweb/kitex_gen/user/userservice"
	"github.com/Sna-ken/videoweb/pkg/db/mysql"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/server"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.Load("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	mysqldb, err := mysql.NewMySQL(&cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}

	userClient, err := userservice.NewClient(
		"user",
		client.WithHostPorts(cfg.RPC.UserAddress),
	)
	if err != nil {
		log.Fatal(err)
	}

	address, err := net.ResolveTCPAddr("tcp", cfg.RPC.AuthAddress)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	authdb := repository.NewAuthDB(mysqldb, redisClient)

	authService := authservice.NewAuthService(authdb)
	authHandler := handler.NewAuthServiceImpl(authService, userClient)
	svr := authserver.NewServer(authHandler, server.WithServiceAddr(address))

	if err := svr.Run(); err != nil {
		log.Println(err.Error())
	}
}
