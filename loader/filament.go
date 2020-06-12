package loader

import (
	"context"
	"fmt"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/dataloader"
	"sync"
)

// LoadFilamentBrand retrieves the filament brand loader from the context, and uses it to load the specified filament brand
func LoadFilamentBrand(ctx context.Context, filamentBrandID string) (printdb.FilamentBrand, error) {
	var filamentBrand printdb.FilamentBrand

	ldr, err := extract(ctx, filamentBrandLoaderKey)
	if err != nil {
		return filamentBrand, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(filamentBrandID))()
	if err != nil {
		return filamentBrand, err
	}

	filamentBrand, ok := data.(printdb.FilamentBrand)
	if !ok {
		return filamentBrand, errors.WrongType(filamentBrand, data)
	}

	return filamentBrand, nil
}

func LoadFilamentBrands(ctx context.Context, filamentBrandID string) (printdb.FilamentBrandPage, error) {
	var filamentBrands printdb.FilamentBrandPage

	ldr, err := extract(ctx, filamentBrandsLoaderKey)
	if err != nil {
		return filamentBrands, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(filamentBrandID))()
	if err != nil {
		return filamentBrands, err
	}

	filamentBrands, ok := data.(printdb.FilamentBrandPage)
	if !ok {
		return filamentBrands, errors.WrongType(filamentBrands, data)
	}

	return filamentBrands, nil
}

type filamentBrandGetter interface {
	GetFilamentBrand(id string) (printdb.FilamentBrand, error)
}

type filamentBrandsGetter interface {
	GetFilamentBrands(id *string) (printdb.FilamentBrandPage, error)
}

type filamentBrandLoader struct {
	get filamentBrandGetter
}

type filamentBrandsLoader struct {
	get filamentBrandsGetter
}

func newFilamentBrandLoader(client filamentBrandGetter) dataloader.BatchFunc {
	return filamentBrandLoader{get: client}.loadBatch
}

func newFilamentBrandsLoader(client filamentBrandsGetter) dataloader.BatchFunc {
	return filamentBrandsLoader{get: client}.loadBatch
}

func (ldr filamentBrandLoader) loadBatch(ctx context.Context, urls dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url dataloader.Key) {
			defer wg.Done()

			data, err := ldr.get.GetFilamentBrand(url.String())
			fmt.Printf("%v\n", data)
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, url)
	}

	wg.Wait()

	return results
}

func (ldr filamentBrandsLoader) loadBatch(ctx context.Context, ids dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(ids)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, id := range ids {
		go func(i int, id dataloader.Key) {
			defer wg.Done()
			idString := id.String()
			data, err := ldr.get.GetFilamentBrands(&idString)
			fmt.Printf("%v\n", data)
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, id)
	}

	wg.Wait()

	return results
}
