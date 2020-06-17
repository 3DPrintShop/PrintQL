package resolver

import (
	"context"
)

// FilamentResolver resolves to a set of queries for querying filament objects
type FilamentResolver struct {
}

// NewFilamentResolver creates a new resolver for filament queries.
func NewFilamentResolver() (*FilamentResolver, error) {
	resolver := FilamentResolver{}
	return &resolver, nil
}

// Brands resolves a set of brands filter by NewFilamentBrandsArgs.
func (r FilamentResolver) Brands(ctx context.Context, args NewFilamentBrandsArgs) (*[]*FilamentBrandResolver, error) {
	return NewFilamentBrands(ctx, args)
}

// Spools resolves a set of spools filter by NewFilamentSpoolsArgs.
func (r FilamentResolver) Spools(ctx context.Context, args NewFilamentSpoolsArgs) (*[]*FilamentSpoolResolver, error) {
	return NewFilamentSpools(ctx, args)
}
