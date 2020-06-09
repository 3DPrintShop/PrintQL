package resolver

import (
	"context"
	"fmt"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	graphql "github.com/graph-gophers/graphql-go"
)

type ComponentResolver struct {
	Component printdb.Component
}

type NewComponentsArgs struct {
	ID *string
}

type NewComponentArgs struct {
	ID string
}

type Component struct {
	ID  graphql.ID
	URL string
}

type ComponentPage struct {
}

func NewComponent(ctx context.Context, args NewComponentArgs) (*ComponentResolver, error) {
	component, errs := loader.LoadComponent(ctx, args.ID)

	return &ComponentResolver{Component: component}, errs
}

func NewComponents(ctx context.Context, args NewComponentsArgs) (*[]*ComponentResolver, error) {
	if args.ID != nil {
		component, err := NewComponent(ctx, NewComponentArgs{ID: *args.ID})
		resolvers := []*ComponentResolver{component}
		return &resolvers, err
	}

	var resolvers []*ComponentResolver
	components, err := loader.LoadComponents(ctx, "")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for _, component := range components {
		resolvers = append(resolvers, &ComponentResolver{Component: component})
	}

	return &resolvers, nil
}

func (r *ComponentResolver) ID() graphql.ID {
	return graphql.ID(r.Component.ID)
}

func (r *ComponentResolver) Name() string {
	return r.Component.Name
}

func (r *ComponentResolver) Type() string {
	return r.Component.Type
}

func (r *ComponentResolver) Projects(ctx context.Context) *[]*ProjectResolver {
	var resolvers []*ProjectResolver
	fmt.Printf("Attempting to get projects for component %s\n", r.Component.ID)
	for _, projectID := range r.Component.Projects.ProjectIds {
		fmt.Printf("Getting component %s\n", projectID)
		project, err := loader.LoadProject(ctx, projectID)
		if err != nil {
			fmt.Printf("Failed to load component by id %s\n", projectID)
			fmt.Println(err)
		} else {
			resolvers = append(resolvers, &ProjectResolver{Project: project})
		}
	}

	return &resolvers
}
