package printdb

import (
	"github.com/boltdb/bolt"
)

const (
	PrinterBucket = "Printers"
	APIKey        = "APIKey"
	Endpoint      = "Endpoint"
	Alias         = "Alias"

	ProjectBucket = "Projects"
	Name          = "Name"
	Metadata      = "Metadata"

	ProjectComponentsBucket = "Components"

	ProjectImagesBucket = "Media"

	ComponentBucket = "Components"
	Type            = "Type"

	ComponentProjectBucket = "Projects"

	ImageBucket = "Media"
	AltText     = "AltText"
)

type Client struct {
	db *bolt.DB
}

func NewClient(db *bolt.DB) (*Client, error) {
	return &Client{db: db}, insureBuckets(db)
}

func insureBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(PrinterBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(ProjectBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(ComponentBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(ImageBucket))
		if err != nil {
			return err
		}
		return nil
	})
}
