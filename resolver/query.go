package resolver

import (
	"context"
	"github.com/graph-gophers/graphql-go"
)

// PrintersQueryArgs are a set of arguments that can be passed in from graphql to dictate how to filter the list of printers.
type PrintersQueryArgs struct {
	ID *graphql.ID
}

// Printers go brrr.
func (r SchemaResolver) Printers(ctx context.Context, args PrintersQueryArgs) (*[]*PrinterResolver, error) {
	if args.ID != nil {
		id := string(*args.ID)
		return NewPrinters(ctx, NewPrintersArgs{ID: &id})
	}
	return NewPrinters(ctx, NewPrintersArgs{ID: nil})
}

// ProjectsQueryArgs are arguments that can be passed via graphQL for how to create the projects resolver.
type ProjectsQueryArgs struct {
	ID *graphql.ID
}

// Projects creates a resolver that resolves a list of ProjectResolvers.
func (r SchemaResolver) Projects(ctx context.Context, args ProjectsQueryArgs) (*[]*ProjectResolver, error) {
	if args.ID != nil {
		id := string(*args.ID)
		return NewProjects(ctx, NewProjectsArgs{ID: &id})
	}
	return NewProjects(ctx, NewProjectsArgs{ID: nil})
}

// ComponentsQueryArgs are the arguments that can be passed in via graphql for how to create the ComponentsResolver
type ComponentsQueryArgs struct {
	ID *graphql.ID
}

// Components creates a resolver that resolves a list of ComponentResolvers
func (r SchemaResolver) Components(ctx context.Context, args ComponentsQueryArgs) (*[]*ComponentResolver, error) {
	if args.ID != nil {
		id := string(*args.ID)
		return NewComponents(ctx, NewComponentsArgs{ID: &id})
	}
	return NewComponents(ctx, NewComponentsArgs{ID: nil})
}

// FilamentBrandQueryArgs are the parameters that are passed in via graphql to control how to filter the list of FilamentBrandResolvers.
type FilamentBrandQueryArgs struct {
	ID *graphql.ID
}

// FilamentBrands resolves a list of FilamentBrandResolves
func (r SchemaResolver) FilamentBrands(ctx context.Context, args FilamentBrandQueryArgs) (*[]*FilamentBrandResolver, error) {
	if args.ID != nil {
		id := string(*args.ID)
		return NewFilamentBrands(ctx, NewFilamentBrandsArgs{ID: &id})
	}
	return NewFilamentBrands(ctx, NewFilamentBrandsArgs{ID: nil})
}

// Filament returns a resolver for filaments.
func (r SchemaResolver) Filament() (*FilamentResolver, error) {
	return NewFilamentResolver()
}

// Self returns a resolver that details the user that is calling graphql.
func (r SchemaResolver) Self(ctx context.Context) (*AccountResolver, error) {
	return NewAccount(ctx)
}
