package loader

import (
	"context"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"sync"
)

func LoadMedia(ctx context.Context, componentId string) (printdb.Image, error) {
	fmt.Println("Component ID to load: " + componentId)
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

	fmt.Printf("Loader: Media: %v", media)

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
			fmt.Printf("%v\n", data)
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, url)
	}

	wg.Wait()

	return results
}
