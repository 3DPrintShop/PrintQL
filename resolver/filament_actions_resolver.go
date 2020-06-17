package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
)

// FilamentActionsResolver is a resolver for actions that can be taken on filament.
type FilamentActionsResolver struct {
}

// NewFilamentActionsResolver creates a new resolver for mutating filament.
func NewFilamentActionsResolver() (*FilamentActionsResolver, error) {
	resolver := FilamentActionsResolver{}
	return &resolver, nil
}

// CreateFilamentBrand creates a new brand of filament.
func (r FilamentActionsResolver) CreateFilamentBrand(ctx context.Context, args createFilamentBrandArgs) (*graphql.ID, error) {
	client := ctx.Value("client").(*printdb.Client)

	filamentBrandID, err := client.CreateFilamentBrand(args.Name)
	brandID := graphql.ID(filamentBrandID)

	return &brandID, err
}

type createFilamentSpoolArgs struct {
	BrandID graphql.ID
}

// CreateFilamentSpool creates a new spool of filament.
func (r FilamentActionsResolver) CreateFilamentSpool(ctx context.Context, args createFilamentSpoolArgs) (*FilamentSpoolResolver, error) {
	var errs errors.Errors

	client := ctx.Value("client").(*printdb.Client)

	spoolID, err := client.CreateFilamentSpool(string(args.BrandID))
	if err != nil {
		errs = append(errs, err)
	}

	resolver, err := NewFilamentSpool(ctx, NewFilamentSpoolArgs{ID: spoolID})
	if err != nil {
		errs = append(errs, err)
	}

	return resolver, errs.Err()
}
