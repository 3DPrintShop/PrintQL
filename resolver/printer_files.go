package resolver

import (
	"context"
	"github.com/vitiock/go-octoprint"
)

// PrinterFileResolver resolves the Printer File type.
type PrinterFileResolver struct {
	File *octoprint.FileInformation
}

// NewPrinterFilesArgs are the arguments required to get the files for a printer.
type NewPrinterFilesArgs struct {
	APIKey   string
	Endpoint string
}

// NewPrinterFiles returns a list of file resolvers
func NewPrinterFiles(ctx context.Context, args NewPrinterFilesArgs) (*[]*PrinterFileResolver, error) {
	client := octoprint.NewClient(args.Endpoint, args.APIKey)
	request := octoprint.FilesRequest{
		Location:  octoprint.Local,
		Recursive: true,
	}

	result, err := request.Do(client)
	if err != nil {
		return nil, err
	}

	var resolvers = make([]*PrinterFileResolver, 0, 1)

	for _, d := range result.Files {
		if !d.IsFolder() {
			resolvers = append(resolvers, &PrinterFileResolver{File: d})
		}
	}

	return &resolvers, err
}

// Hash resolves the files hash string.
func (r *PrinterFileResolver) Hash() string {
	return r.File.Hash
}

// Path resolves the path of this file.
func (r *PrinterFileResolver) Path() string {
	return r.File.Path
}

// EstimatedPrintTime resolves the estimated time it takes to print this file.
func (r *PrinterFileResolver) EstimatedPrintTime() float64 {
	return r.File.GCodeAnalysis.EstimatedPrintTime
}

// LastPrintTime resolves the amount of time it took to print last time.
func (r *PrinterFileResolver) LastPrintTime() float64 {
	return r.File.Print.Last.PrintTime
}

// AveragePrintTime resolves what the average print time for the file was.
func (r *PrinterFileResolver) AveragePrintTime() float64 {
	return r.File.Print.Last.PrintTime
}
