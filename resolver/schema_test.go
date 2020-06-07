package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/vitiock/PrintQL/schema"
	"github.com/vitiock/oauthgraphql/resolver"
	"testing"
)

func TestResolversSatisfySchema(t *testing.T) {
	rootResolver := &resolver.QueryResolver{}
	_, err := graphql.ParseSchema(schema.String(), rootResolver)
	if err != nil {
		t.Error(err)
	}
}
