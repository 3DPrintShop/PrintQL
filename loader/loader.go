package loader

import (
	"context"
	"fmt"
	"github.com/graph-gophers/dataloader"
)

type key string

const (
	printerLoaderKey        key = "printer"
	printerPageLoaderKey    key = "printerPage"
	projectLoaderKey        key = "project"
	projectPageLoaderKey    key = "projectPage"
	componentLoaderKey      key = "component"
	componentPageLoaderKey  key = "componentPage"
	mediaLoaderKey          key = "media"
	filamentBrandLoaderKey  key = "filamentBrand"
	filamentBrandsLoaderKey key = "filamentBrands"
	filamentSpoolLoaderKey  key = "filamentSpool"
)

type Client interface {
	printerGetter
	printerPageGetter
	projectGetter
	projectPageGetter
	componentGetter
	componentPageGetter
	mediaGetter
	filamentBrandGetter
	filamentBrandsGetter
}

func Initialize(boltClient Client) Collection {
	return Collection{
		lookup: map[key]dataloader.BatchFunc{
			printerLoaderKey:        newPrinterLoader(boltClient),
			printerPageLoaderKey:    newPrinterPageLoader(boltClient),
			projectLoaderKey:        newProjectLoader(boltClient),
			projectPageLoaderKey:    newProjectPageLoader(boltClient),
			componentLoaderKey:      newComponentLoader(boltClient),
			componentPageLoaderKey:  newComponentPageLoader(boltClient),
			mediaLoaderKey:          newMediaLoader(boltClient),
			filamentBrandLoaderKey:  newFilamentBrandLoader(boltClient),
			filamentBrandsLoaderKey: newFilamentBrandsLoader(boltClient),
		},
	}
}

type Collection struct {
	lookup map[key]dataloader.BatchFunc
}

func (c Collection) Attach(ctx context.Context) context.Context {
	for k, batchFn := range c.lookup {
		ctx = context.WithValue(ctx, k, dataloader.NewBatchedLoader(batchFn))
	}

	return ctx
}

func extract(ctx context.Context, k key) (*dataloader.Loader, error) {
	ldr, ok := ctx.Value(k).(*dataloader.Loader)
	if !ok {
		return nil, fmt.Errorf("unable to find loader: (%s) in the request context", k)
	}

	return ldr, nil
}

func (k key) String() string {
	return string(k)
}
