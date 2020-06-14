package printdb_test

import (
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_TestFilamentBrand(t *testing.T) {
	type test struct {
		name           string
		brandsToCreate int
		startingWeight int32
		purchaseLink   string
	}

	tests := []test{
		{name: "one filament", brandsToCreate: 1, startingWeight: 10, purchaseLink: "amazon"},
		{name: "two filaments", brandsToCreate: 3, startingWeight: 1000, purchaseLink: "prusa"},
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
			t.Run("Create Filament Brand", func(t *testing.T) {
				for i := 0; i < test.brandsToCreate; i++ {
					filamentID, err := client.CreateFilamentBrand(TestName)

					if err != nil {
						t.Error(err)
					}

					assert.NotEqual(t, "", filamentID)

					err = client.SetFilamentStartWeight(filamentID, test.startingWeight)
					if err != nil {
						t.Error(err)
					}

					err = client.SetFilamentPurchaseLink(filamentID, test.purchaseLink)
					if err != nil {
						t.Error(err)
					}
				}
			})

			t.Run("Get Filament Brands", func(t *testing.T) {
				filamentBrandPage, err := client.GetFilamentBrands(nil)

				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.brandsToCreate, len(filamentBrandPage.IDs))
			})

			t.Run("Get Each Filament Brand", func(t *testing.T) {
				filamentBrandPage, err := client.GetFilamentBrands(nil)

				if err != nil {
					t.Error(err)
				}

				for _, filamentID := range filamentBrandPage.IDs {
					filament, err := client.GetFilamentBrand(filamentID)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, filamentID, filament.ID)
					assert.Equal(t, TestName, filament.Name)
					assert.Equal(t, test.startingWeight, filament.StartWeight)
					assert.Equal(t, test.purchaseLink, filament.PurchaseLink)
				}
			})
		})

		teardown(context)

	}
}

func TestClient_TestFilamentSpool(t *testing.T) {
	type test struct {
		name           string
		brandsToCreate int
		spoolsPerBrand int
		startingWeight int32
	}

	tests := []test{
		{name: "one brand 10 spools", brandsToCreate: 1, spoolsPerBrand: 10, startingWeight: 10},
		{name: "three brands one thousand spools", brandsToCreate: 3, spoolsPerBrand: 50, startingWeight: 1000},
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
			t.Run("Create Brands and Spools", func(t *testing.T) {
				for i := 0; i < test.brandsToCreate; i++ {
					filamentID, err := client.CreateFilamentBrand(TestName)

					if err != nil {
						t.Error(err)
					}

					assert.NotEqual(t, "", filamentID)

					err = client.SetFilamentStartWeight(filamentID, test.startingWeight)
					if err != nil {
						t.Error(err)
					}

					for x := 0; x < test.spoolsPerBrand; x++ {
						spoolID, err := client.CreateFilamentSpool(filamentID)

						if err != nil {
							t.Error(err)
						}

						assert.NotEqual(t, "", spoolID)
					}
				}
			})

			t.Run("Get Spools", func(t *testing.T) {
				filamentSpoolPage, err := client.GetFilamentSpools(nil)

				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.spoolsPerBrand*test.brandsToCreate, len(filamentSpoolPage.IDs))
			})

			t.Run("Get Each Spool", func(t *testing.T) {
				filamentSpoolPage, err := client.GetFilamentSpools(nil)

				if err != nil {
					t.Error(err)
				}

				for _, filamentID := range filamentSpoolPage.IDs {
					filament, err := client.GetFilamentSpool(filamentID)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, filamentID, filament.ID)
					assert.Equal(t, test.startingWeight, filament.RemainingWeight)
					assert.NotNil(t, filament.FilamentBrand)
					assert.NotEqual(t, "", filament.FilamentBrand)
				}
			})
		})

		teardown(context)

	}
}
