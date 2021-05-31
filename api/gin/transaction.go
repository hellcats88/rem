package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellcats88/abstracte/api"
	"github.com/hellcats88/abstracte/logging"
	"github.com/hellcats88/abstracte/runtime"
	"github.com/hellcats88/abstracte/storage"
)

type txNoOp struct{}

func (txNoOp) Ref() interface{}                    { return nil }
func (txNoOp) Commit() error                       { return nil }
func (txNoOp) Rollback() error                     { return nil }
func (txNoOp) Begin() (storage.Transaction, error) { return nil, nil }

type transactionHandler struct {
	log logging.Logger
	db  storage.Context
}

func (g transactionHandler) createNoTransaction(ctx *gin.Context) {
	ctx.Set(api.TxKey, txNoOp{})
	ctx.Next()
}

func (g transactionHandler) createManagedTransaction(ctx *gin.Context) {
	tx, err := g.db.Tx()

	rCtx, _ := ctx.Get(api.RuntimeKey)
	svcCtx := rCtx.(runtime.Context)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.Model{
			Error: api.ErrorModel{
				Code:   api.ApiErrorUnexpected,
				Msg:    "Failed to open new managed transaction",
				DevMsg: err.Error(),
				CorrId: svcCtx.Log().CorrID(),
			},
		})

		return
	}

	ctx.Set(api.TxKey, tx)
	ctx.Next()
}

func (g transactionHandler) createUnmanagedTransaction(ctx *gin.Context) {
	tx, err := g.db.UnmanagedTx()

	rCtx, _ := ctx.Get(api.RuntimeKey)
	svcCtx := rCtx.(runtime.Context)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.Model{
			Error: api.ErrorModel{
				Code:   api.ApiErrorUnexpected,
				Msg:    "Failed to open new unmanaged transaction",
				DevMsg: err.Error(),
				CorrId: svcCtx.Log().CorrID(),
			},
		})

		return
	}

	ctx.Set(api.TxKey, tx)
	ctx.Next()
}

func (g transactionHandler) createCommitTx(ctx *gin.Context) {
	ctx.Next()

	result, _ := ctx.Get(api.ServiceResultKey)
	svcRes := result.(api.ServiceOutput)

	tx, _ := ctx.Get(api.TxKey)
	svcTx := tx.(storage.Transaction)

	rCtx, _ := ctx.Get(api.RuntimeKey)
	svcCtx := rCtx.(runtime.Context)

	if svcRes.Status() != api.ApiErrorNoError {
		err := svcTx.Rollback()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.Model{
				Error: api.ErrorModel{
					Code:   api.ApiErrorUnexpected,
					Msg:    "Failed to rollback changes",
					DevMsg: err.Error(),
					CorrId: svcCtx.Log().CorrID(),
				},
			})
			return
		}

	} else {
		err := svcTx.Commit()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.Model{
				Error: api.ErrorModel{
					Code:   api.ApiErrorUnexpected,
					Msg:    "Failed to commit changes",
					DevMsg: err.Error(),
					CorrId: svcCtx.Log().CorrID(),
				},
			})
			return
		}
	}
}
