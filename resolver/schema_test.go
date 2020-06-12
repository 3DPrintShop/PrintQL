package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/schema"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"net/http"
	"testing"
)

func setup(loaderId key, loader dataloader.BatchFunc) (context.Context, error) {
	request, err := http.NewRequest("POST", "/graphql", nil)
	if err != nil {
		return nil, err
	}

	ctx := request.Context()
	ctx = context.WithValue(ctx, loaderId, dataloader.NewBatchedLoader(loader))

	return ctx, nil
}

func TestResolversSatisfySchema(t *testing.T) {
	rootResolver := &SchemaResolver{}
	_, err := graphql.ParseSchema(schema.String(), rootResolver)
	if err != nil {
		t.Error(err)
	}
}
