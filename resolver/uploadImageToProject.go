package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/graphql-go"
	graphqlupload "github.com/smithaitufe/go-graphql-upload"
	"path"
)

type uploadImageRequest struct {
	AltText *string
	Image   graphqlupload.GraphQLUpload
}

type uploadImageToProjectArgs struct {
	ProjectID graphql.ID
	Request   uploadImageRequest
}

// UploadImageToProject uploads an image then associates that image with a project.
func (r SchemaResolver) UploadImageToProject(ctx context.Context, args uploadImageToProjectArgs) (*graphql.ID, error) {
	client := ctx.Value("client").(*printdb.Client)

	imageID, err := client.CreateImage(printdb.NewImageRequest{AltText: args.Request.AltText, Type: path.Ext(args.Request.Image.FileName)})
	if err != nil {
		return nil, err
	}
	err = client.AssociateImageWithProject(printdb.AssociateImageWithProjectRequest{ProjectID: string(args.ProjectID), ImageID: imageID, Type: path.Ext(args.Request.Image.FileName)})
	if err != nil {
		return nil, err
	}

	rd, err := args.Request.Image.CreateReadStream()
	if err != nil {
		// Need to delete this as part of the transaction?
		return nil, err
	}
	if rd != nil {
		args.Request.Image.WriteFile("./media/" + imageID + path.Ext(args.Request.Image.FileName))
	}
	return nil, nil
}
