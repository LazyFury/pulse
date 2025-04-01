package framework

import "github.com/gin-gonic/gin"

type PluginInfo struct {
	Name        string
	Description string
	Version     string
	Author      string
	Website     string
}

type IPlugin interface {
	Install() error
	Uninstall() error

	// 注册路由
	RegisterRouter(router *RouterGroup)
	GetMiddlewares() []gin.HandlerFunc

	// GetPluginInfo 获取插件信息
	GetPluginInfo() *PluginInfo

	// GetPluginName 获取插件名称
	GetPluginName() string

	// Invoke 执行插件方法,event 在系统中定义，比如 用户已登录，用户已下单，用户已付款等
	Invoke(event string, data map[string]interface{})
}

type Plugin struct {
	Info        *PluginInfo
	Middlewares []gin.HandlerFunc
}

// Invoke implements IPlugin.
func (p *Plugin) Invoke(event string, data map[string]interface{}) {
	panic("unimplemented")
}

// GetPluginInfo implements IPlugin.
func (p *Plugin) GetPluginInfo() *PluginInfo {
	return p.Info
}

// GetPluginName implements IPlugin.
func (p *Plugin) GetPluginName() string {
	return p.Info.Name
}

// Middlewares implements IPlugin.
func (p *Plugin) GetMiddlewares() []gin.HandlerFunc {
	return p.Middlewares
}

// RegisterRouter implements IPlugin.
func (p *Plugin) RegisterRouter(router *RouterGroup) {

}

func (p *Plugin) Install() error {
	return nil
}

func (p *Plugin) Uninstall() error {
	return nil
}

// with middlewares
func (p *Plugin) WithMiddlewares(middlewares []gin.HandlerFunc) *Plugin {
	p.Middlewares = append(p.Middlewares, middlewares...)
	return p
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

var _ IPlugin = (*Plugin)(nil)

type IPluginManager interface {
	Install() error
	Uninstall() error
}

type PluginManager struct {
	Plugins []IPlugin
}

func NewPluginManager(plugins ...IPlugin) *PluginManager {
	return &PluginManager{
		Plugins: plugins,
	}
}

func (p *PluginManager) Install() error {
	for _, plugin := range p.Plugins {
		if err := plugin.Install(); err != nil {
			return err
		}
	}
	return nil
}

func (p *PluginManager) Uninstall() error {
	for _, plugin := range p.Plugins {
		if err := plugin.Uninstall(); err != nil {
			return err
		}
	}
	return nil
}

type PayMethod struct {
	Name string
}

type IPaymentPlugin interface {
	GetPayMethods() []PayMethod
	Pay(data map[string]interface{}) error
}

type PaymentPlugin struct {
	*Plugin
}

func NewPaymentPlugin() *PaymentPlugin {
	return &PaymentPlugin{
		Plugin: NewPlugin(),
	}
}

func (p *PaymentPlugin) GetPayMethods() []PayMethod {
	return []PayMethod{
		{
			Name: "alipay",
		},
	}
}

func (p *PaymentPlugin) Pay(data map[string]interface{}) error {
	return nil
}
