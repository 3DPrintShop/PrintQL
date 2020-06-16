package resolver_test

import (
	"context"
	"fmt"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/3DPrintShop/PrintQL/resolver"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/magiconair/properties/assert"
	"testing"
)

const (
	testID      = "testID"
	testName    = "Test Name"
	testAltText = "alt text what it is"
	testType    = "Jiff"
)

var (
	testMetadata = printdb.ProjectMetadata{
		Public: false,
	}
	componentIDs   = []string{"component1", "component2", "component3"}
	testComponents = printdb.ProjectComponentPage{
		ComponentIds: componentIDs,
	}

	imageIDs   = []string{"media1", "media2", "media3"}
	testImages = printdb.ProjectMediaPage{
		MediaIds: imageIDs,
	}
)

type mockProjectGetter struct {
	loaded bool
}

func (mock mockProjectGetter) Component(componentId string) (printdb.Component, error) {
	return printdb.Component{
		ID:   componentId,
		Name: testName,
		Type: testType,
	}, nil
}

func (mock mockProjectGetter) Image(mediaId string) (printdb.Image, error) {
	return printdb.Image{
		ID:      mediaId,
		AltText: testAltText,
		Type:    testType,
	}, nil
}

func (mock mockProjectGetter) GetProject(id string) (printdb.Project, error) {
	if !mock.loaded {
		mock.loaded = true
		return printdb.Project{
			ID:         id,
			Name:       testName,
			Metadata:   testMetadata,
			Components: testComponents,
			Images:     testImages,
		}, nil
	}
	return printdb.Project{ID: "Loaded twice"}, fmt.Errorf("Function was called a second time")
}

func (mock mockProjectGetter) GetProjects(pageId *string) (printdb.ProjectPage, error) {
	ids := []string{"test", "test2", "test3"}
	page := printdb.ProjectPage{
		ProjectIDs: ids,
		NextKey:    pageId,
	}
	return page, nil
}

func TestProjectResolver_TestFilamentBrand(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, loader.ProjectLoaderKey, dataloader.NewBatchedLoader(loader.NewProjectLoader(mockProjectGetter{loaded: false})))
	ctx = context.WithValue(ctx, loader.ProjectPageLoaderKey, dataloader.NewBatchedLoader(loader.NewProjectPageLoader(mockProjectGetter{loaded: false})))
	ctx = context.WithValue(ctx, loader.MediaLoaderKey, dataloader.NewBatchedLoader(loader.NewMediaLoader(mockProjectGetter{loaded: false})))
	ctx = context.WithValue(ctx, loader.ComponentLoaderKey, dataloader.NewBatchedLoader(loader.NewComponentLoader(mockProjectGetter{loaded: false})))

	t.Run("Project Resolver", func(t *testing.T) {
		project, err := resolver.NewProject(ctx, resolver.NewProjectArgs{ID: testID})

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, graphql.ID(testID), project.ID())
		assert.Equal(t, testName, project.Name())
	})

	t.Run("Projects resolver", func(t *testing.T) {
		projects, err := resolver.NewProjects(ctx, resolver.NewProjectsArgs{ID: nil})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 3, len(*projects), fmt.Sprintf("Wanted 3 projects but got %d", len(*projects)))

		for _, v := range *projects {
			assert.Equal(t, testName, v.Name())
			assert.Equal(t, false, v.Public())
			assert.Equal(t, 3, len(*v.Images(ctx)), "images for the test project.")

			for _, m := range *v.Images(ctx) {
				assert.Equal(t, testType, m.Type())
				assert.Equal(t, testAltText, m.AltText())
				assert.Matches(t, string(m.ID()), "media.*")
				assert.Matches(t, m.Path(), "/media/media.Jiff")
			}

			for _, c := range *v.Components(ctx) {
				assert.Equal(t, testName, c.Name())
				assert.Equal(t, testType, c.Type())
				assert.Matches(t, string(c.ID()), "component.*")
			}
		}
	})
}
