package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/boltdb/bolt"
	"github.com/graph-gophers/graphql-go"
	graphqlupload "github.com/smithaitufe/go-graphql-upload"
	"io/ioutil"
)

type CreatePrintersQueryArgs struct {
	Name     string
	ApiKey   string
	Endpoint string
}

func (r SchemaResolver) CreatePrinter(ctx context.Context, args CreatePrintersQueryArgs) (*PrinterResolver, error) {
	client := ctx.Value("client").(*printdb.Client)

	printerID, err := client.CreatePrinter(printdb.NewPrinterRequest{Endpoint: args.Endpoint, Name: args.Name, APIKey: args.ApiKey})

	if err != nil {
		return nil, err
	}

	return NewPrinter(ctx, NewPrinterArgs{ID: printerID})
}

type DeletePrintersQueryArgs struct {
	ID graphql.ID
}

func (r SchemaResolver) DeletePrinter(ctx context.Context, args DeletePrintersQueryArgs) (*graphql.ID, error) {
	db := ctx.Value("db").(*bolt.DB)

	return &args.ID, db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Printers"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("Printers"))
			if err != nil {
				return nil
			}
		}

		return b.DeleteBucket([]byte(args.ID))
	})
}

type CreateProjectQueryArgs struct {
	Name string
}

func (r SchemaResolver) CreateProject(ctx context.Context, args CreateProjectQueryArgs) (*ProjectResolver, error) {
	client := ctx.Value("client").(*printdb.Client)

	projectId, err := client.CreateProject(args.Name)

	if err != nil {
		return nil, err
	}

	return NewProject(ctx, NewProjectArgs{ID: projectId})
}

type uploadComponentArgs struct {
	ProjectID graphql.ID
	Component graphqlupload.GraphQLUpload
}

func (r SchemaResolver) UploadComponent(ctx context.Context, args uploadComponentArgs) (*graphql.ID, error) {
	client := ctx.Value("client").(*printdb.Client)

	componentId, err := client.CreateComponent(printdb.NewComponentRequest{Name: args.Component.FileName})
	client.AssociateComponentWithProject(string(args.ProjectID), componentId)

	rd, err := args.Component.CreateReadStream()
	if err != nil {
		return nil, err
	}
	if rd != nil {
		b2, err := ioutil.ReadAll(rd)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(args.Component.FileName, b2[:], 0666)

		// method 2: using WriteFile function. Easily write to any location in the local file system
		args.Component.WriteFile("./uploads/" + componentId + ".stl")
	}
	return nil, nil
}

func (r SchemaResolver) FilamentActions(ctx context.Context) (*FilamentActionsResolver, error) {
	return NewFilamentActionsResolver()
}
