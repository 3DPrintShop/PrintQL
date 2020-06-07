package resolver

import (
	"context"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	graphqlupload "github.com/smithaitufe/go-graphql-upload"
	"github.com/vitiock/PrintQL/printdb"
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

func (r SchemaResolver) UploadImageToProject(ctx context.Context, args uploadImageToProjectArgs) (*graphql.ID, error) {
	fmt.Println("Attempting to create file for image upload")
	fmt.Printf("FileName: %s\n", args.Request.Image.FileName)
	fmt.Printf("FilePath: %s\n", args.Request.Image.FilePath)

	client := ctx.Value("client").(*printdb.Client)

	imageId, err := client.CreateImage(printdb.NewImageRequest{AltText: args.Request.AltText, Type: path.Ext(args.Request.Image.FilePath)})
	client.AssociateImageWithProject(printdb.AssociateImageWithProjectRequest{ProjectId: string(args.ProjectID), ImageId: imageId, Type: path.Ext(args.Request.Image.FileName)})

	rd, err := args.Request.Image.CreateReadStream()
	if err != nil {
		fmt.Println(err.Error())
	}
	if rd != nil {
		args.Request.Image.WriteFile("./media/" + imageId + path.Ext(args.Request.Image.FileName))
	} else {
		fmt.Println("failure to create reader for component")
	}
	return nil, nil
}
