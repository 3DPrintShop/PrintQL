package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	graphql "github.com/graph-gophers/graphql-go"
)

// ProjectResolver resolves to the project type.
type ProjectResolver struct {
	Project printdb.Project
}

// NewProjectsArgs are the arguments that can be passed in to define how to filter a list of projects.
type NewProjectsArgs struct {
	ID *string
}

// NewProjectArgs are the arguments that can be passed in to specify which project ro resolve.
type NewProjectArgs struct {
	ID string
}

// NewProject returns a project resolver for the project specified in the args.
func NewProject(ctx context.Context, args NewProjectArgs) (*ProjectResolver, error) {
	project, errs := loader.LoadProject(ctx, args.ID)

	return &ProjectResolver{Project: project}, errs
}

// NewProjects returns a set of project resolvers based on the criteria passed in.
func NewProjects(ctx context.Context, args NewProjectsArgs) (*[]*ProjectResolver, error) {
	if args.ID != nil {
		project, err := NewProject(ctx, NewProjectArgs{ID: *args.ID})
		resolvers := []*ProjectResolver{project}
		return &resolvers, err
	}

	var resolvers []*ProjectResolver
	projects, err := loader.LoadProjects(ctx, "")

	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		resolvers = append(resolvers, &ProjectResolver{Project: project})
	}

	return &resolvers, nil
}

// ID returns the identifier of the project.
func (r *ProjectResolver) ID() graphql.ID {
	return graphql.ID(r.Project.ID)
}

// Name returns the Name of the project.
func (r *ProjectResolver) Name() string {
	return r.Project.Name
}

// Public returns if the project is publicly visible.
func (r *ProjectResolver) Public() bool {
	return r.Project.Metadata.Public
}

// Components returns a list of components that are connected with the project.
func (r *ProjectResolver) Components(ctx context.Context) *[]*ComponentResolver {
	var resolvers []*ComponentResolver
	for _, componetID := range r.Project.Components.ComponentIds {
		component, err := loader.LoadComponent(ctx, componetID)
		if err == nil {
			resolvers = append(resolvers, &ComponentResolver{Component: component})
		}
	}

	return &resolvers
}

// Images returns a list of images that are connected with the project.
func (r *ProjectResolver) Images(ctx context.Context) *[]*MediaResolver {
	var resolvers []*MediaResolver
	for _, mediaID := range r.Project.Images.MediaIds {
		mediaResolver, err := NewMedia(ctx, NewMediaArgs{ID: mediaID})
		if err == nil {
			resolvers = append(resolvers, mediaResolver)
		}
	}

	return &resolvers
}
