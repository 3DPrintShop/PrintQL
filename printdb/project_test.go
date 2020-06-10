package printdb_test

import (
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_TestProjectCreationAndRetrieval(t *testing.T) {
	type test struct {
		name             string
		projectsToCreate int
		componentsToAdd  int
		imagesToAdd      int
	}

	tests := []test{
		{name: "single project no components or images", projectsToCreate: 1, componentsToAdd: 0, imagesToAdd: 0},
		{name: "3 projects no components or images", projectsToCreate: 3, componentsToAdd: 0, imagesToAdd: 0},
		{name: "Single project 2 images", projectsToCreate: 1, componentsToAdd: 0, imagesToAdd: 2},
		{name: "Single project 2 components", projectsToCreate: 1, componentsToAdd: 2, imagesToAdd: 0},
		{name: "Single Project 2 components 2 images", projectsToCreate: 1, componentsToAdd: 2, imagesToAdd: 2},
		{name: "3 Projects 2 components 2 images", projectsToCreate: 3, componentsToAdd: 2, imagesToAdd: 2},
	}

	for _, test := range tests {
		context, err := setup()
		if err != nil {
			t.Error(err)
		}

		client, err := printdb.NewClient(context.db)

		if err != nil {
			t.Error(err)
			return
		}

		t.Run(test.name, func(t *testing.T) {

			t.Run("Create Project", func(t *testing.T) {
				for i := 0; i < test.projectsToCreate; i++ {
					projectID, err := client.CreateProject(printdb.NewProjectRequest{
						Name: TestName,
					})

					if err != nil {
						t.Error(err)
						return
					}

					if projectID == "" {
						t.Error("ProjectID was returned as an empty string.")
					}
					t.Run("Add images to project", func(t *testing.T) {
						for x := 0; x < test.imagesToAdd; x++ {
							altText := string(x)

							imageId, err := client.CreateImage(printdb.NewImageRequest{
								Type:    ".png",
								AltText: &altText,
							})

							if err != nil {
								t.Error(err)
							}

							err = client.AssociateImageWithProject(printdb.AssociateImageWithProjectRequest{ProjectId: projectID, ImageId: imageId, Type: ".png"})
							if err != nil {
								t.Error(err)
							}
						}
					})

					t.Run("Add components to project", func(t *testing.T) {
						for x := 0; x < test.componentsToAdd; x++ {
							componentId, err := client.CreateComponent(printdb.NewComponentRequest{
								Type: "STL",
								Name: string(x),
							})

							if err != nil {
								t.Error(err)
								continue
							}

							err = client.AssociateComponentWithProject(printdb.AssociateComponentWithProjectRequest{
								ProjectId:   projectID,
								ComponentId: componentId,
							})

							if err != nil {
								t.Error(err)
							}
						}
					})
				}
			})

			t.Run("Get Projects", func(t *testing.T) {
				projectPage, err := client.Projects(nil)

				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, len(projectPage.ProjectIds), test.projectsToCreate)
				assert.Nil(t, projectPage.NextKey)

			})

			t.Run("Get Each Project", func(t *testing.T) {
				projectPage, err := client.Projects(nil)

				if err != nil {
					t.Error(err)
				}

				for _, projectID := range projectPage.ProjectIds {
					project, err := client.Project(projectID)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, project.ID, projectID)
					assert.Equal(t, project.Name, TestName)
					assert.Equal(t, project.Metadata.Public, false)
					assert.Equal(t, len(project.Images.MediaIds), test.imagesToAdd)
					assert.Equal(t, len(project.Components.ComponentIds), test.componentsToAdd)
				}
			})
		})
		teardown(context)
	}
}
