package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/logging"
	"github.com/hellcats88/abstracte/storage"
	"github.com/hellcats88/abstracte/tenant"
	"github.com/hellcats88/rem/runtime"
)

type runtimeHandler struct {
}

func (g runtimeHandler) createRuntimeContext(ctx *gin.Context) {
	// logging key is always populated, don't check the exist return value.
	iLogCtx, _ := ctx.Get(api.LogKey)
	logCtx := iLogCtx.(logging.Context)

	// tx key is always populated, don't check the exist return value.
	iTxCtx, _ := ctx.Get(api.TxKey)
	txCtx := iTxCtx.(storage.Transaction)

	// tenant key is always populated, don't check the exist return value.
	iTenantCtx, _ := ctx.Get(api.TenantKey)
	tenantCtx := iTenantCtx.(tenant.Context)

	rCtx := runtime.New(logCtx, txCtx, tenantCtx)
	ctx.Set(api.RuntimeKey, rCtx)
	ctx.Next()
}
