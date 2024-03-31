package server

import (
	"backendService/internals/setup/config"
	"backendService/internals/setup/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

type server struct {
	Config *config.AppConfig
	Db     *database.Database
	App    *gin.Engine
}

func NewServer(config *config.AppConfig, db *database.Database) *server {
	return &server{
		Config: config,
		Db:     db,
		App:    gin.Default(),
	}
}

func (s *server) Start() {
	address := ":" + strconv.Itoa(config.Config.App.Port)
	s.App.Run(address)
}
