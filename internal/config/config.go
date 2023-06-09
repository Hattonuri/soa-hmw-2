package config

import "github.com/kelseyhightower/envconfig"

type Client struct {
	User     string `default:"defaultUser" envconfig:"USER"`
	Role     string `default:"defaultRole" envconfig:"ROLE"`
	Hostname string `default:"localhost" envconfig:"HOSTNAME"`
}

func InitClient(c *Client) error {
	return envconfig.Process("", c)
}

type Server struct {
	MaxPlayers uint64 `default:"4" envconfig:"MAX_PLAYERS"`
}

func InitServer(c *Server) error {
	return envconfig.Process("", c)
}
