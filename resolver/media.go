package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
)

// MediaResolver resolves the media data object.
type MediaResolver struct {
	Media printdb.Image
}

// NewMediaArgs are the arguments passed in to identify a media object.
type NewMediaArgs struct {
	ID string
}

// NewMedia takes in NewMediaArgs and returns a media resolver.
func NewMedia(ctx context.Context, args NewMediaArgs) (*MediaResolver, error) {
	media, errs := loader.LoadMedia(ctx, args.ID)

	return &MediaResolver{Media: media}, errs
}

// ID resolves the media's identifier.
func (r *MediaResolver) ID() graphql.ID {
	return graphql.ID(r.Media.ID)
}

// AltText resolves the alt text to use for an image.
func (r *MediaResolver) AltText() string {
	return r.Media.AltText
}

// Type resolves the type of media the image is.
func (r *MediaResolver) Type() string {
	return r.Media.Type
}

// Path resolves the relative url for the media file.
func (r *MediaResolver) Path() string {
	return "/media/" + r.Media.ID + r.Media.Type
}
