package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/runtime"
)

type inputParamsHandler struct {
	requestedInputParams []string
}

func (g inputParamsHandler) loadParams(ctx *gin.Context) {
	params := make(map[string]string)

	rCtx, _ := ctx.Get(api.RuntimeKey)
	svcCtx := rCtx.(runtime.Context)

	for _, p := range g.requestedInputParams {
		pV := ctx.Param(p)
		if pV == "" {
			ctx.AbortWithStatusJSON(http.StatusNotFound, api.Model{
				Error: api.ErrorModel{
					Code:   api.ApiErrorMissingRequiredItem,
					Msg:    "Missing part of URL",
					DevMsg: fmt.Sprintf("Cannot find parameter %s", p),
					CorrId: svcCtx.Log().CorrID(),
				},
			})
		}

		params[p] = pV
	}

	ctx.Set(api.InputParamsKey, params)
	ctx.Next()
}
