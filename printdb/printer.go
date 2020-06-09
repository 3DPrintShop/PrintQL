package printdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

type Printer struct {
	ID       string `json:"id"`
	Alias    string `json:"alias"`
	APIKey   string `json:"apiKey"`
	Endpoint string `json:"endpoint"`
}

func (c *Client) Printer(printerId string) (Printer, error) {
	//Load printer from boltdb here
	var printer Printer

	c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PrinterBucket))
		if b == nil {
			return fmt.Errorf("database not setup")
		}
		pv := b.Bucket([]byte(printerId))

		if pv == nil {
			return fmt.Errorf("printer by that id does not exist")
		}

		printer = Printer{
			ID:       printerId,
			Alias:    string(pv.Get([]byte(Alias))),
			APIKey:   string(pv.Get([]byte(APIKey))),
			Endpoint: string(pv.Get([]byte(Endpoint))),
		}

		return nil
	})

	return printer, nil
}

type PrinterPage struct {
	PrinterIds []string
	NextKey    *string
}

func (c *Client) Printers(pageId *string) (PrinterPage, error) {

	var printerIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PrinterBucket))
		c := b.Cursor()

		if pageId != nil {
			c.Seek([]byte(*pageId))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			printerIds = append(printerIds, string(k))
		}

		return nil
	})

	return PrinterPage{PrinterIds: printerIds, NextKey: nil}, err
}

type NewPrinterRequest struct {
	Name     string
	APIKey   string
	Endpoint string
}

func (c *Client) CreatePrinter(request NewPrinterRequest) (string, error) {
	fmt.Printf("Name: %s, apiKey: %s, endpoint: %s\n", request.Name, request.APIKey, request.Endpoint)

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
