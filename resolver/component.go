package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	graphql "github.com/graph-gophers/graphql-go"
)

// ComponentResolver resolves the component type.
type ComponentResolver struct {
	Component printdb.Component
}

// NewComponentsArgs are the arguments that can be passed in when querying for a component.
type NewComponentsArgs struct {
	ID *string
}

// NewComponentArgs are the arguments that can be passed in when querying for a paginated list of components.
type NewComponentArgs struct {
	ID string
}

// NewComponent returns a new resolver for a component specified by the args passed in.
func NewComponent(ctx context.Context, args NewComponentArgs) (*ComponentResolver, error) {
	component, errs := loader.LoadComponent(ctx, args.ID)
	return &ComponentResolver{Component: component}, errs
}

// NewComponents returns a set of resolvers filtered by the args passed in.
func NewComponents(ctx context.Context, args NewComponentsArgs) (*[]*ComponentResolver, error) {
	if args.ID != nil {
		component, err := NewComponent(ctx, NewComponentArgs{ID: *args.ID})
		resolvers := []*ComponentResolver{component}
		return &resolvers, err
	}

	var resolvers []*ComponentResolver
	components, err := loader.LoadComponents(ctx, "")

	if err != nil {
		return nil, err
	}

	for _, component := range components {
		resolvers = append(resolvers, &ComponentResolver{Component: component})
	}

	return &resolvers, nil
}

// ID resolves the component's ID.
func (r *ComponentResolver) ID() graphql.ID {
	return graphql.ID(r.Component.ID)
}

// Name resolves the component's name.
func (r *ComponentResolver) Name() string {
	return r.Component.Name
}

// Type resolves the component's type.
func (r *ComponentResolver) Type() string {
	return r.Component.Type
}

// Projects resolves a set of project resolvers that a component is part of.
func (r *ComponentResolver) Projects(ctx context.Context) *[]*ProjectResolver {
	var resolvers []*ProjectResolver
	for _, projectID := range r.Component.Projects.ProjectIds {
		project, err := loader.LoadProject(ctx, projectID)
		if err != nil {
			continue
		}
		resolvers = append(resolvers, &ProjectResolver{Project: project})
	}

	return &resolvers
}
