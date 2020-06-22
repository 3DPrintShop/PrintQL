package resolver

import (
	"context"
	"fmt"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/vitiock/go-octoprint"
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

// IntegrationType tells what integration is being used to communicate with the printer.
func (r *PrinterResolver) IntegrationType(ctx context.Context) (string, error) {
	return r.Printer.IntegrationType, nil
}

func (r *PrinterResolver) LoadedFile(ctx context.Context) (*string, error) {
	client := octoprint.NewClient(r.Printer.Endpoint, r.Printer.APIKey)

	jobRequest := octoprint.JobRequest{}
	jobStatus, err := jobRequest.Do(client)

	if err != nil {
		return nil, err
	}

	return &jobStatus.Job.File.Name, nil
}

func (r *PrinterResolver) Job(ctx context.Context) (*PrintJobResolver, error) {
	return NewPrintJobResolver(ctx, NewPrintJobArgs{PrinterID: r.Printer.ID})
}

func (r *PrinterResolver) Bed(ctx context.Context) (*PrintBedResolver, error) {
	return NewPrintBedResolver(ctx, NewPrintBedArgs{PrinterID: string(r.ID())})
}

func (r *PrinterResolver) Tools(ctx context.Context) (*[]*PrintHeadResolver, error) {
	printerState, err := loader.LoadOctoprintPrinter(ctx, r.Printer.ID)

	var resolvers []*PrintHeadResolver

	if err != nil {
		if err.Error() == "Printer is not operational" {
			return nil, nil
		}
		return nil, err
	}

	for k, v := range printerState.Temperature.Current {
		fmt.Println(k)
		if k == "bed" {
			fmt.Println("continue.")
			continue
		}

		resolver, err := NewPrintHeadResolver(ctx, NewPrintHeadArgs{name: k, target: &v.Target, actual: &v.Actual})

		if err != nil {
			continue
		}

		resolvers = append(resolvers, resolver)
	}

	return &resolvers, nil
}
