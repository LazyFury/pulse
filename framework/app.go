package framework

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/lazyfury/pulse/framework/constants"
	"github.com/lazyfury/pulse/framework/driver"
	"github.com/spf13/viper"
)

type IApp interface {
	GetName() string
	GetDB() *driver.DB
	Migrate()
	GetMiddlewares() []gin.HandlerFunc
	GetConfig() *viper.Viper
	Register(router *RouterGroup)
	GetHandlers() []Handler
	GetPlugins() []IPlugin
}

type App struct {
	AppName     string
	AppConfig   *viper.Viper
	db          *driver.DB
	Handlers    []Handler
	Plugins     []IPlugin
	Middlewares []gin.HandlerFunc
}

func NewApp(name string, config *viper.Viper, db *driver.DB, handlers []Handler) *App {
	return &App{
		AppName:   name,
		AppConfig: config,
		db:        db,
		Handlers:  handlers,
	}
}

// WithPlugins
func (a *App) WithPlugins(plugins []IPlugin) *App {
	a.Plugins = append(a.Plugins, plugins...)
	return a
}

// with plugin
func (a *App) WithPlugin(plugin IPlugin) *App {
	a.Plugins = append(a.Plugins, plugin)
	return a
}

// with middlewares
func (a *App) WithMiddlewares(middlewares []gin.HandlerFunc) *App {
	a.Middlewares = append(a.Middlewares, middlewares...)
	return a
}

func (a *App) GetConfig() *viper.Viper {
	return a.AppConfig
}

func (a *App) GetName() string {
	return a.AppName
}

func (a *App) GetDB() *driver.DB {
	return a.db
}

func (a *App) GetHandlers() []Handler {
	return a.Handlers
}

func (a *App) GetPlugins() []IPlugin {
	return a.Plugins
}

func (a *App) GetMiddlewares() []gin.HandlerFunc {
	return a.Middlewares
}

type RegisterApp struct {
	Router *RouterGroup
	App    IApp
}

func (instance *RegisterApp) Register() {
	log.Debug("====App:【" + instance.App.GetName() + "】 register Start====================================")
	log.Debug("注册app", "app", instance.App.GetName())
	if instance.App.GetConfig().GetBool(constants.DEBUG) {
		log.Debug("迁移中")
		instance.App.Migrate()
		log.Debug("迁移成功 success")
	}
	router := instance.Router.Group("")
	router.Use(instance.App.GetMiddlewares()...)
	instance.App.Register(router)

	for _, handler := range instance.App.GetHandlers() {
		handler.Register(router.Group(handler.GetRouterGroupName()))
	}

	for _, plugin := range instance.App.GetPlugins() {
		plugin.RegisterRouter(router.Group(""))
	}

	log.Debug("注册成功", "app", instance.App.GetName())
	log.Debug("====App:【" + instance.App.GetName() + "】 register End====================================")
	log.Debug("\n")
}

func RegisterApps(apps []*RegisterApp) {
	for _, app := range apps {
		app.Register()
	}
}
