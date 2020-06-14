package loader

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/dataloader"
	"sync"
)

// LoadFilamentBrand retrieves the filament brand loader from the context, and uses it to load the specified filament brand
func LoadFilamentSpool(ctx context.Context, filamentBrandID string) (printdb.FilamentSpool, error) {
	var filamentSpool printdb.FilamentSpool

	ldr, err := extract(ctx, filamentSpoolLoaderKey)
	if err != nil {
		return filamentSpool, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(filamentBrandID))()
	if err != nil {
		return filamentSpool, err
	}

	filamentSpool, ok := data.(printdb.FilamentSpool)
	if !ok {
		return filamentSpool, errors.WrongType(filamentSpool, data)
	}

	return filamentSpool, nil
}

func LoadFilamentSpools(ctx context.Context, filamentBrandID string) (printdb.IdentifierPage, error) {
	var filamentSpools printdb.IdentifierPage

	ldr, err := extract(ctx, filamentSpoolsLoaderKey)
	if err != nil {
		return filamentSpools, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(filamentBrandID))()
	if err != nil {
		return filamentSpools, err
	}

	filamentBrands, ok := data.(printdb.IdentifierPage)
	if !ok {
		return filamentBrands, errors.WrongType(filamentBrands, data)
	}

	return filamentBrands, nil
}

type filamentSpoolGetter interface {
	GetFilamentSpool(id string) (printdb.FilamentSpool, error)
}

type filamentSpoolsGetter interface {
	GetFilamentSpools(id *string) (printdb.IdentifierPage, error)
}

type filamentSpoolLoader struct {
	get filamentSpoolGetter
}

type filamentSpoolsLoader struct {
	get filamentSpoolsGetter
}

func newFilamentSpoolLoader(client filamentSpoolGetter) dataloader.BatchFunc {
	return filamentSpoolLoader{get: client}.loadBatch
}

func newFilamentSpoolsLoader(client filamentSpoolsGetter) dataloader.BatchFunc {
	return filamentSpoolsLoader{get: client}.loadBatch
}

func (ldr filamentSpoolLoader) loadBatch(ctx context.Context, urls dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url dataloader.Key) {
			defer wg.Done()

			data, err := ldr.get.GetFilamentSpool(url.String())
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, url)
	}

	wg.Wait()

	return results
}

func (ldr filamentSpoolsLoader) loadBatch(ctx context.Context, ids dataloader.Keys) []*dataloader.Result {
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
			data, err := ldr.get.GetFilamentSpools(&idString)
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, id)
	}

	wg.Wait()

	return results
}
