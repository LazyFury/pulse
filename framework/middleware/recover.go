package middleware

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/lazyfury/pulse/helper/response"
)

func RecoverHandlerFunc(ctx *gin.Context) {
	defer func(ctx *gin.Context) {
		if err := recover(); err != nil {
			msg := ""
			if e, ok := err.(error); ok {
				msg = e.Error()
			}
			log.Error(err)
			response.Error(response.ErrInternalError, msg).Done(ctx)
			ctx.Abort()
			panic(err)
		}
	}(ctx)

	ctx.Next()
}
