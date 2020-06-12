package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
)

// FilamentBrandResolver resolves the filament brand type.
type FilamentBrandResolver struct {
	filamentBrand printdb.FilamentBrand
}

// NewFilamentBrandsArgs are the arguments passed in to get a page of filament brands.
type NewFilamentBrandsArgs struct {
	ID *string
}

// NewFilamentBrandArgs are the arguments passed in to get a specific filament brand.
type NewFilamentBrandArgs struct {
	ID string
}

// NewFilamentBrand creates a resolver for a filament brand specified by the passed in args.
func NewFilamentBrand(ctx context.Context, args NewFilamentBrandArgs) (*FilamentBrandResolver, error) {
	filamentBrand, errs := loader.LoadFilamentBrand(ctx, args.ID)

	return &FilamentBrandResolver{filamentBrand: filamentBrand}, errs
}

// NewFilamentBrands gets a list of filaments specified by the args, if no args are passed in all are returned.
func NewFilamentBrands(ctx context.Context, args NewFilamentBrandsArgs) (*[]*FilamentBrandResolver, error) {
	if args.ID != nil {
		filamentBrand, err := NewFilamentBrand(ctx, NewFilamentBrandArgs{ID: *args.ID})
		resolvers := []*FilamentBrandResolver{filamentBrand}
		return &resolvers, err
	}

	var resolvers []*FilamentBrandResolver
	var errs errors.Errors
	filamentBrands, err := loader.LoadFilamentBrands(ctx, "")

	if err != nil {
		return nil, err
	}

	for _, filamentBrandID := range filamentBrands.FilamentBrandIDs {
		filamentBrand, err := NewFilamentBrand(ctx, NewFilamentBrandArgs{ID: filamentBrandID})
		if err != nil {
			errs = append(errs, err)
		}
		resolvers = append(resolvers, filamentBrand)
	}

	return &resolvers, errs.Err()
}

// ID returns the ID of the filament brand.
func (r *FilamentBrandResolver) ID() graphql.ID {
	return graphql.ID(r.filamentBrand.ID)
}

// Name returns the name of the filament brand.
func (r *FilamentBrandResolver) Name() string {
	return r.filamentBrand.Name
}

// PurchaseLink returns a url to purchase the filament.
func (r *FilamentBrandResolver) PurchaseLink() *string {
	return &r.filamentBrand.PurchaseLink
}

// StartWeight returns the weight that a spool of the filament starts at.
func (r *FilamentBrandResolver) StartWeight() *int32 {
	return &r.filamentBrand.StartWeight
}

// SpoolWeight is the weight of the spool that the filament comes on.
func (r *FilamentBrandResolver) SpoolWeight() *int32 {
	return &r.filamentBrand.StartWeight
}
