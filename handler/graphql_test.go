package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/3DPrintShop/PrintQL/handler"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/3DPrintShop/PrintQL/resolver"
	"github.com/3DPrintShop/PrintQL/schema"
	"github.com/boltdb/bolt"
	"github.com/graph-gophers/graphql-go"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testContext struct {
	db *bolt.DB
}

func setup() (*testContext, error) {
	path := tempfile()
	db, err := bolt.Open(path, 0600, nil)

	if err != nil {
		return nil, err
	}

	return &testContext{
		db: db,
	}, nil
}

func teardown(ctx *testContext) {
	defer os.Remove(ctx.db.Path())
	ctx.db.Close()
}

// tempfile returns a temporary file path.
func tempfile() string {
	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}

func runQuery(h handler.GraphQL, query string) (string, error) {
	requestBodyObj := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     query,
		Variables: nil,
	}

	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(requestBodyObj); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "/graphql", &requestBody)
	if err != nil {
		return "", err
	}

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	return rr.Body.String(), nil
}

func TestGraphQL_TestGraphQL(t *testing.T) {
	root, err := resolver.NewRoot()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx, err := setup()
	defer teardown(ctx)

	if err != nil {
		t.Error(err)
		return
	}

	printClient, err := printdb.NewClient(ctx.db)

	if err != nil {
		t.Error(err)
		return
	}

	if printClient == nil {
		t.Errorf("Failed to create a print client.")
		return
	}

	h := handler.GraphQL{
		Schema:  graphql.MustParseSchema(schema.String(), root),
		Loaders: loader.Initialize(printClient),
		Client:  printClient,
		DB:      ctx.db,
	}

	type simpleQueryTest struct {
		name   string
		query  string
		result string
	}

	tests := []simpleQueryTest{
		{
			name:   "no query",
			query:  "{}",
			result: "{\"data\":{}}",
		},
		{
			name:   "projects",
			query:  "{\n  projects{\n  \tid,\n  \tname,\n    components{\n      id,\n      name,\n      type\n    }\n\t}\n}",
			result: "{\"data\":{\"projects\":[]}}",
		},
		{
			name:   "create project",
			query:  "mutation {\n  createProject(name:\"testName\"){    \n    name\n  }\n}",
			result: "{\"data\":{\"createProject\":{\"name\":\"testName\"}}}",
		},
		{
			name:   "get project with result",
			query:  "{\n  projects{\n  \tname,\n    components{\n      id,\n      name,\n      type\n    }\n\t}\n}",
			result: "{\"data\":{\"projects\":[{\"name\":\"testName\",\"components\":[]}]}}",
		},
		{
			name:   "get filament brands",
			query:  "{\n\tfilamentBrands{\n    startWeight,\n    spoolWeight,\n    name\n  }\n}",
			result: "{\"data\":{\"filamentBrands\":[]}}",
		},
		{
			name:   "get filament spools",
			query:  "{\n\tfilament{\n    spools {\n      brand {\n        name\n      }\n    }\n  }\n}",
			result: "{\"data\":{\"filament\":{\"spools\":[]}}}",
		},
		{
			name:   "get printers",
			query:  "{\n\tprinters{\n    name,\n    endpoint,    \n  }\n}",
			result: "{\"data\":{\"printers\":[]}}",
		},
		{
			name:   "get components",
			query:  "{\n  components{\n    name,\n    type,\n    projects {\n      name\n    }\n  }\n}",
			result: "{\"data\":{\"components\":[]}}",
		},
		{
			name:   "create filament brand",
			query:  "mutation{\n  filamentActions{\n    createFilamentBrand(name: \"Filament brand\"){\n      name\n    }\n  }\n}",
			result: "{\"data\":{\"filamentActions\":{\"createFilamentBrand\":{\"name\":\"Filament brand\"}}}}",
		},
		{
			name:   "get filament brands with results",
			query:  "{\n\tfilamentBrands{\n    startWeight,\n    spoolWeight,\n    name\n  }\n}",
			result: "{\"data\":{\"filamentBrands\":[{\"startWeight\":0,\"spoolWeight\":0,\"name\":\"Filament brand\"}]}}",
		},
		{
			name:   "create printer",
			query:  "mutation{\n  createPrinter(name: \"printerName\", apiKey:\"API Key\", endpoint:\"EndPoint\", integrationType:\"Octoprint\"){\n    name,\n    endpoint,    \n    integrationType\n  }\n}",
			result: "{\"data\":{\"createPrinter\":{\"name\":\"printerName\",\"endpoint\":\"EndPoint\",\"integrationType\":\"Octoprint\"}}}",
		},
		{
			name:   "get printers with results",
			query:  "{\n\tprinters{\n    name,\n    endpoint,    \n    integrationType,\n  }\n}",
			result: "{\"data\":{\"printers\":[{\"name\":\"printerName\",\"endpoint\":\"EndPoint\",\"integrationType\":\"Octoprint\"}]}}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := runQuery(h, test.query)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, result, test.result)
		})
	}
}
