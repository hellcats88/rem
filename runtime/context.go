package runtime

import (
	"github.com/hellcats88/abstracte/env"
	"github.com/hellcats88/abstracte/logging"
	"github.com/hellcats88/abstracte/runtime"
	"github.com/hellcats88/abstracte/storage"
	"github.com/hellcats88/abstracte/tenant"
)

// Context defines a runtime group of common information
type context struct {
	log    logging.Context
	tx     storage.Transaction
	tenant tenant.Context
	env    env.Context
}

func (c context) Log() logging.Context {
	return c.log
}

func (c context) Tx() storage.Transaction {
	return c.tx
}

func (c context) Tenant() tenant.Context {
	return c.tenant
}

func (c context) Env() env.Context {
	return c.env
}

// New creates an instance of runtime context with the associated logging info
func New(log logging.Context, tx storage.Transaction, tenant tenant.Context) runtime.Context {
	return context{
		log:    log,
		tx:     tx,
		tenant: tenant,
		env:    env.New("Global"),
	}
}

// New creates an instance of runtime context with the associated logging info for a specific environment
func NewWithEnv(log logging.Context, tx storage.Transaction, tenant tenant.Context, env env.Context) runtime.Context {
	return context{
		log:    log,
		tx:     tx,
		tenant: tenant,
		env:    env,
	}
}

// New clones a context using new transaction
func NewFromTx(from runtime.Context, tx storage.Transaction) runtime.Context {
	return context{
		log:    from.Log,
		env:    from.Env,
		tx:     tx,
		tenant: from.Tenant,
	}
}
