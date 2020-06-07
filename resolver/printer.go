package resolver

import (
	"context"
	"fmt"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
)

type PrinterResolver struct {
	Printer printdb.Printer
}

type NewPrintersArgs struct {
	ID *string
}

type NewPrinterArgs struct {
	ID string
}

type Printer struct {
	Id graphql.ID
	Alias string
	APIKey string
	Endpoint string
}

func NewPrinter(ctx context.Context, args NewPrinterArgs) (*PrinterResolver, error) {
	fmt.Printf("Printer request: %s\n", args.ID)
	printer, errs := loader.LoadPrinter(ctx, args.ID)

	return &PrinterResolver{Printer: printer}, errs
}

func NewPrinters(ctx context.Context, args NewPrintersArgs) (*[]*PrinterResolver, error) {
	if args.ID != nil {
		printer, err := NewPrinter(ctx, NewPrinterArgs{ ID: *args.ID})
		resolvers := []*PrinterResolver{printer}
		return &resolvers, err
	}

	var resolvers []*PrinterResolver
	printers, err := loader.LoadPrinters(ctx, "")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for _, printer := range printers {
		resolvers = append(resolvers, &PrinterResolver{Printer: printer})
	}

	return &resolvers, nil
}

func (r *PrinterResolver) ID() graphql.ID {
	return graphql.ID(r.Printer.ID)
}

func (r *PrinterResolver) Name() string {
	return r.Printer.Alias
}

func (r *PrinterResolver) Endpoint() string {
	return r.Printer.Endpoint
}


type PrinterFilesQueryArgs struct {
	Path *string
}

func (r *PrinterResolver) Files(ctx context.Context, args PrinterFilesQueryArgs) (*[]*PrinterFileResolver, error) {
	return NewPrinterFiles(ctx, NewPrinterFilesArgs{Endpoint: r.Printer.Endpoint, APIKey: r.Printer.APIKey})
}

func (r *PrinterResolver) State(ctx context.Context) (*PrinterStateResolver, error) {
	return NewPrinterState(ctx, NewPrinterStateArgs{Endpoint: r.Printer.Endpoint, APIKey: r.Printer.APIKey})
}