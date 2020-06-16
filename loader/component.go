package loader

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/dataloader"
	"sync"
)

func LoadComponent(ctx context.Context, componentId string) (printdb.Component, error) {
	var component printdb.Component

	ldr, err := extract(ctx, ComponentLoaderKey)
	if err != nil {
		return component, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(componentId))()
	if err != nil {
		return component, err
	}

	component, ok := data.(printdb.Component)
	if !ok {
		return component, errors.WrongType(component, data)
	}

	return component, nil
}

func LoadComponents(ctx context.Context, componentPageID string) ([]printdb.Component, error) {
	var components []printdb.Component

	ldr, err := extract(ctx, componentPageLoaderKey)
	if err != nil {
		return components, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(componentPageID))()
	if err != nil {
		return components, err
	}

	componentPage, ok := data.(printdb.ComponentPage)
	if !ok {
		return components, errors.WrongType(componentPage, data)
	}

	for _, v := range componentPage.ComponentIds {
		component, err := LoadComponent(ctx, v)
		if err != nil {
			return components, err
		}
		components = append(components, component)
	}

	return components, nil
}

type componentPageGetter interface {
	Components(pageId *string) (printdb.ComponentPage, error)
}

type componentGetter interface {
	Component(componentId string) (printdb.Component, error)
}

type componentPageLoader struct {
	get componentPageGetter
}

type componentLoader struct {
	get componentGetter
}

func NewComponentLoader(client componentGetter) dataloader.BatchFunc {
	return componentLoader{get: client}.loadBatch
}

func newComponentPageLoader(client componentPageGetter) dataloader.BatchFunc {
	return componentPageLoader{get: client}.loadBatch
}

func (ldr componentPageLoader) loadBatch(ctx context.Context, pageIDs dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(pageIDs)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, pageID := range pageIDs {
		go func(i int, url dataloader.Key) {
			defer wg.Done()

			pageKey := pageID.String()
			data, err := ldr.get.Components(&pageKey)
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, pageID)
	}

	wg.Wait()

	return results
}

func (ldr componentLoader) loadBatch(ctx context.Context, urls dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url dataloader.Key) {
			defer wg.Done()

			data, err := ldr.get.Component(url.String())
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, url)
	}

	wg.Wait()

	return results
}
