package printdb_test

import (
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestName          = "Name"
	TestAPIKey        = "APIKey"
	TestEndpoint      = "EndPoint"
	TestComponentType = "STL"
	TestComponentName = "TestName"
)

func TestClient_TestPrinterCreationAndRetrieval(t *testing.T) {
	type test struct {
		name             string
		printersToCreate int
	}

	tests := []test{
		{name: "one printer", printersToCreate: 1},
		{name: "three printers", printersToCreate: 3},
	}

	for _, test := range tests {
		context, err := setup()
		if err != nil {
			t.Error(err)
		}

		client, err := printdb.NewClient(context.db)

		if err != nil {
			t.Error(err)
			return
		}

		t.Run(test.name, func(t *testing.T) {
			t.Run("Create Printer", func(t *testing.T) {
				for i := 0; i < test.printersToCreate; i++ {
					printerID, err := client.CreatePrinter(printdb.NewPrinterRequest{
						Name:     TestName,
						Endpoint: TestEndpoint,
						APIKey:   TestAPIKey,
					})

					if err != nil {
						t.Error(err)
					}

					assert.NotEqual(t, "", printerID)
				}
			})

			t.Run("Get Printers", func(t *testing.T) {
				printerPage, err := client.Printers(nil)

				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.printersToCreate, len(printerPage.PrinterIds))
			})

			t.Run("Get Each Printer", func(t *testing.T) {
				printerPage, err := client.Printers(nil)

				if err != nil {
					t.Error(err)
				}

				for _, printerID := range printerPage.PrinterIds {
					printer, err := client.Printer(printerID)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, TestAPIKey, printer.APIKey)
					assert.Equal(t, printerID, printer.ID)
					assert.Equal(t, TestName, printer.Alias)
					assert.Equal(t, TestEndpoint, printer.Endpoint)
				}
			})
		})

		teardown(context)

	}
}

func TestClient_CreatePrinterPrinter(t *testing.T) {
	context, err := setup()
	if err != nil {
		t.Error(err)
	}

	client, err := printdb.NewClient(context.db)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = client.CreatePrinter(printdb.NewPrinterRequest{
		Name:     TestName,
		APIKey:   TestAPIKey,
		Endpoint: TestEndpoint,
	})

	if err != nil {
		defer teardown(context)
		t.Error(err)
		return
	}

	teardown(context)
}
