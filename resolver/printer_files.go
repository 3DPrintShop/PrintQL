package resolver

import (
	"context"
	"github.com/vitiock/go-octoprint"
)

// FilmResolver resolves the Film type.
type PrinterFileResolver struct {
	File *octoprint.FileInformation
}

type NewPrinterFilesArgs struct {
	APIKey   string
	Endpoint string
}

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

// ID resolves the film's unique identifier.
func (r *PrinterFileResolver) Hash() string {
	return r.File.Hash
}

// Episode resolves the episode number of this film.
func (r *PrinterFileResolver) Path() string {
	return r.File.Path
}

func (r *PrinterFileResolver) EstimatedPrintTime() float64 {
	return r.File.GCodeAnalysis.EstimatedPrintTime
}

func (r *PrinterFileResolver) LastPrintTime() float64 {
	return r.File.Print.Last.PrintTime
}

func (r *PrinterFileResolver) AveragePrintTime() float64 {
	return r.File.Print.Last.PrintTime
}
