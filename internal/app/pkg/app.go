package app

//создает подключение к бд, веб-сервер, создает конфиг. может создаваться и стартоваться
import (
	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/handler"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config  *config.Config
	Logger  *logrus.Logger
	Router  *gin.Engine
	Handler *handler.Handler
}

func (a *Application) Run() {
	a.Logger.Info("Server start up")
	a.Handler.Register(a.Router)
	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	if err := a.Router.Run(serverAddress); err != nil {
		a.Logger.Fatalln(err)
	}

	a.Logger.Info("Server down")
}

func New(c *config.Config, r *gin.Engine, l *logrus.Logger, h *handler.Handler) *Application {
	return &Application{
		Config:  c,
		Logger:  l,
		Router:  r,
		Handler: h,
	}
}
