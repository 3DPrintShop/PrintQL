package printdb_test

import (
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_TestFilamentBrand(t *testing.T) {
	type test struct {
		name             string
		printersToCreate int
		startingWeight   int
	}

	tests := []test{
		{name: "one filament", printersToCreate: 1, startingWeight: 10},
		{name: "two filaments", printersToCreate: 3, startingWeight: 1000},
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
			t.Run("Create Printer", func(t *testing.T) {
				for i := 0; i < test.printersToCreate; i++ {
					filamentID, err := client.CreateFilamentBrand(TestName)

					if err != nil {
						t.Error(err)
					}

					assert.NotEqual(t, "", filamentID)

					client.SetFilamentStartWeight(filamentID, test.startingWeight)
				}
			})

			t.Run("Get Printers", func(t *testing.T) {
				filamentBrandPage, err := client.GetFilamentBrands(nil)

				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.printersToCreate, len(filamentBrandPage.FilamentBrandIDs))
			})

			t.Run("Get Each Printer", func(t *testing.T) {
				filamentBrandPage, err := client.GetFilamentBrands(nil)

				if err != nil {
					t.Error(err)
				}

				for _, filamentID := range filamentBrandPage.FilamentBrandIDs {
					filament, err := client.GetFilamentBrand(filamentID)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, filamentID, filament.ID)
					assert.Equal(t, TestName, filament.Name)
					assert.Equal(t, test.startingWeight, filament.StartWeight)
				}
			})
		})

		teardown(context)

	}
}
