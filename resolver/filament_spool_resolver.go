package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
)

type FilamentSpoolResolver struct {
	filamentSpool *printdb.FilamentSpool
}

// NewFilamentSpoolsArgs are the arguments passed in to get a page of filament spools.
type NewFilamentSpoolsArgs struct {
	ID *string
}

// NewFilamentSpoolArgs are the arguments passed in to get a specific filament spool.
type NewFilamentSpoolArgs struct {
	ID string
}

// NewFilamentSpool creates a resolver for a filament spool specified by the passed in args.
func NewFilamentSpool(ctx context.Context, args NewFilamentSpoolArgs) (*FilamentSpoolResolver, error) {
	filamentSpool, errs := loader.LoadFilamentSpool(ctx, args.ID)

	return &FilamentSpoolResolver{filamentSpool: &filamentSpool}, errs
}

// NewFilamentSpools gets a list of spools specified by the args, if no args are passed in all spools are returned.
func NewFilamentSpools(ctx context.Context, args NewFilamentSpoolsArgs) (*[]*FilamentSpoolResolver, error) {
	if args.ID != nil {
		filamentBrand, err := NewFilamentSpool(ctx, NewFilamentSpoolArgs{ID: *args.ID})
		resolvers := []*FilamentSpoolResolver{filamentBrand}
		return &resolvers, err
	}

	var resolvers []*FilamentSpoolResolver
	var errs errors.Errors
	filamentSpools, err := loader.LoadFilamentSpools(ctx, "")

	if err != nil {
		return nil, err
	}

	for _, filamentSpoolID := range filamentSpools.IDs {
		filamentSpool, err := NewFilamentSpool(ctx, NewFilamentSpoolArgs{ID: filamentSpoolID})
		if err != nil {
			errs = append(errs, err)
			continue
		}
		resolvers = append(resolvers, filamentSpool)
	}

	return &resolvers, errs.Err()
}

// ID returns the ID of the filament spool.
func (r *FilamentSpoolResolver) ID() graphql.ID {
	return graphql.ID(r.filamentSpool.ID)
}

// Brand returns a resolver for the brand the filament is.
func (r *FilamentSpoolResolver) Brand(ctx context.Context) (*FilamentBrandResolver, error) {
	return NewFilamentBrand(ctx, NewFilamentBrandArgs{ID: r.filamentSpool.FilamentBrand})
}
