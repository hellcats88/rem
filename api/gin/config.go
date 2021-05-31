package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/logging"
)

type ginConfig struct {
	log         gin.HandlerFunc
	tenant      gin.HandlerFunc
	tx          gin.HandlerFunc
	commit      gin.HandlerFunc
	headers     gin.HandlerFunc
	model       gin.HandlerFunc
	params      gin.HandlerFunc
	queryParams gin.HandlerFunc
	beforeRun   []gin.HandlerFunc
	afterRun    []gin.HandlerFunc
}

func (g ginConfig) Valid() bool {
	return true
}

type ginConfigBuilder struct {
	config ginConfig
	log    logging.Logger
}

func NewConfigBuilder(log logging.Logger) api.ConfigBuilder {
	return &ginConfigBuilder{
		log: log,
		config: ginConfig{
			log:    logHandler{}.createLogContext,
			tx:     transactionHandler{log: log}.createNoTransaction,
			tenant: tenantHandler{log: log}.createNoTenant,
		},
	}
}

func (b *ginConfigBuilder) Log(p api.ConfigLog) api.ConfigBuilder {
	b.config.log = logHandler{}.createLogContext
	return b
}

func (b *ginConfigBuilder) CustomLog(p api.C) api.ConfigBuilder {
	b.config.log = p.Handler.(gin.HandlerFunc)
	return b
}

func (b *ginConfigBuilder) Tenant(p api.ConfigTenant) api.ConfigBuilder {
	if p == api.ConfigTenantFromHeaders {
		b.config.tenant = tenantHandler{log: b.log}.createTenantFromHeaders
	}
	return b
}

func (b *ginConfigBuilder) CustomTenant(p api.C) api.ConfigBuilder {
	b.config.tenant = p.Handler.(gin.HandlerFunc)
	return b
}

func (b *ginConfigBuilder) Tx(p api.ConfigTx) api.ConfigBuilder {
	if p == api.ConfigTxManaged {
		b.config.tx = transactionHandler{log: b.log}.createManagedTransaction
		b.config.commit = transactionHandler{log: b.log}.createCommitTx
	} else if p == api.ConfigTxUnmanaged {
		b.config.tx = transactionHandler{log: b.log}.createUnmanagedTransaction
	}
	return b
}

func (b *ginConfigBuilder) CustomTx(p api.C) api.ConfigBuilder {
	b.config.tx = p.Handler.(gin.HandlerFunc)
	return b
}

func (b *ginConfigBuilder) Headers(p interface{}) api.ConfigBuilder {
	b.config.headers = headersHandler{requestedModel: p, log: b.log}.loadHEaders
	return b
}

func (b *ginConfigBuilder) CustomHeaders(p api.C) api.ConfigBuilder {
	b.config.headers = p.Handler.(gin.HandlerFunc)
	return b
}

func (b *ginConfigBuilder) InputModel(p interface{}) api.ConfigBuilder {
	b.config.model = modelHandler{requestedModel: p}.getModel
	return b
}

func (b *ginConfigBuilder) CustomInputModel(p api.C) api.ConfigBuilder {
	b.config.model = p.Handler.(gin.HandlerFunc)
	return b
}

func (b *ginConfigBuilder) InputParams(name []string) api.ConfigBuilder {
	b.config.params = inputParamsHandler{requestedInputParams: name}.loadParams
	return b
}

func (b *ginConfigBuilder) CustomInputParam(p api.C) api.ConfigBuilder {
	b.config.params = p.Handler.(gin.HandlerFunc)
	return b
}

func (b *ginConfigBuilder) QueryParams(p interface{}) api.ConfigBuilder {
	b.config.queryParams = queryParamsHandler{requestedModel: p, log: b.log}.getQueryParams
	return b
}

func (b *ginConfigBuilder) CustomQueryParams(p api.C) api.ConfigBuilder {
	b.config.queryParams = p.Handler.(gin.HandlerFunc)
	return b
}

func (b *ginConfigBuilder) CustomBeforeRun(p api.C) api.ConfigBuilder {
	b.config.beforeRun = append(b.config.beforeRun, p.Handler.(gin.HandlerFunc))
	return b
}

func (b *ginConfigBuilder) CustomAfterRun(p api.C) api.ConfigBuilder {
	b.config.afterRun = append(b.config.beforeRun, p.Handler.(gin.HandlerFunc))
	return b
}

func (b *ginConfigBuilder) Build() api.Config {
	return b.config
}
