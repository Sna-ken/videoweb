package client

import (
	"github.com/Sna-ken/videoweb/kitex_gen/auth/authservice"
	"github.com/cloudwego/kitex/client"
)

var AuthServiceClient authservice.Client

func InitAuthServiceClient(addr string) error {
	cli, err := authservice.NewClient(
		"auth",
		client.WithHostPorts(addr),
	)
	if err != nil {
		return err
	}

	AuthServiceClient = cli
	return nil
}
