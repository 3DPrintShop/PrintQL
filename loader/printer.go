package loader

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/dataloader"
	"sync"
)

func LoadPrinter(ctx context.Context, printerId string) (printdb.Printer, error) {
	var printer printdb.Printer

	ldr, err := extract(ctx, printerLoaderKey)
	if err != nil {
		return printer, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(printerId))()
	if err != nil {
		return printer, err
	}

	printer, ok := data.(printdb.Printer)
	if !ok {
		return printer, errors.WrongType(printer, data)
	}

	return printer, nil
}

func LoadPrinters(ctx context.Context, printerPageID string) ([]printdb.Printer, error) {
	var printers []printdb.Printer

	ldr, err := extract(ctx, printerPageLoaderKey)
	if err != nil {
		return printers, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(printerPageID))()
	if err != nil {
		return printers, err
	}

	printerPage, ok := data.(printdb.PrinterPage)
	if !ok {
		return printers, errors.WrongType(printerPage, data)
	}

	for _, v := range printerPage.PrinterIds {
		printer, err := LoadPrinter(ctx, v)
		if err != nil {
			return printers, err
		}
		printers = append(printers, printer)
	}

	return printers, nil
}

type printerPageGetter interface {
	Printers(pageId *string) (printdb.PrinterPage, error)
}

type printerGetter interface {
	Printer(printerId string) (printdb.Printer, error)
}

type printerPageLoader struct {
	get printerPageGetter
}

type printerLoader struct {
	get printerGetter
}

func newPrinterLoader(client printerGetter) dataloader.BatchFunc {
	return printerLoader{get: client}.loadBatch
}

func newPrinterPageLoader(client printerPageGetter) dataloader.BatchFunc {
	return printerPageLoader{get: client}.loadBatch
}

func (ldr printerPageLoader) loadBatch(ctx context.Context, pageIDs dataloader.Keys) []*dataloader.Result {
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
			data, err := ldr.get.Printers(&pageKey)
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, pageID)
	}

	wg.Wait()

	return results
}

func (ldr printerLoader) loadBatch(ctx context.Context, urls dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url dataloader.Key) {
			defer wg.Done()

			data, err := ldr.get.Printer(url.String())
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, url)
	}

	wg.Wait()

	return results
}
