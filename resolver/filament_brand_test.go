package resolver_test

import (
	"fmt"
	"github.com/3DPrintShop/PrintQL/resolver"
	"github.com/graph-gophers/graphql-go"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFilamentBrandResolver_TestFilamentBrand(t *testing.T) {
	ctx := getContext()

	t.Run("Filament Brand Resolver", func(t *testing.T) {
		brand, err := resolver.NewFilamentBrand(ctx, resolver.NewFilamentBrandArgs{ID: testID})

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, graphql.ID(testID), brand.ID())
		assert.Equal(t, testName, brand.Name())
		assert.Equal(t, int32(testStartWeight), *brand.StartWeight())
		assert.Equal(t, int32(testSpoolWeight), *brand.SpoolWeight())
		assert.Equal(t, testPurchaseLink, *brand.PurchaseLink())
	})

	t.Run("Projects resolver", func(t *testing.T) {
		brands, err := resolver.NewFilamentBrands(ctx, resolver.NewFilamentBrandsArgs{ID: nil})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 5, len(*brands), fmt.Sprintf("Wanted 5 brands but got %d", len(*brands)))

		for _, v := range *brands {
			assert.Equal(t, testName, v.Name())
			assert.Equal(t, int32(testStartWeight), *v.StartWeight())
			assert.Equal(t, int32(testSpoolWeight), *v.SpoolWeight())
			assert.Equal(t, testPurchaseLink, *v.PurchaseLink())
		}
	})

	t.Run("Projects resolver single ID", func(t *testing.T) {
		testBrandID := testID
		brands, err := resolver.NewFilamentBrands(ctx, resolver.NewFilamentBrandsArgs{ID: &testBrandID})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 1, len(*brands), fmt.Sprintf("Wanted 5 brands but got %d", len(*brands)))

		for _, v := range *brands {
			assert.Equal(t, testName, v.Name())
			assert.Equal(t, int32(testStartWeight), *v.StartWeight())
			assert.Equal(t, int32(testSpoolWeight), *v.SpoolWeight())
			assert.Equal(t, testPurchaseLink, *v.PurchaseLink())
		}
	})
}
