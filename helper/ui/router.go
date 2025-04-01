package ui

import (
	"encoding/json"

	"github.com/lazyfury/pulse/framework/types"
)

type IRouter interface {
	GetPath() string
	GetParams() map[string]string
	GetComponentName() string
}

type Router struct {
	Key         string            `json:"key"`
	Title       string            `json:"title"`
	Icon        string            `json:"icon"`
	Path        string            `json:"path"`
	Params      map[string]string `json:"params"`
	Page        IPage             `json:"page"`
	Children    []IRouter         `json:"children"`
	Description string            `json:"description"`
	Apis        map[string]string `json:"apis"`
}

var _ IRouter = (*Router)(nil)
var _ types.MarshalJSON = (*Router)(nil)

func (r *Router) GetPath() string {
	return r.Path
}

func (r *Router) GetParams() map[string]string {
	return r.Params
}

func (r *Router) GetComponentName() string {
	return r.Page.GetComponentName()
}

func (r *Router) MarshalJSON() ([]byte, error) {
	type Alias Router
	type A struct {
		*Alias
		Component string `json:"component"`
	}
	var a *A = &A{
		Alias:     (*Alias)(r),
		Component: r.Page.GetComponentName(),
	}
	return json.Marshal(a)
}

func (r *Router) AddChildren(children ...IRouter) *Router {
	r.Children = append(r.Children, children...)
	return r
}

func (r *Router) WithDescription(description string) *Router {
	r.Description = description
	return r
}

func (r *Router) WithApis(apis map[string]string) *Router {
	r.Apis = apis
	return r
}

func NewRouter(key, title, icon, path string, params map[string]string, page IPage) *Router {
	return &Router{
		Key:    key,
		Title:  title,
		Icon:   icon,
		Path:   path,
		Params: params,
		Page:   page,
	}
}
