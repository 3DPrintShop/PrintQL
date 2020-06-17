package resolver_test

import (
	"fmt"
	"github.com/3DPrintShop/PrintQL/resolver"
	"github.com/graph-gophers/graphql-go"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestProjectResolver_TestGetProject(t *testing.T) {
	ctx := getContext()

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

	t.Run("Projects resolver with ID", func(t *testing.T) {
		testProjectID := testID
		projects, err := resolver.NewProjects(ctx, resolver.NewProjectsArgs{ID: &testProjectID})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 1, len(*projects), fmt.Sprintf("Wanted 3 projects but got %d", len(*projects)))

		for _, v := range *projects {
			assert.Equal(t, testID, string(v.ID()))
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
