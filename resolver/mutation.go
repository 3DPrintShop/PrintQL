package resolver

import (
	"context"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/boltdb/bolt"
	"github.com/graph-gophers/graphql-go"
	graphqlupload "github.com/smithaitufe/go-graphql-upload"
	"io/ioutil"
)

// CreatePrintersQueryArgs are the parameters passed in as part of a create printer query that specify details about the printer.
type CreatePrintersQueryArgs struct {
	Name     string
	APIKey   string
	Endpoint string
}

// CreatePrinter creates a printer and returns a resolver to the newly created printer.
func (r SchemaResolver) CreatePrinter(ctx context.Context, args CreatePrintersQueryArgs) (*PrinterResolver, error) {
	client := ctx.Value("client").(*printdb.Client)

	printerID, err := client.CreatePrinter(printdb.NewPrinterRequest{Endpoint: args.Endpoint, Name: args.Name, APIKey: args.APIKey})

	if err != nil {
		return nil, err
	}

	return NewPrinter(ctx, NewPrinterArgs{ID: printerID})
}

// DeletePrintersQueryArgs are the parameters passed in as part of a delete printer query that specify which printer to delete.
type DeletePrintersQueryArgs struct {
	ID graphql.ID
}

// DeletePrinter deletes a printer and returns the id of the printer deleted.
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

// CreateProjectQueryArgs are the required args passed in as part of a create project mutation.
type CreateProjectQueryArgs struct {
	Name string
}

// CreateProject creates a new project and returns a resolver to that project.
func (r SchemaResolver) CreateProject(ctx context.Context, args CreateProjectQueryArgs) (*ProjectResolver, error) {
	client := ctx.Value("client").(*printdb.Client)

	projectID, err := client.CreateProject(args.Name)

	if err != nil {
		return nil, err
	}

	return NewProject(ctx, NewProjectArgs{ID: projectID})
}

type uploadComponentArgs struct {
	ProjectID graphql.ID
	Component graphqlupload.GraphQLUpload
}

// UploadComponent creates a component in the printdb, then saves it to a file with a matching id, and returns that id.
func (r SchemaResolver) UploadComponent(ctx context.Context, args uploadComponentArgs) (*graphql.ID, error) {
	client := ctx.Value("client").(*printdb.Client)

	componentID, err := client.CreateComponent(printdb.NewComponentRequest{Name: args.Component.FileName})
	if err != nil {
		return nil, err
	}
	err = client.AssociateComponentWithProject(string(args.ProjectID), componentID)
	if err != nil {
		return nil, err
	}

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
		args.Component.WriteFile("./uploads/" + componentID + ".stl")
	}
	return nil, nil
}

// FilamentActions is a resolver to mutations specifically for manipulating filament.
func (r SchemaResolver) FilamentActions(ctx context.Context) (*FilamentActionsResolver, error) {
	return NewFilamentActionsResolver()
}
