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
	StartWeight int32
	// SpoolWeight is the weight of an empty spool of this filament in grams.
	SpoolWeight int32
}

// FilamentSpool represents a spool of filament of a specific brand.
type FilamentSpool struct {
	// ID is the identifier for this instance of a filament.
	ID string
	// FilamentBrand is the identifier for the type of filament it is.
	FilamentBrand string
	// RemainingWeight is the weight in grams of the current spool
	RemainingWeight int32
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
func (c *Client) SetFilamentStartWeight(id string, weight int32) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentBrandBucket))
		fb := b.Bucket([]byte(id))

		return fb.Put([]byte(StartingWeight), itob(weight))
	})
}

// SetFilamentPurchaseLink sets the purchase link for a brand of filament.
func (c *Client) SetFilamentPurchaseLink(id string, purchaseLink string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(FilamentBrandBucket)).Bucket([]byte(id)).Put([]byte(PurchaseLink), []byte(purchaseLink))
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
func (c Client) GetFilamentBrand(id string) (FilamentBrand, error) {
	var filamentBrand FilamentBrand
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentBrandBucket))
		fb := b.Bucket([]byte(id))
		filamentBrand = FilamentBrand{
			ID:           id,
			Name:         string(fb.Get([]byte(Alias))),
			StartWeight:  btoi(fb.Get([]byte(StartingWeight))),
			PurchaseLink: string(fb.Get([]byte(PurchaseLink))),
		}

		return nil
	})

	return filamentBrand, err
}

func (c *Client) CreateFilamentSpool(brandID string) (string, error) {
	id := uuid.New()

	filimentBrand, err := c.GetFilamentBrand(brandID)
	if err != nil {
		return "", err
	}

	err = c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentSpoolBucket))
		sb, err := b.CreateBucket([]byte(id.String()))

		if err != nil {
			return err
		}

		err = sb.Put([]byte(BrandID), []byte(brandID))
		if err != nil {
			return err
		}

		return sb.Put([]byte(RemainingWeight), itob(filimentBrand.StartWeight))
	})

	return id.String(), err
}

// FilamentSpoolPage is a paginated list of filament spool identifiers, as well as the identifier needed to get the next page if it exists.
type FilamentSpoolPage struct {
	//FilamentSpoolIDs is the list of filament spool ids that are part of the page.
	FilamentSpoolIDs []string
	//NextPage is the identifier used to get the next page of identifiers.
	NextPage *string
}

// GetFilamentSpools gets a paginated list of fillament spool ids.
func (c *Client) GetFilamentSpools(nextPageId *string) (FilamentSpoolPage, error) {
	var spoolIDs []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentSpoolBucket))
		c := b.Cursor()

		if nextPageId != nil {
			c.Seek([]byte(*nextPageId))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			spoolIDs = append(spoolIDs, string(k))
		}

		return nil
	})

	return FilamentSpoolPage{FilamentSpoolIDs: spoolIDs, NextPage: nil}, err
}

// GetFilamentBrand returns details about a filament brand given it's id.
func (c *Client) GetFilamentSpool(id string) (FilamentSpool, error) {
	var filamentSpool FilamentSpool
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilamentSpoolBucket))
		fs := b.Bucket([]byte(id))
		filamentSpool = FilamentSpool{
			ID:              id,
			FilamentBrand:   string(fs.Get([]byte(BrandID))),
			RemainingWeight: btoi(fs.Get([]byte(RemainingWeight))),
		}

		return nil
	})

	return filamentSpool, err
}
