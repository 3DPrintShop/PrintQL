package loader

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/dataloader"
	"sync"
)

func LoadMedia(ctx context.Context, componentId string) (printdb.Image, error) {
	var media printdb.Image

	ldr, err := extract(ctx, mediaLoaderKey)
	if err != nil {
		return media, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(componentId))()
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
	Image(mediaId string) (printdb.Image, error)
}

type mediaLoader struct {
	get mediaGetter
}

func newMediaLoader(client mediaGetter) dataloader.BatchFunc {
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
