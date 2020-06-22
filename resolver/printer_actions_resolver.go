package resolver

import (
	"context"
	"fmt"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
	"github.com/vitiock/go-octoprint"
)

type JobAction string

const (
	START  JobAction = "START"
	CANCEL JobAction = "CANCEL"
	PAUSE  JobAction = "PAUSE"
	RESUME JobAction = "RESUME"
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

// CreatePrinter creates a printer and returns a resolver to the newly created printer.
func (r PrinterActionsResolver) CreatePrinter(ctx context.Context, args CreatePrintersQueryArgs) (*PrinterResolver, error) {
	client := ctx.Value("client").(*printdb.Client)

	printerID, err := client.CreatePrinter(printdb.NewPrinterRequest{Endpoint: args.Endpoint, Name: args.Name, APIKey: args.APIKey, IntegrationType: args.IntegrationType})

	if err != nil {
		return nil, err
	}

	return NewPrinter(ctx, NewPrinterArgs{ID: printerID})
}

type SelectFileQueryArgs struct {
	FilePath string
}

func (r PrinterActionsResolver) SelectFile(ctx context.Context, args SelectFileQueryArgs) (*string, error) {
	printer, err := loader.LoadPrinter(ctx, string(*r.printerID))

	if err != nil {
		return nil, err
	}

	client := octoprint.NewClient(printer.Endpoint, printer.APIKey)

	fileRequest := octoprint.SelectFileRequest{Path: args.FilePath, Location: octoprint.Local, Print: false}
	err = fileRequest.Do(client)

	if err != nil {
		return nil, err
	}

	return &args.FilePath, nil
}

type SendJobActionQueryArgs struct {
	Action *JobAction
}

func (r PrinterActionsResolver) SendJobAction(ctx context.Context, args SendJobActionQueryArgs) (*PrintJobResolver, error) {
	printer, err := loader.LoadPrinter(ctx, string(*r.printerID))

	if err != nil {
		return nil, err
	}

	client := octoprint.NewClient(printer.Endpoint, printer.APIKey)

	switch *args.Action {
	case START:
		request := octoprint.StartRequest{}
		request.Do(client)
	case PAUSE:
		request := octoprint.PauseRequest{Action: octoprint.Pause}
		request.Do(client)
	case RESUME:
		request := octoprint.PauseRequest{Action: octoprint.Resume}
		request.Do(client)
	case CANCEL:
		request := octoprint.CancelRequest{}
		request.Do(client)
	}

	return NewPrintJobResolver(ctx, NewPrintJobArgs{PrinterID: string(*r.printerID)})
}
