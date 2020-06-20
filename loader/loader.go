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
	ProjectLoaderKey        key = "project"
	ProjectPageLoaderKey    key = "projectPage"
	ComponentLoaderKey      key = "component"
	componentPageLoaderKey  key = "componentPage"
	MediaLoaderKey          key = "media"
	FilamentBrandLoaderKey  key = "filamentBrand"
	FilamentBrandsLoaderKey key = "filamentBrands"
	FilamentSpoolLoaderKey  key = "filamentSpool"
	filamentSpoolsLoaderKey key = "filamentSpools"
)

// Client is an interface composed of all of the interfaces required by the loaders.
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
	filamentSpoolGetter
	filamentSpoolsGetter
}

// Initialize creates a loader lookup collection
func Initialize(boltClient Client) Collection {
	return Collection{
		lookup: map[key]dataloader.BatchFunc{
			printerLoaderKey:        newPrinterLoader(boltClient),
			printerPageLoaderKey:    newPrinterPageLoader(boltClient),
			ProjectLoaderKey:        NewProjectLoader(boltClient),
			ProjectPageLoaderKey:    NewProjectPageLoader(boltClient),
			ComponentLoaderKey:      NewComponentLoader(boltClient),
			componentPageLoaderKey:  newComponentPageLoader(boltClient),
			MediaLoaderKey:          NewMediaLoader(boltClient),
			FilamentBrandLoaderKey:  NewFilamentBrandLoader(boltClient),
			FilamentBrandsLoaderKey: NewFilamentBrandsLoader(boltClient),
			FilamentSpoolLoaderKey:  NewFilamentSpoolLoader(boltClient),
			filamentSpoolsLoaderKey: newFilamentSpoolsLoader(boltClient),
		},
	}
}

// Collection represents a collection of loaders
type Collection struct {
	lookup map[key]dataloader.BatchFunc
}

// Attach attaches the loaders in the collection to the context.
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
