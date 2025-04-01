package framework

type Handler interface {
	Register(router *RouterGroup)
	GetRouterGroupName() string
}
