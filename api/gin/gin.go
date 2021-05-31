package gin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/logging"
	"github.com/hellcats88/abstracte/runtime"
)

type ginHttp struct {
	engine *gin.Engine
	log    logging.Logger
}

func New(log logging.Logger) api.Http {
	engine := gin.Default()

	entity := ginHttp{
		engine: engine,
		log:    log,
	}

	return entity
}

func reverse(items []gin.HandlerFunc) []gin.HandlerFunc {
	for i := 0; i < len(items)/2; i++ {
		j := len(items) - i - 1
		items[i], items[j] = items[j], items[i]
	}
	return items
}

func (g *ginHttp) wrapService(service api.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ignore exist result because runtime context is mandatory
		// and the user cannot remove it
		ctx, _ := c.Get(api.RuntimeKey)
		rCtx := ctx.(runtime.Context)

		output := service(rCtx, ginServiceInput{ctx: c})
		c.Set(api.ServiceResultKey, output)
	}
}

func (g ginHttp) AddRoute(method string, path string, config api.Config, service api.Service) error {
	var handlers []gin.HandlerFunc
	ginCnf := config.(ginConfig)

	handlers = append(handlers, ginCnf.log, ginCnf.tenant, ginCnf.tx)

	if ginCnf.headers != nil {
		handlers = append(handlers, ginCnf.headers)
	}

	if ginCnf.model != nil {
		handlers = append(handlers, ginCnf.model)
	}

	if ginCnf.params != nil {
		handlers = append(handlers, ginCnf.params)
	}

	if ginCnf.queryParams != nil {
		handlers = append(handlers, ginCnf.queryParams)
	}

	if ginCnf.beforeRun != nil && len(ginCnf.beforeRun) > 0 {
		handlers = append(handlers, ginCnf.beforeRun...)
	}

	handlers = append(handlers, runtimeHandler{}.createRuntimeContext, g.wrapService(service))

	//reverse order due to recursive logic of gin middlewares
	handlers = append(handlers, resultHandler{log: g.log}.handleResult)

	if ginCnf.afterRun != nil && len(ginCnf.afterRun) > 0 {
		handlers = append(handlers, reverse(ginCnf.afterRun)...)
	}

	if ginCnf.commit != nil {
		handlers = append(handlers, ginCnf.commit)
	}

	g.engine.Handle(method, path, handlers...)
	return nil
}

func (g ginHttp) Listen(port int, address string) error {
	return g.engine.Run(fmt.Sprintf("%s:%d", address, port))
}
