package tenant

import "github.com/hellcats88/abstracte/tenant"

type context struct {
	id     string
	userID string
}

func New(id string, userID string) tenant.Context {
	return context{id: id, userID: userID}
}

func NewEmpty() tenant.Context {
	return context{id: "Global", userID: "system"}
}

func (c context) ID() string {
	return c.id
}

func (c context) UserID() string {
	return c.userID
}
