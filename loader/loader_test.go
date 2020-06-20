package loader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"net/http"
)

func setup(loaderID key, loader dataloader.BatchFunc) (context.Context, error) {
	request, err := http.NewRequest("POST", "/graphql", nil)
	if err != nil {
		return nil, err
	}

	ctx := request.Context()
	ctx = context.WithValue(ctx, loaderID, dataloader.NewBatchedLoader(loader))

	return ctx, nil
}
