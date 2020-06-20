package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/handler"
)

// AccountResolver resolves the account type.
type AccountResolver struct {
	user *handler.User
}

// NewAccount creates a resolver for the account you are currently logged in as.
func NewAccount(ctx context.Context) (*AccountResolver, error) {
	user := ctx.Value("user").(*handler.User)
	return &AccountResolver{user: user}, nil
}

// Email resolves the email address of the account.
func (r *AccountResolver) Email() *string {
	if r.user == nil {
		return nil
	}
	return &r.user.Email
}

// ID resolves the ID used to uniquely identify the account.
func (r *AccountResolver) ID() *string {
	if r.user == nil {
		return nil
	}
	return &r.user.ID
}
