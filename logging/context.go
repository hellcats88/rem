package logging

import (
	"github.com/google/uuid"
	"github.com/hellcats88/abstracte/logging"
)

type context struct {
	corrID string
	extra  []logging.K
}

// NewContext creates new logger context from a correlation ID
func NewContext(corrID string) logging.Context {
	return &context{
		corrID: corrID,
	}
}

// NewContextUUID creates new logger context generating
// a new correlation ID from UUID package
func NewContextUUID() logging.Context {
	return &context{
		corrID: uuid.New().String(),
	}
}

// CorrID returns the context correlation ID
func (ctx context) CorrID() string {
	return ctx.corrID
}

// CloneEmpty create a new independent context with the same original correlation ID without extra parameters
func (ctx context) CloneNoExtra() logging.Context {
	return NewContext(ctx.CorrID())
}

// Clone create a new independent context with the same original correlation ID and extra parameters
func (ctx context) Clone() logging.Context {
	return &context{
		corrID: ctx.corrID,
		extra:  ctx.extra,
	}
}

// AddExtra appends to the internal memory the input list of extra parameters. Duplicates are admitted
func (ctx *context) AddExtra(extras ...logging.K) {
	ctx.extra = append(ctx.extra, extras...)
}

// GetExtras returns the list of stored extras in this context
func (ctx context) Extras() []logging.K {
	return ctx.extra
}
