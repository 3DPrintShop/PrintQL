package loader

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/graph-gophers/dataloader"
	"github.com/vitiock/go-octoprint"
	"sync"
)

// LoadOctoprintPrinter loads the printer state from octoprint connection
func LoadOctoprintJob(ctx context.Context, printerID string) (*octoprint.JobResponse, error) {
	var printer *octoprint.JobResponse

	ldr, err := extract(ctx, octoprintJobLoaderKey)
	if err != nil {
		return printer, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(printerID))()
	if err != nil {
		return printer, err
	}

	printer, ok := data.(*octoprint.JobResponse)
	if !ok {
		return printer, errors.WrongType(printer, data)
	}

	return printer, nil
}

type octoprintJobLoader struct {
}

func newOctoprintJobLoader() dataloader.BatchFunc {
	return octoprintJobLoader{}.loadBatch
}

func (ldr octoprintJobLoader) loadBatch(ctx context.Context, printerIDs dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(printerIDs)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, printerID := range printerIDs {
		go func(i int, printerID dataloader.Key) {
			defer wg.Done()
			printer, err := LoadPrinter(ctx, printerID.String())

			if err != nil {
				results[i] = &dataloader.Result{Data: nil, Error: err}
				return
			}

			client := octoprint.NewClient(printer.Endpoint, printer.APIKey)
			jr := octoprint.JobRequest{}

			jobResult, err := jr.Do(client)

			results[i] = &dataloader.Result{Data: jobResult, Error: err}
		}(i, printerID)
	}

	wg.Wait()

	return results
}
