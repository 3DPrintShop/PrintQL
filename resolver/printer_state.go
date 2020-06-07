package resolver

import (
	"fmt"
	"context"
	"github.com/vitiock/go-octoprint"
)

// FilmResolver resolves the Film type.
type PrinterStateResolver struct {
	ConnectionState octoprint.ConnectionState
	PrinterState string
}

type NewPrinterStateArgs struct {
	APIKey string
	Endpoint string
}

func NewPrinterState(ctx context.Context, args NewPrinterStateArgs) (*PrinterStateResolver, error) {
	client := octoprint.NewClient(args.Endpoint, args.APIKey)
	cr := octoprint.ConnectionRequest{}

	resolver := PrinterStateResolver{
	}

	connectionResult, err := cr.Do(client)

	if(err != nil){
		fmt.Printf(err.Error())
		return &PrinterStateResolver{
			ConnectionState: octoprint.ConnectionState("INVALID"),
			PrinterState: "INVALID",
		}, nil
	}

	resolver.ConnectionState = connectionResult.Current.State

	sr := octoprint.StateRequest{
		History: false,
	}

	stateResult, err := sr.Do(client)

	if(err != nil){
		fmt.Printf(err.Error())
	}

	if(stateResult != nil){
		resolver.PrinterState = stateResult.State.Text
	} else {
		resolver.PrinterState = "OFFLINE"
	}

	return &resolver, nil
}

// ID resolves the film's unique identifier.
func (r *PrinterStateResolver) Connection() string {
	return string(r.ConnectionState)
}

func (r *PrinterStateResolver) State() string {
	return r.PrinterState
}