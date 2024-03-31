package server

import (
	"backendService/internals/setup/config"
	"backendService/internals/setup/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	Server *server
)

type server struct {
	Config *config.AppConfig
	Db     *database.Database
	App    *gin.Engine
}

func NewServer(config *config.AppConfig, db *database.Database) *server {
	Server = &server{
		Config: config,
		Db:     db,
		App:    gin.Default(),
	}
	return Server
}

func (s *server) Start() {
	address := ":" + strconv.Itoa(config.Config.App.Port)
	s.App.Run(address)
}
