package printdb

import (
	"testing"
)

const (
	TestName     = "Name"
	TestAPIKey   = "APIKey"
	TestEndpoint = "EndPoint"
	TestComponentType = "STL"
	TestComponentName = "TestName"
)

func TestClient_CreatePrinterPrinter(t *testing.T) {
	context, err := setup()
	if err != nil {
		t.Error(err)
	}

	client, err := NewClient(context.db)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = client.CreatePrinter(NewPrinterRequest{
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
