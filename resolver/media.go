package resolver

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	"github.com/vitiock/PrintQL/loader"
	"github.com/vitiock/PrintQL/printdb"
)

type MediaResolver struct {
	Media printdb.Image
}

type NewMediaArgs struct {
	ID string
}

func NewMedia(ctx context.Context, args NewMediaArgs) (*MediaResolver, error) {
	media, errs := loader.LoadMedia(ctx, args.ID)

	return &MediaResolver{Media: media}, errs
}

func (r *MediaResolver) ID() graphql.ID {
	return graphql.ID(r.Media.ID)
}

func (r *MediaResolver) AltText() string {
	return r.Media.AltText
}

func (r *MediaResolver) Type() string {
	return r.Media.Type
}

func (r *MediaResolver) Path() string {
	return "/media/" + r.Media.ID + "." + r.Media.Type
}
