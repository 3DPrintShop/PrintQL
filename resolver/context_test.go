package resolver_test

import (
	"context"
	"fmt"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/dataloader"
)

const (
	testID           = "testID"
	testName         = "Test Name"
	testAltText      = "alt text what it is"
	testType         = "Jiff"
	testPurchaseLink = "purchase link"
	testStartWeight  = 15
	testSpoolWeight  = 25
)

var (
	testMetadata = printdb.ProjectMetadata{
		Public: false,
	}
	componentIDs   = []string{"component1", "component2", "component3"}
	testComponents = printdb.ProjectComponentPage{
		ComponentIds: componentIDs,
	}

	imageIDs   = []string{"media1", "media2", "media3"}
	testImages = printdb.ProjectMediaPage{
		MediaIds: imageIDs,
	}
)

type mockPrintDB struct {
	loaded bool
}

var brandIDS = []string{"brand1", "brand2", "brand3", "brand4", "brand5"}

func (mock mockPrintDB) GetFilamentBrands(id *string) (printdb.IdentifierPage, error) {
	return printdb.IdentifierPage{
		IDs: brandIDS,
	}, nil
}

func (mock mockPrintDB) GetFilamentBrand(id string) (printdb.FilamentBrand, error) {
	return printdb.FilamentBrand{
		ID:           id,
		Name:         testName,
		PurchaseLink: testPurchaseLink,
		StartWeight:  testStartWeight,
		SpoolWeight:  testSpoolWeight,
	}, nil
}

func (mock mockPrintDB) Component(componentId string) (printdb.Component, error) {
	return printdb.Component{
		ID:   componentId,
		Name: testName,
		Type: testType,
	}, nil
}

func (mock mockPrintDB) Image(mediaId string) (printdb.Image, error) {
	return printdb.Image{
		ID:      mediaId,
		AltText: testAltText,
		Type:    testType,
	}, nil
}

func (mock mockPrintDB) GetProject(id string) (printdb.Project, error) {
	if !mock.loaded {
		mock.loaded = true
		return printdb.Project{
			ID:         id,
			Name:       testName,
			Metadata:   testMetadata,
			Components: testComponents,
			Images:     testImages,
		}, nil
	}
	return printdb.Project{ID: "Loaded twice"}, fmt.Errorf("Function was called a second time")
}

func (mock mockPrintDB) GetProjects(pageId *string) (printdb.ProjectPage, error) {
	ids := []string{"test", "test2", "test3"}
	page := printdb.ProjectPage{
		ProjectIDs: ids,
		NextKey:    pageId,
	}
	return page, nil
}

func getContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, loader.ProjectLoaderKey, dataloader.NewBatchedLoader(loader.NewProjectLoader(mockPrintDB{loaded: false})))
	ctx = context.WithValue(ctx, loader.ProjectPageLoaderKey, dataloader.NewBatchedLoader(loader.NewProjectPageLoader(mockPrintDB{loaded: false})))
	ctx = context.WithValue(ctx, loader.MediaLoaderKey, dataloader.NewBatchedLoader(loader.NewMediaLoader(mockPrintDB{loaded: false})))
	ctx = context.WithValue(ctx, loader.ComponentLoaderKey, dataloader.NewBatchedLoader(loader.NewComponentLoader(mockPrintDB{loaded: false})))
	ctx = context.WithValue(ctx, loader.FilamentBrandLoaderKey, dataloader.NewBatchedLoader(loader.NewFilamentBrandLoader(mockPrintDB{loaded: false})))
	ctx = context.WithValue(ctx, loader.FilamentBrandsLoaderKey, dataloader.NewBatchedLoader(loader.NewFilamentBrandsLoader(mockPrintDB{loaded: false})))
	return ctx
}
