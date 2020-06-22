package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/loader"
)

type PrintJobResolver struct {
	PrinterID string
}

type NewPrintJobArgs struct {
	PrinterID string
}

func NewPrintJobResolver(ctx context.Context, args NewPrintJobArgs) (*PrintJobResolver, error) {
	return &PrintJobResolver{PrinterID: args.PrinterID}, nil
}

func (r *PrintJobResolver) FileName(ctx context.Context) (*string, error) {
	job, err := loader.LoadOctoprintJob(ctx, r.PrinterID)

	if err != nil {
		return nil, err
	}

	return &job.Job.File.Path, nil
}

func (r *PrintJobResolver) State(ctx context.Context) (*string, error) {
	job, err := loader.LoadOctoprintJob(ctx, r.PrinterID)

	if err != nil {
		return nil, err
	}

	return &job.State, nil
}

func (r *PrintJobResolver) PrintTime(ctx context.Context) (*float64, error) {
	job, err := loader.LoadOctoprintJob(ctx, r.PrinterID)

	if err != nil {
		return nil, err
	}

	if job.Job.LastPrintTime == 0 {
		return &job.Job.EstimatedPrintTime, nil
	}
	return &job.Job.LastPrintTime, nil
}

func (r *PrintJobResolver) ElapsedTime(ctx context.Context) (*float64, error) {
	job, err := loader.LoadOctoprintJob(ctx, r.PrinterID)

	if err != nil {
		return nil, err
	}

	if job.Progress.PrintTime == 0 {
		return nil, nil
	}
	return &job.Progress.PrintTime, nil
}
