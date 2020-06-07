package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/handler"
)

type AccountResolver struct {
	user *handler.User
}

func NewAccount(ctx context.Context) (*AccountResolver, error) {
	user := ctx.Value("user").(*handler.User)
	return &AccountResolver{user: user}, nil
}

func (r *AccountResolver) Email() *string {
	if r.user == nil {
		return nil
	}
	return &r.user.Email
}

func (r *AccountResolver) ID() *string {
	if r.user == nil {
		return nil
	}
	return &r.user.Id
}
