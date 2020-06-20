package printdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// Image represents an image that can be displayed to the user
type Image struct {
	ID      string
	AltText string
	Type    string
}

// Image retrieves a an image corresponding with imageID.
func (c *Client) Image(imageID string) (Image, error) {
	var image Image

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
		if b == nil {
			return fmt.Errorf("database not setup")
		}
		pv := b.Bucket([]byte(imageID))

		if pv == nil {
			return fmt.Errorf("project by that id does not exist")
		}

		image = Image{
			ID:      imageID,
			AltText: string(pv.Get([]byte(AltText))),
			Type:    string(pv.Get([]byte(Type))),
		}

		return nil
	})
	return image, err
}

// GetImages returns a paginated list of image IDs.
func (c *Client) GetImages(pageID *string) (IdentifierPage, error) {
	return c.GetIDsFromBaseBucket(ImageBucket, pageID, nil)
}

// NewImageRequest is the data needed to create a new image in the database.
type NewImageRequest struct {
	AltText *string
	Type    string
}

// CreateImage creates an image based on the data in request, and returns an ID to that image.
func (c *Client) CreateImage(request NewImageRequest) (string, error) {
	id := uuid.New()
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
		ib, err := b.CreateBucket([]byte(id.String()))
		if err != nil {
			return err
		}

		ib.Put([]byte(AltText), []byte(*request.AltText))
		ib.Put([]byte(Type), []byte(request.Type))
		return nil
	})

	return id.String(), err
}
