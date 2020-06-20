package loader

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/dataloader"
	"sync"
)

// LoadMedia loads media using the loader associated with the context.
func LoadMedia(ctx context.Context, mediaID string) (printdb.Image, error) {
	var media printdb.Image

	ldr, err := extract(ctx, MediaLoaderKey)
	if err != nil {
		return media, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(mediaID))()
	if err != nil {
		return media, err
	}

	media, ok := data.(printdb.Image)
	if !ok {
		return media, errors.WrongType(media, data)
	}

	return media, nil
}

type mediaGetter interface {
	Image(mediaID string) (printdb.Image, error)
}

type mediaLoader struct {
	get mediaGetter
}

// NewMediaLoader creates a dataloader.BatchFunc for loading media.
func NewMediaLoader(client mediaGetter) dataloader.BatchFunc {
	return mediaLoader{get: client}.loadBatch
}

func (ldr mediaLoader) loadBatch(ctx context.Context, urls dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url dataloader.Key) {
			defer wg.Done()

			data, err := ldr.get.Image(url.String())
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, url)
	}

	wg.Wait()

	return results
}
