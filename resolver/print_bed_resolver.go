package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/loader"
)

type PrintBedResolver struct {
	PrinterID string
}

type NewPrintBedArgs struct {
	PrinterID string
}

func NewPrintBedResolver(ctx context.Context, args NewPrintBedArgs) (*PrintBedResolver, error) {
	return &PrintBedResolver{PrinterID: args.PrinterID}, nil
}

func (r *PrintBedResolver) Temperature(ctx context.Context) (*float64, error) {
	printerState, err := loader.LoadOctoprintPrinter(ctx, r.PrinterID)

	if err != nil {
		if err.Error() == "Printer is not operational" {
			return nil, nil
		}
		return nil, err
	}

	temp := printerState.Temperature.Current["bed"].Actual

	return &temp, nil
}

func (r *PrintBedResolver) Target(ctx context.Context) (*float64, error) {
	printerState, err := loader.LoadOctoprintPrinter(ctx, r.PrinterID)

	if err != nil {
		if err.Error() == "Printer is not operational" {
			return nil, nil
		}
		return nil, err
	}

	temp := printerState.Temperature.Current["bed"].Target

	if temp == 0 {
		return nil, nil
	}
	return &temp, nil
}
