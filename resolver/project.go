package resolver

import (
	"context"
	"fmt"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
)

type ProjectResolver struct {
	Project printdb.Project
}

type NewProjectsArgs struct {
	ID *string
}

type NewProjectArgs struct {
	ID string
}

type Project struct {
	Id   graphql.ID
	Name string
}

func NewProject(ctx context.Context, args NewProjectArgs) (*ProjectResolver, error) {
	fmt.Printf("Project request: %s\n", args.ID)
	project, errs := loader.LoadProject(ctx, args.ID)

	return &ProjectResolver{Project: project}, errs
}

func NewProjects(ctx context.Context, args NewProjectsArgs) (*[]*ProjectResolver, error) {
	if args.ID != nil {
		project, err := NewProject(ctx, NewProjectArgs{ID: *args.ID})
		resolvers := []*ProjectResolver{project}
		return &resolvers, err
	}

	var resolvers []*ProjectResolver
	projects, err := loader.LoadProjects(ctx, "")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for _, project := range projects {
		resolvers = append(resolvers, &ProjectResolver{Project: project})
	}

	return &resolvers, nil
}

func (r *ProjectResolver) ID() graphql.ID {
	return graphql.ID(r.Project.ID)
}

func (r *ProjectResolver) Name() string {
	return r.Project.Name
}

func (r *ProjectResolver) Public() bool {
	return r.Project.Metadata.Public
}

func (r *ProjectResolver) Components(ctx context.Context) *[]*ComponentResolver {
	var resolvers []*ComponentResolver
	fmt.Printf("Attempting to get components for project %s\n", r.Project.ID)
	for _, componetId := range r.Project.Components.ComponentIds {
		fmt.Printf("Getting componenet %s\n", componetId)
		component, err := loader.LoadComponent(ctx, componetId)
		if err != nil {
			fmt.Printf("Failed to load component by id %s\n", componetId)
			fmt.Println(err)
		} else {
			resolvers = append(resolvers, &ComponentResolver{Component: component})
		}
	}

	return &resolvers
}

func (r *ProjectResolver) Images(ctx context.Context) *[]*MediaResolver {
	var resolvers []*MediaResolver
	fmt.Printf("Attempting to get components for project %s\n", r.Project.ID)
	for _, mediaID := range r.Project.Images.MediaIds {
		mediaResolver, err := NewMedia(ctx, NewMediaArgs{ID: mediaID})
		if err == nil {
			resolvers = append(resolvers, mediaResolver)
		} else {
			fmt.Printf("Failed to get media resolver for id: %s\n", mediaID)
		}
	}

	return &resolvers
}
