package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	graphql "github.com/graph-gophers/graphql-go"
)

// PrinterResolver resolves the Printer type.
type PrinterResolver struct {
	Printer printdb.Printer
}

// NewPrintersArgs represent the arguments passed into the Printers query.
type NewPrintersArgs struct {
	ID *string
}

// NewPrinterArgs represent the required arguments needed to create a PrinterResolver.
type NewPrinterArgs struct {
	ID string
}

// Printer is a representation of a 3DPrinter.
type Printer struct {
	ID       graphql.ID
	Alias    string
	APIKey   string
	Endpoint string
}

// NewPrinter creates a new PrinterResolver.
func NewPrinter(ctx context.Context, args NewPrinterArgs) (*PrinterResolver, error) {
	printer, errs := loader.LoadPrinter(ctx, args.ID)

	return &PrinterResolver{Printer: printer}, errs
}

// NewPrinters creates a list of printer resolvers filtered by NewPrinterArgs
func NewPrinters(ctx context.Context, args NewPrintersArgs) (*[]*PrinterResolver, error) {
	if args.ID != nil {
		printer, err := NewPrinter(ctx, NewPrinterArgs{ID: *args.ID})
		resolvers := []*PrinterResolver{printer}
		return &resolvers, err
	}

	var resolvers []*PrinterResolver
	printers, err := loader.LoadPrinters(ctx, "")

	if err != nil {
		return nil, err
	}

	for _, printer := range printers {
		resolvers = append(resolvers, &PrinterResolver{Printer: printer})
	}

	return &resolvers, nil
}

// ID resolves the printer's identifier.
func (r *PrinterResolver) ID() graphql.ID {
	return graphql.ID(r.Printer.ID)
}

// Name resolves the printer's name.
func (r *PrinterResolver) Name() string {
	return r.Printer.Alias
}

// Endpoint resolves the endpoint of the printer integration.
func (r *PrinterResolver) Endpoint() string {
	return r.Printer.Endpoint
}

// PrinterFilesQueryArgs are the supported arguments on how to filter the printer files.
type PrinterFilesQueryArgs struct {
	Path *string
}

// Files creates a list of file resolvers for files that exist on the printer.
func (r *PrinterResolver) Files(ctx context.Context, args PrinterFilesQueryArgs) (*[]*PrinterFileResolver, error) {
	return NewPrinterFiles(ctx, NewPrinterFilesArgs{Endpoint: r.Printer.Endpoint, APIKey: r.Printer.APIKey})
}

// State resolves the current state of the printer and it's connections.
func (r *PrinterResolver) State(ctx context.Context) (*PrinterStateResolver, error) {
	return NewPrinterState(ctx, NewPrinterStateArgs{Endpoint: r.Printer.Endpoint, APIKey: r.Printer.APIKey})
}
