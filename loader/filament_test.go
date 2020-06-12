package loader

import (
	"fmt"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testBrandID      = "test brand"
	testName         = "testName"
	testPurchaseLink = "purchase link"
	testStartWeight  = int32(10)
	testSpoolWeight  = int32(5)
)

type mockFilamentBrandGetter struct {
	loaded     bool
	loadedPage bool
}

func (mock mockFilamentBrandGetter) GetFilamentBrand(id string) (printdb.FilamentBrand, error) {
	if !mock.loaded {
		mock.loaded = true
		return printdb.FilamentBrand{
			ID:           id,
			Name:         testName,
			PurchaseLink: testPurchaseLink,
			StartWeight:  testStartWeight,
			SpoolWeight:  testSpoolWeight,
		}, nil
	}
	return printdb.FilamentBrand{ID: "Loaded twice"}, fmt.Errorf("Function was called a second time")
}

func (mock mockFilamentBrandGetter) GetFilamentBrands(id *string) (printdb.FilamentBrandPage, error) {
	if !mock.loadedPage {
		mock.loadedPage = true
		return printdb.FilamentBrandPage{
			FilamentBrandIDs: []string{"test", "test2", "test3"},
			NextPage:         nil,
		}, nil
	}
	return printdb.FilamentBrandPage{FilamentBrandIDs: nil, NextPage: nil}, fmt.Errorf("Function was called a second time")
}

func TestLoader_TestFilamentBrand(t *testing.T) {
	context, err := setup(filamentBrandLoaderKey, newFilamentBrandLoader(mockFilamentBrandGetter{loaded: false}))
	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Loads filament", func(t *testing.T) {
		filamentBrand, err := LoadFilamentBrand(context, testBrandID)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, testBrandID, filamentBrand.ID)
		assert.Equal(t, testName, filamentBrand.Name)
		assert.Equal(t, testPurchaseLink, filamentBrand.PurchaseLink)
		assert.Equal(t, testSpoolWeight, filamentBrand.SpoolWeight)
		assert.Equal(t, testStartWeight, filamentBrand.StartWeight)
	})

	t.Run("Loads filament only once", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			filamentBrand, err := LoadFilamentBrand(context, testBrandID)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, testBrandID, filamentBrand.ID)
			assert.Equal(t, testName, filamentBrand.Name)
			assert.Equal(t, testPurchaseLink, filamentBrand.PurchaseLink)
			assert.Equal(t, testSpoolWeight, filamentBrand.SpoolWeight)
			assert.Equal(t, testStartWeight, filamentBrand.StartWeight)
		}
	})
}

func TestLoader_TestFilamentBrands(t *testing.T) {
	context, err := setup(filamentBrandsLoaderKey, newFilamentBrandsLoader(mockFilamentBrandGetter{loadedPage: false}))
	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Loads filament page", func(t *testing.T) {
		filamentBrands, err := LoadFilamentBrands(context, testBrandID)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 3, len(filamentBrands.FilamentBrandIDs))
		if len(filamentBrands.FilamentBrandIDs) != 3 {
			return
		}
		assert.Equal(t, "test", filamentBrands.FilamentBrandIDs[0])
		assert.Equal(t, "test2", filamentBrands.FilamentBrandIDs[1])
		assert.Equal(t, "test3", filamentBrands.FilamentBrandIDs[2])
		assert.Nil(t, filamentBrands.NextPage)
	})

	t.Run("Loads filament page only once", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			filamentBrands, err := LoadFilamentBrands(context, testBrandID)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, "test", filamentBrands.FilamentBrandIDs[0])
			assert.Equal(t, "test2", filamentBrands.FilamentBrandIDs[1])
			assert.Equal(t, "test3", filamentBrands.FilamentBrandIDs[2])
			assert.Nil(t, filamentBrands.NextPage)
		}
	})
}
