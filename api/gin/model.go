package gin

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/runtime"
)

type modelHandler struct {
	requestedModel interface{}
}

func (g modelHandler) getModel(ctx *gin.Context) {
	emptyModel := reflect.New(reflect.TypeOf(g.requestedModel)).Interface()

	rCtx, _ := ctx.Get(api.RuntimeKey)
	svcCtx := rCtx.(runtime.Context)

	if err := ctx.ShouldBindJSON(emptyModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.Model{
			Error: api.ErrorModel{
				Code:   api.ApiErrorUnexpected,
				Msg:    "API needs a valid payload",
				DevMsg: fmt.Sprintf("Failed to transform payload model from JSON. %v", err),
				CorrId: svcCtx.Log().CorrID(),
			},
		})
	}

	ctx.Set(api.InputModelKey, emptyModel)
	ctx.Next()
}
