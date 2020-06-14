package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
)

type FilamentActionsResolver struct {
}

func NewFilamentActionsResolver() (*FilamentActionsResolver, error) {
	resolver := FilamentActionsResolver{}
	return &resolver, nil
}

func (r FilamentActionsResolver) CreateFilamentBrand(ctx context.Context, args createFilamentBrandArgs) (*graphql.ID, error) {
	client := ctx.Value("client").(*printdb.Client)

	filamentBrandID, err := client.CreateFilamentBrand(args.Name)
	brandId := graphql.ID(filamentBrandID)

	return &brandId, err
}

type createFilamentSpoolArgs struct {
	BrandID graphql.ID
}

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
