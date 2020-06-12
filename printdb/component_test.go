package printdb_test

import (
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_TestComponentCreationAndRetrieval(t *testing.T) {
	type test struct {
		name               string
		componentsToCreate int
	}

	tests := []test{
		{name: "one component", componentsToCreate: 1},
		{name: "three components", componentsToCreate: 3},
	}

	for _, test := range tests {
		context, err := setup()
		if err != nil {
			t.Error(err)
		}

		client, err := printdb.NewClient(context.db)

		t.Run(test.name, func(t *testing.T) {
			t.Run("Create Component", func(t *testing.T) {
				for i := 0; i < test.componentsToCreate; i++ {
					printerID, err := client.CreateComponent(printdb.NewComponentRequest{
						Name: TestComponentName,
						Type: TestComponentType,
					})

					if err != nil {
						t.Error(err)
					}

					assert.NotEqual(t, "", printerID)
				}
			})

			t.Run("Get Components", func(t *testing.T) {
				componentPage, err := client.Components(nil)

				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.componentsToCreate, len(componentPage.ComponentIds))
			})

			t.Run("Get Each Component", func(t *testing.T) {
				componentPage, err := client.Components(nil)

				if err != nil {
					t.Error(err)
				}

				for _, componentID := range componentPage.ComponentIds {
					component, err := client.Component(componentID)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, componentID, component.ID)
					assert.Equal(t, TestComponentType, component.Type)
					assert.Equal(t, TestComponentName, component.Name)
					assert.Equal(t, 0, len(component.Projects.ProjectIds))
				}
			})
		})

		if err != nil {
			t.Error(err)
			return
		}

		teardown(context)
	}
}
