package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
)

type createFilamentBrandArgs struct {
	Name string
}

// CreateFilamentBrand is a resolver that creates a new brand of filament then returns that filament object.
func (r SchemaResolver) CreateFilamentBrand(ctx context.Context, args createFilamentBrandArgs) (*graphql.ID, error) {
	client := ctx.Value("client").(*printdb.Client)

	filamentBrandID, err := client.CreateFilamentBrand(args.Name)
	brandID := graphql.ID(filamentBrandID)

	return &brandID, err
}
