package upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyfury/pulse/framework"
	"github.com/lazyfury/pulse/helper/response"
)

type UploadPlugin struct {
	*framework.Plugin
	uploader *Uploader
}

var _ framework.IPlugin = (*UploadPlugin)(nil)

func NewUploadPlugin() *UploadPlugin {
	return &UploadPlugin{
		Plugin: framework.NewPlugin(),
		uploader: &Uploader{
			BaseDir:      "./static/upload",
			UploadMethod: DefaultUpload,
			GetFile:      DefaultGetFile,
		},
	}
}

func (p *UploadPlugin) RegisterRouter(router *framework.RouterGroup) {
	router.Doc(&framework.DocItem{
		Method: http.MethodPost,
		Path:   "/upload",
	}).POST("/upload", p.upload)
}

func (p *UploadPlugin) upload(ctx *gin.Context) {
	path, err := p.uploader.Default(ctx.Request)
	if err != nil {
		response.Error(response.ErrBadRequest, err.Error()).Done(ctx)
		return
	}
	response.Success(map[string]any{
		"url": path,
	}).Done(ctx)
}
