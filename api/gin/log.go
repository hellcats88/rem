package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	alog "github.com/hellcats88/abstracte/logging"
	"github.com/hellcats88/rem/logging"
)

type logHandler struct {
}

func (g logHandler) createLogContext(ctx *gin.Context) {
	var lCtx alog.Context

	if corrId := ctx.GetHeader("X-Correlation-ID"); corrId != "" {
		lCtx = logging.NewContext(corrId)
	} else {
		lCtx = logging.NewContextUUID()
	}

	ctx.Set(api.LogKey, lCtx)
	ctx.Next()
}
