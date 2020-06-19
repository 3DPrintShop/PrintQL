package printdb

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
)

const (
	PrinterBucket = "Printers"
	APIKey        = "APIKey"
	Endpoint      = "Endpoint"
	Alias         = "Alias"
	LoadedSpool   = "LoadedSpool"

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

	FilamentBrandBucket = "FilamentBrand"
	StartingWeight      = "StartingWeight"
	PurchaseLink        = "PurchaseLink"

	FilamentSpoolBucket = "FilamentSpool"
	BrandID             = "BrandID"
	RemainingWeight     = "RemainingWeight"
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
		_, err = tx.CreateBucketIfNotExists([]byte(FilamentBrandBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(FilamentSpoolBucket))
		if err != nil {
			return err
		}
		return nil
	})
}

func itob(v int32) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int32 {
	if b == nil {
		return 0
	}
	return int32(binary.BigEndian.Uint64(b))
}
