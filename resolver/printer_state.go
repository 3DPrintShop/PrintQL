package resolver

import (
	"context"
	"github.com/vitiock/go-octoprint"
)

// PrinterStateResolver resolves the PrinterState type.
type PrinterStateResolver struct {
	ConnectionState octoprint.ConnectionState
	PrinterState    string

	APIKey   string
	Endpoint string
}

// NewPrinterStateArgs is a structure that represents the needed variables to access a printer's state.
type NewPrinterStateArgs struct {
	APIKey   string
	Endpoint string
}

// NewPrinterState returns a resolver that can be used to access the state of a printer using Octoprint with the specified api key and endpoint.
func NewPrinterState(ctx context.Context, args NewPrinterStateArgs) (*PrinterStateResolver, error) {
	//TODO: Move a bunch of these queries to the resolver pieces so they aren't loaded if not needed
	client := octoprint.NewClient(args.Endpoint, args.APIKey)
	cr := octoprint.ConnectionRequest{}

	resolver := PrinterStateResolver{}

	connectionResult, err := cr.Do(client)

	if err != nil {
		return &PrinterStateResolver{
			ConnectionState: octoprint.ConnectionState("INVALID"),
			PrinterState:    "INVALID",
		}, nil
	}

	resolver.ConnectionState = connectionResult.Current.State

	sr := octoprint.StateRequest{
		History: false,
	}

	stateResult, err := sr.Do(client)

	if err != nil {
		resolver.PrinterState = err.Error()
		return &resolver, nil
	}

	if stateResult != nil {
		resolver.PrinterState = stateResult.State.Text
	} else {
		resolver.PrinterState = "OFFLINE"
	}

	return &resolver, nil
}

// Connection resolves the current connection state of the printer.
func (r *PrinterStateResolver) Connection() string {
	return string(r.ConnectionState)
}

// State resolves the current state of the printer.
func (r *PrinterStateResolver) State() string {
	return r.PrinterState
}

func (r *PrinterStateResolver) LoadedFile() (string, error) {
	client := octoprint.NewClient(r.Endpoint, r.APIKey)

	jobRequest := octoprint.JobRequest{}
	jobStatus, err := jobRequest.Do(client)

	if err != nil {
		return "", err
	}

	return jobStatus.Job.File.Name, nil
}
