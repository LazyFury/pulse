package framework

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	*gin.RouterGroup
}

func Group(router *gin.Engine, path string) *RouterGroup {
	return &RouterGroup{RouterGroup: router.Group(path)}
}

func (r *RouterGroup) Doc(doc *DocItem) *RouterGroup {
	// register handler
	// r.Use(middlewares...).Handle(doc.Method, doc.Path, func(ctx *gin.Context) {
	// 	// 似乎还可以在这里添加验证和直接填参数 handler(ctx,page,limit,params...)
	// 	handler(ctx)
	// })

	// write doc
	basePath := reflect.ValueOf(r).Elem().FieldByName("basePath")
	doc.Path = basePath.String() + doc.Path

	if doc.Body != nil {
		doc.Body = reflectSwagSpec(doc.Body)
	}

	if doc.Params != nil {
		doc.Params = reflectSwagSpec(doc.Params)
	}

	DocInstance.Add(doc)
	return r
}

// reflectSwagSpec reflects the struct and get the json tag and swag tag, and make it as a map.
// it will return a map like this:
//
//	{
//		"key1": {
//			"type": "string",
//			"label": "label1",
//			"required": true
//		},
//		"key2": {
//			"type": "integer",
//			"label": "label2",
//			"required": false
//		},
//		...
//	}
func reflectSwagSpec[T any](v T) map[string]map[string]any {
	var _map = make(map[string]map[string]any)
	var fieldLenth = reflect.TypeOf(v).NumField()
	var fields []reflect.StructField

	for i := 0; i < fieldLenth; i++ {
		fields = append(fields, reflect.TypeOf(v).Field(i))
	}
	for _, field := range fields {
		specstr := field.Tag.Get("swag")
		specs := strings.Split(specstr, ",")

		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}

		if _map[key] == nil {
			_map[key] = make(map[string]any)
		}

		for _, spec := range specs {
			switch spec {
			case "required":
				_map[key]["required"] = true
			case "string":
				_map[key]["type"] = "string"
			case "int":
				_map[key]["type"] = "integer"

			case "float64":
				_map[key]["type"] = "number"

			case "interface{}":
				_map[key]["type"] = "object"

			case "array":
				_map[key]["type"] = "array"

			default:
				_map[key]["label"] = spec
			}
		}

	}
	return _map

}

// GROUP
func (r *RouterGroup) Group(path string) *RouterGroup {
	return &RouterGroup{RouterGroup: r.RouterGroup.Group(path)}
}

// USE
func (r *RouterGroup) Use(middleware ...gin.HandlerFunc) *RouterGroup {
	r.RouterGroup.Use(middleware...)
	return &RouterGroup{RouterGroup: r.RouterGroup}
}

// GET
func (r *RouterGroup) GET(path string, handlers ...gin.HandlerFunc) *RouterGroup {
	r.RouterGroup.GET(path, handlers...)
	return r
}

// POST
func (r *RouterGroup) POST(path string, handlers ...gin.HandlerFunc) *RouterGroup {
	r.RouterGroup.POST(path, handlers...)
	return r
}

// PUT
func (r *RouterGroup) PUT(path string, handlers ...gin.HandlerFunc) *RouterGroup {
	r.RouterGroup.PUT(path, handlers...)
	return r
}

// DELETE
func (r *RouterGroup) DELETE(path string, handlers ...gin.HandlerFunc) *RouterGroup {
	r.RouterGroup.DELETE(path, handlers...)
	return r
}

// PATCH
func (r *RouterGroup) PATCH(path string, handlers ...gin.HandlerFunc) *RouterGroup {
	r.RouterGroup.PATCH(path, handlers...)
	return r
}

// Handle
func (r *RouterGroup) Handle(method, path string, handlers ...gin.HandlerFunc) *RouterGroup {
	r.RouterGroup.Handle(method, path, handlers...)
	return r
}
