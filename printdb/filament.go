package printdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// Filament stores information about a specific brand of filament.
type FilamentBrand struct {
	// ID is the identifier for the filament.
	ID string
	// Name is the display name for the filament.
	Name string
	// PurchaseLink is a url to purchase more.
	PurchaseLink string
	// StartWeight is the starting weight in grams of the filament when new.
	StartWeight int
	// SpoolWeight is the weight of an empty spool of this filament in grams.
	SpoolWeight int
}

// FilamentSpool represents a spool of filament of a specific brand.
type FilamentSpool struct {
	// ID is the identifier for this instance of a filament.
	ID string
	// FilamentBrand is the identifier for the type of filament it is.
	FilamentBrand string
	// RemainingWeight is the weight in grams of the current spool
	RemainingWeight int
}

// CreateFilamentBrand creates a new filament brand entry.
func (c *Client) CreateFilamentBrand(name string) (string, error) {
	id := uuid.New()

	c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentBrandBucket))
		pv, err := b.CreateBucket([]byte(id.String()))

		if err != nil {
			return err
		}

		pv.Put([]byte(Alias), []byte(name))

		return nil
	})

	return id.String(), nil
}

// SetFilamentStartWeight sets the start weight for a brand of filament
func (c *Client) SetFilamentStartWeight(id string, weight int) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentBrandBucket))
		fb := b.Bucket([]byte(id))

		return fb.Put([]byte(StartingWeight), itob(weight))
	})
}

// FilamentBrandPage is a paginated list of filament brand identifiers, as well as the identifier needed to get the next page if it exists.
type FilamentBrandPage struct {
	//FilamentBrandIDs is the list of filament brand ids that are part of the page.
	FilamentBrandIDs []string
	//NextPage is the identifier used to get the next page of identifiers.
	NextPage *string
}

// GetFilamentBrands returns a paginated set of identifiers for filament brands, and takes in an identifier to get subsequent pages, that identifier is returned from this function.
func (c *Client) GetFilamentBrands(nextPageId *string) (FilamentBrandPage, error) {
	var filamentIDs []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentBrandBucket))
		c := b.Cursor()

		if nextPageId != nil {
			c.Seek([]byte(*nextPageId))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			filamentIDs = append(filamentIDs, string(k))
		}

		return nil
	})

	return FilamentBrandPage{FilamentBrandIDs: filamentIDs, NextPage: nil}, err
}

// GetFilamentBrand returns details about a filament brand given it's id.
func (c *Client) GetFilamentBrand(id string) (FilamentBrand, error) {
	var filamentBrand FilamentBrand
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentBrandBucket))
		fb := b.Bucket([]byte(id))
		filamentBrand = FilamentBrand{
			ID:          id,
			Name:        string(fb.Get([]byte(Alias))),
			StartWeight: btoi(fb.Get([]byte(StartingWeight))),
		}

		return nil
	})

	return filamentBrand, err
}
