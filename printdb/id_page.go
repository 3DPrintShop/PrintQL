package printdb

import (
	"github.com/boltdb/bolt"
)

// IdentifierPage is a set of paginated IDs that allow for looking up specific entities.
type IdentifierPage struct {
	IDs      []string
	NextPage *string
}

// GetIDsFromBaseBucket is a helper function to get the identifiers from a bucket that exists on the base level of the database.
func (c *Client) GetIDsFromBaseBucket(bucketName string, startID *string, pageSize *int) (IdentifierPage, error) {
	var idPage IdentifierPage
	var err error

	err = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		idPage, err = GetIDsFromBucket(b, startID, pageSize)
		return err
	})

	return idPage, err
}

// GetIDsFromBucket gets a list of ids from a bucket.
func GetIDsFromBucket(bucket *bolt.Bucket, startID *string, pageSize *int) (IdentifierPage, error) {
	var ids []string

	if bucket == nil {
		return IdentifierPage{
			IDs:      ids,
			NextPage: nil,
		}, nil
	}

	c := bucket.Cursor()

	count := 0
	var nextPage *string = nil

	var (
		k   []byte
		val []byte
	)

	if startID != nil {
		k, val = c.Seek([]byte(*startID))
	} else {
		k, val = c.First()
	}

	for ; k != nil; k, val = c.Next() {
		if val != nil {
			continue
		}

		if pageSize != nil && count >= *pageSize {
			next := string(k)
			nextPage = &next
			break
		}
		count++
		ids = append(ids, string(k))
	}

	return IdentifierPage{
		IDs:      ids,
		NextPage: nextPage,
	}, nil
}
