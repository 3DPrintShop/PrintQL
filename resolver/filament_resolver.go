package resolver

import (
	"context"
)

type FilamentResolver struct {
}

func NewFilamentResolver() (*FilamentResolver, error) {
	resolver := FilamentResolver{}
	return &resolver, nil
}

func (r FilamentResolver) Brands(ctx context.Context, args NewFilamentBrandsArgs) (*[]*FilamentBrandResolver, error) {
	return NewFilamentBrands(ctx, args)
}

func (r FilamentResolver) Spools(ctx context.Context, args NewFilamentSpoolsArgs) (*[]*FilamentSpoolResolver, error) {
	return NewFilamentSpools(ctx, args)
}
