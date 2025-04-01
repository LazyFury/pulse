package framework

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lazyfury/pulse/framework/types"
)

type Doc struct {
	Title   string
	URL     string
	Version string
	Items   []*DocItem
}

// MarshalJSON implements MarshalJSON.
// 可能需要整理成 openapi 适合的结构，暂时不处理
func (d *Doc) MarshalJSON() ([]byte, error) {
	type Alias Doc
	var a = &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	return json.Marshal(a)
}

var _ types.MarshalJSON = (*Doc)(nil)

type DocItem struct {
	Title       string
	Content     string
	Description string
	Tag         string
	Path        string
	Method      string
	Params      any
	Body        any
	Success     any
	Failure     any
}

var DocInstance *Doc = nil
var mutex sync.RWMutex

func NewDoc() *Doc {
	DocInstance = &Doc{}
	return DocInstance
}

func (d *Doc) Add(item *DocItem) {
	mutex.Lock() // i`m not sure
	defer mutex.Unlock()
	d.Items = append(d.Items, item)
}

func (d *Doc) Handler(ctx *gin.Context) {
	mutex.RLock()
	defer mutex.RUnlock()
	ctx.JSON(http.StatusOK, d)
}
