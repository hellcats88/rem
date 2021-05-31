package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/logging"
	"github.com/hellcats88/abstracte/runtime"
)

type resultHandler struct {
	log logging.Logger
}

func (g resultHandler) handleResult(ctx *gin.Context) {
	ctx.Next()

	result, _ := ctx.Get(api.ServiceResultKey)
	svcRes := result.(api.ServiceOutput)

	rCtx, _ := ctx.Get(api.RuntimeKey)
	svcCtx := rCtx.(runtime.Context)

	if svcRes.Status() != api.ApiErrorNoError {
		httpCode := http.StatusInternalServerError

		switch svcRes.Status() {
		case api.ApiErrorAuthFailed:
			httpCode = http.StatusForbidden
		case api.ApiErrorEntityAlreadyExists:
			httpCode = http.StatusConflict
		case api.ApiErrorEntityDoesNotExists:
			httpCode = http.StatusNotFound
		case api.ApiErrorMissingRequiredItem:
			httpCode = http.StatusBadRequest
		case api.ApiErrorUnexpected:
			httpCode = http.StatusInternalServerError
		case api.ApiErrorUnknownItemRequested:
			httpCode = http.StatusBadRequest
		}

		ctx.AbortWithStatusJSON(httpCode, api.Model{
			Error: api.ErrorModel{
				Code:   svcRes.Status(),
				Msg:    svcRes.ErrMessage(),
				DevMsg: svcRes.Err().Error(),
				CorrId: svcCtx.Log().CorrID(),
			},
		})

	} else {
		ctx.JSON(http.StatusOK, api.Model{
			Error: api.ErrorModel{
				Code:   api.ApiErrorNoError,
				CorrId: svcCtx.Log().CorrID(),
			},
			Data: svcRes.ResponseModel(),
		})
	}
}
