package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/logging"
	"github.com/hellcats88/rem/tenant"
)

type tenantHandler struct {
	log logging.Logger
}

func (g tenantHandler) createNoTenant(ctx *gin.Context) {
	ctx.Set(api.TenantKey, tenant.NewEmpty())
	ctx.Next()
}

func (g tenantHandler) createTenantFromHeaders(ctx *gin.Context) {
	// logging key is always populated, don't check the exist return value.
	iLogCtx, _ := ctx.Get(api.LogKey)
	logCtx := iLogCtx.(logging.Context)
	tCtx := tenant.New(ctx.GetHeader("X-Tenant-ID"), ctx.GetHeader("X-Tenant-UserID"))

	if tCtx.ID() == "" || tCtx.UserID() == "" {
		g.log.Error(logCtx, "Rejected request caused by missing tenant informations")

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.Model{
			Error: api.ErrorModel{
				Code:   api.ApiErrorAuthFailed,
				Msg:    "Failed to get user information",
				DevMsg: "Missing X-Tenant-ID or X-Tenant-UserID headers",
				CorrId: logCtx.CorrID(),
			},
		})
	}

	ctx.Set(api.TenantKey, tCtx)
	ctx.Next()
}
