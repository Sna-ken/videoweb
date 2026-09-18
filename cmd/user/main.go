package main

import (
	"log"
	"net"

	"github.com/Sna-ken/videoweb/config"
	handler "github.com/Sna-ken/videoweb/internal/user"
	"github.com/Sna-ken/videoweb/internal/user/repository"
	userservice "github.com/Sna-ken/videoweb/internal/user/service"
	userserver "github.com/Sna-ken/videoweb/kitex_gen/user/userservice"
	"github.com/Sna-ken/videoweb/pkg/db/mysql"
	"github.com/cloudwego/kitex/server"
)

func main() {
	cfg, err := config.Load("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	db, err := mysql.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}

	address, err := net.ResolveTCPAddr("tcp", cfg.RPC.UserAddress)
	if err != nil {
		log.Fatal(err)
	}

	userdb := repository.NewUserDB(db)
	userService := userservice.NewUserService(userdb)
	userHandler := handler.NewUserServiceImpl(userService)
	svr := userserver.NewServer(userHandler, server.WithServiceAddr(address))

	if err := svr.Run(); err != nil {
		log.Println(err.Error())
	}
}
