package resolver_test

import (
	"fmt"
	"github.com/3DPrintShop/PrintQL/resolver"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrinterActionsResolver_TestPrinterActions(t *testing.T) {
	ctx := getContext()

	t.Run("Printer Actions", func(t *testing.T) {
		id := graphql.ID(testID)

		printerActions, err := resolver.NewPrinterActionsResolver(resolver.NewPrinterActionsArgs{PrinterID: &id})
		if err != nil {
			t.Error(err)
			return
		}

		spool, err := printerActions.LoadSpool(ctx, resolver.NewLoadSpoolArgs{SpoolID: testSpoolID})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, string(spool.ID()), testSpoolID)
	})

	t.Run("Printer Actions no printer id", func(t *testing.T) {
		printerActions, err := resolver.NewPrinterActionsResolver(resolver.NewPrinterActionsArgs{PrinterID: nil})
		if err != nil {
			t.Error(err)
			return
		}

		_, err = printerActions.LoadSpool(ctx, resolver.NewLoadSpoolArgs{SpoolID: testSpoolID})
		assert.Equal(t, err, fmt.Errorf("attempting to call load spool without a printer id"))
	})
}
