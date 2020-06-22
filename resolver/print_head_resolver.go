package resolver

import "context"

type PrintHeadResolver struct {
	actual *float64
	target *float64
	name   string
}

type NewPrintHeadArgs struct {
	actual *float64
	target *float64
	name   string
}

func NewPrintHeadResolver(ctx context.Context, args NewPrintHeadArgs) (*PrintHeadResolver, error) {
	return &PrintHeadResolver{name: args.name, actual: args.actual, target: args.target}, nil
}

func (r PrintHeadResolver) Name() string {
	return r.name
}

func (r PrintHeadResolver) Actual() *float64 {
	return r.actual
}

func (r PrintHeadResolver) Target() *float64 {
	if *r.target == 0 {
		return nil
	}
	return r.target
}
