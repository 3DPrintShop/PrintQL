package printdb_test

import (
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_TestImageCreationAndRetrieval(t *testing.T) {
	type test struct {
		name           string
		imagesToCreate int
	}

	tests := []test{
		{name: "one image", imagesToCreate: 1},
		{name: "three images", imagesToCreate: 3},
	}

	for _, test := range tests {
		context, err := setup()
		if err != nil {
			t.Error(err)
		}

		client, err := printdb.NewClient(context.db)

		t.Run(test.name, func(t *testing.T) {
			t.Run("Create Images", func(t *testing.T) {
				for i := 0; i < test.imagesToCreate; i++ {
					altText := TestAltText
					printerID, err := client.CreateImage(printdb.NewImageRequest{
						Type:    TestComponentType,
						AltText: &altText,
					})

					if err != nil {
						t.Error(err)
					}

					assert.NotEqual(t, "", printerID)
				}
			})

			t.Run("Get Components", func(t *testing.T) {
				imagePage, err := client.GetImages(nil)

				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.imagesToCreate, len(imagePage.ImageIDs))
			})

			t.Run("Get Each Printer", func(t *testing.T) {
				imagePage, err := client.GetImages(nil)

				if err != nil {
					t.Error(err)
				}

				for _, imageID := range imagePage.ImageIDs {
					image, err := client.Image(imageID)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, imageID, image.ID)
					assert.Equal(t, TestComponentType, image.Type)
					assert.Equal(t, TestAltText, image.AltText)
				}
			})
		})

		if err != nil {
			t.Error(err)
			return
		}

		teardown(context)
	}
}
