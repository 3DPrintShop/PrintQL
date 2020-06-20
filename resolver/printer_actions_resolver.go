package resolver

import (
	"context"
	"fmt"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
)

// PrinterActionsResolver is a resolver for actions that can be taken on a printer.
type PrinterActionsResolver struct {
	printerID *graphql.ID
}

// NewPrinterActionsArgs are the arguments passed in when query for PrinterActions.
type NewPrinterActionsArgs struct {
	PrinterID *graphql.ID
}

// NewPrinterActionsResolver creates a new resolver for mutating printers.
func NewPrinterActionsResolver(args NewPrinterActionsArgs) (*PrinterActionsResolver, error) {
	resolver := PrinterActionsResolver{printerID: args.PrinterID}
	return &resolver, nil
}

// NewLoadSpoolArgs are the arguments passed into the load spool mutation.
type NewLoadSpoolArgs struct {
	SpoolID string
}

// LoadSpool associates a print with a spool of filament so future prints can be counter against it.
func (r PrinterActionsResolver) LoadSpool(ctx context.Context, args NewLoadSpoolArgs) (*FilamentSpoolResolver, error) {
	client := ctx.Value("client").(printdb.PrintDB)
	if r.printerID == nil {
		return nil, fmt.Errorf("attempting to call load spool without a printer id")
	}

	err := client.LoadSpoolInPrinter(string(*r.printerID), args.SpoolID)
	if err != nil {
		return nil, err
	}

	spoolResolver, err := NewFilamentSpool(ctx, NewFilamentSpoolArgs{ID: args.SpoolID})
	return spoolResolver, err
}
