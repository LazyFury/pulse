package framework

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	Engine *gin.Engine
	Config *viper.Viper
}

func NewHttpServer(config *viper.Viper, engine *gin.Engine) *HttpServer {

	return &HttpServer{
		Engine: engine,
		Config: config,
	}
}

func (h *HttpServer) Start(port string) *gin.Engine {
	if port == "0" {
		port = "8900"
	}
	log.Debug(fmt.Sprintf("start http server on port: http://127.1:%s", port))
	h.Engine.Run(fmt.Sprintf(":%s", port))
	return h.Engine
}
