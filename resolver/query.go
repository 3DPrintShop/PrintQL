package resolver

import (
	"context"
	"fmt"
	"github.com/graph-gophers/graphql-go"
)

type PrintersQueryArgs struct {
	ID *graphql.ID
}

func (r SchemaResolver) Printers(ctx context.Context, args PrintersQueryArgs) (*[]*PrinterResolver, error) {
	if args.ID != nil {
		fmt.Println("ID passed into printers")
		id := string(*args.ID)
		return NewPrinters(ctx, NewPrintersArgs{ID: &id})
	}
	return NewPrinters(ctx, NewPrintersArgs{ID: nil})
}

type ProjectsQueryArgs struct {
	ID *graphql.ID
}

func (r SchemaResolver) Projects(ctx context.Context, args ProjectsQueryArgs) (*[]*ProjectResolver, error) {
	if args.ID != nil {
		fmt.Println("ID passed into projects")
		id := string(*args.ID)
		return NewProjects(ctx, NewProjectsArgs{ID: &id})
	}
	return NewProjects(ctx, NewProjectsArgs{ID: nil})
}

type ComponentsQueryArgs struct {
	ID *graphql.ID
}

func (r SchemaResolver) Components(ctx context.Context, args ComponentsQueryArgs) (*[]*ComponentResolver, error) {
	if args.ID != nil {
		fmt.Println("ID passed into components")
		id := string(*args.ID)
		return NewComponents(ctx, NewComponentsArgs{ID: &id})
	}
	return NewComponents(ctx, NewComponentsArgs{ID: nil})
}

type FilamentBrandQueryArgs struct {
	ID *graphql.ID
}

func (r SchemaResolver) FilamentBrands(ctx context.Context, args FilamentBrandQueryArgs) (*[]*FilamentBrandResolver, error) {
	if args.ID != nil {
		fmt.Println("ID passed into components")
		id := string(*args.ID)
		return NewFilamentBrands(ctx, NewFilamentBrandsArgs{ID: &id})
	}
	return NewFilamentBrands(ctx, NewFilamentBrandsArgs{ID: nil})
}

func (r SchemaResolver) Self(ctx context.Context) (*AccountResolver, error) {
	return NewAccount(ctx)
}
