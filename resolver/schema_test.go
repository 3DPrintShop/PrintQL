package resolver

import (
	"github.com/3DPrintShop/PrintQL/schema"
	"github.com/graph-gophers/graphql-go"
	"testing"
)

func TestResolversSatisfySchema(t *testing.T) {
	rootResolver := &SchemaResolver{}
	_, err := graphql.ParseSchema(schema.String(), rootResolver)
	if err != nil {
		t.Error(err)
	}
}
