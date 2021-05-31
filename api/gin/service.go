package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
)

type ginServiceInput struct {
	ctx *gin.Context
}

func (g ginServiceInput) RawCtx() interface{} {
	return g.ctx
}

func (g ginServiceInput) Model() interface{} {
	model, ok := g.ctx.Get(api.InputModelKey)
	if !ok {
		panic("Missing required Input Model. Is pipeline correct?")
	}

	return model
}

func (g ginServiceInput) InputParams() map[string]string {
	model, ok := g.ctx.Get(api.InputParamsKey)
	if !ok {
		panic("Missing required Input Params. Is pipeline correct?")
	}

	return model.(map[string]string)
}

func (g ginServiceInput) QueryParams() interface{} {
	model, ok := g.ctx.Get(api.QueryParamsKey)
	if !ok {
		panic("Missing required Query Params. Is pipeline correct?")
	}

	// Query params are optional, but the model must be stored in the context.
	// Users should check if the model is valid using default value of the
	// model.
	return model
}

func (g ginServiceInput) Headers() interface{} {
	model, ok := g.ctx.Get(api.HeadersModelKey)
	if !ok {
		panic("Missing required Headers. Is pipeline correct?")
	}

	return model
}
