package printdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// Printer represents a 3d printer.
type Printer struct {
	ID            string `json:"id"`
	Alias         string `json:"alias"`
	APIKey        string `json:"apiKey"`
	Endpoint      string `json:"endpoint"`
	LoadedSpoolID string
}

// Printer retrieves a printer by id from the database.
func (c *Client) Printer(printerID string) (Printer, error) {
	//Load printer from boltdb here
	var printer Printer

	c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PrinterBucket))
		if b == nil {
			return fmt.Errorf("database not setup")
		}
		pv := b.Bucket([]byte(printerID))

		if pv == nil {
			return fmt.Errorf("printer by that id does not exist")
		}

		printer = Printer{
			ID:            printerID,
			Alias:         string(pv.Get([]byte(Alias))),
			APIKey:        string(pv.Get([]byte(APIKey))),
			Endpoint:      string(pv.Get([]byte(Endpoint))),
			LoadedSpoolID: string(pv.Get([]byte(LoadedSpool))),
		}

		return nil
	})

	return printer, nil
}

// Printers gets a paginated list of printer IDs.
func (c *Client) Printers(pageID *string) (IdentifierPage, error) {
	return c.GetIDsFromBaseBucket(PrinterBucket, pageID, nil)
}

// NewPrinterRequest is the required data to create a new printer.
type NewPrinterRequest struct {
	Name     string
	APIKey   string
	Endpoint string
}

// CreatePrinter creates a printer in the database.
func (c *Client) CreatePrinter(request NewPrinterRequest) (string, error) {
	id := uuid.New()

	c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PrinterBucket))
		pv, err := b.CreateBucket([]byte(id.String()))

		if err != nil {
			return err
		}

		pv.Put([]byte(Alias), []byte(request.Name))
		pv.Put([]byte(APIKey), []byte(request.APIKey))
		pv.Put([]byte(Endpoint), []byte(request.Endpoint))

		return nil
	})

	return id.String(), nil
}

// LoadSpoolInPrinter associates a spool of filament with a printer so that print operations adjust it's weight, and only gcode for that filament can be printed.
func (c *Client) LoadSpoolInPrinter(printerID string, spoolID string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PrinterBucket))
		pb := b.Bucket([]byte(printerID))

		if pb == nil {
			return fmt.Errorf("no printer exists by id: %s", printerID)
		}

		pb.Put([]byte(LoadedSpool), []byte(spoolID))
		return nil
	})
}
