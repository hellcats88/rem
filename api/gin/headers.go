package gin

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/logging"
)

type headersHandler struct {
	requestedModel interface{}
	log            logging.Logger
}

func (g headersHandler) loadHEaders(ctx *gin.Context) {
	// logging key is always populated, don't check the exist return value.
	iLogCtx, _ := ctx.Get(api.LogKey)
	logCtx := iLogCtx.(logging.Context)

	emptyModel := reflect.New(reflect.TypeOf(g.requestedModel)).Interface()

	if err := ctx.ShouldBindHeader(emptyModel); err != nil {
		g.log.Warn(logCtx, "No headers found to be parsed. %v", err)
	}

	ctx.Set(api.HeadersModelKey, emptyModel)
	ctx.Next()
}
