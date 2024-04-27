package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/syahyudi09/ChatPintar/App/config"
	"github.com/syahyudi09/ChatPintar/App/delivery/controller"
	"github.com/syahyudi09/ChatPintar/App/delivery/websocket"
	"github.com/syahyudi09/ChatPintar/App/manager"
)

type Server interface {
	Run()
}

type serverImpl struct {
	engine  *gin.Engine
	usecase manager.UseceaseManager
	config  config.Config
}

func (s *serverImpl) Run() {

	controller.NewAuthController(s.engine, s.usecase.GetAuthUsecase())
	controller.NewChatGroupController(s.engine, s.usecase.GetChatGroupUsecase())
	websocket.NewWebSocketController(s.engine, s.usecase.GetPrivateUsecase(), s.usecase.GetAuthUsecase())

	s.engine.Run(":8080")
}

func NewServer() Server {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	infra := manager.NewInfraManager(config)
	repo := manager.NewRepoManager(infra)
	usecase := manager.NewUsecaseManager(repo)

	return &serverImpl{
		engine:  r,
		usecase: usecase,
		config:  config,
	}
}
