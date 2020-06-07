package printdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

type Image struct {
	ID      string
	AltText string
	Type    string
}

func (c *Client) Image(componentId string) (Image, error) {
	var image Image

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
		if b == nil {
			return fmt.Errorf("database not setup")
		}
		pv := b.Bucket([]byte(componentId))

		if pv == nil {
			return fmt.Errorf("project by that id does not exist")
		}

		image = Image{
			ID:      componentId,
			AltText: string(pv.Get([]byte(AltText))),
			Type:    string(pv.Get([]byte(Type))),
		}

		fmt.Println(image)

		return nil
	})
	fmt.Println(image)
	return image, err
}

type NewImageRequest struct {
	AltText *string
	Type    string
}

func (c *Client) CreateImage(request NewImageRequest) (string, error) {
	id := uuid.New()
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
		ib, err := b.CreateBucket([]byte(id.String()))
		if err != nil {
			fmt.Printf("Error creating component bucket")
			return err
		}

		ib.Put([]byte(AltText), []byte(*request.AltText))
		ib.Put([]byte(Type), []byte(request.Type))
		return nil
	})

	return id.String(), err
}
