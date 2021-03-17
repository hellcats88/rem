package env

import (
	"os"

	"github.com/hellcats88/abstracte/env"
)

type context struct {
	name string
}

func New(name string) env.Context {
	return context{name: name}
}

func NewFromEnvVar(name string) env.Context {
	return context{name: os.Getenv(name)}
}

func (c context) Name() string {
	return c.name
}
