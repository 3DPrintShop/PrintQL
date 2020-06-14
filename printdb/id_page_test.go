package printdb_test

import (
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testBucket = "TestBucket"

func TestClient_TestIDPage(t *testing.T) {
	type test struct {
		name                string
		identifiersToCreate int
		pageSize            *int
	}

	oneHundred := 100

	tests := []test{
		{name: "1 Identifier", identifiersToCreate: 1, pageSize: nil},
		{name: "100 Identifiers", identifiersToCreate: 100, pageSize: nil},
		{name: "5000 identifiers, 100 page size", identifiersToCreate: 5000, pageSize: &oneHundred},
		{name: "6666 identifiers, 100 page size", identifiersToCreate: 6666, pageSize: &oneHundred},
	}

	for _, test := range tests {
		context, err := setup()
		if err != nil {
			t.Error(err)
		}

		t.Run(test.name, func(t *testing.T) {
			err := context.db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucket([]byte(testBucket))

				if err != nil {
					t.Error(err)
				}

				for i := 0; i < test.identifiersToCreate; i++ {
					b.CreateBucket([]byte(string(i)))
				}
				return nil
			})

			if err != nil {
				t.Error(err)
			}

			t.Run("Query bucket for ids", func(t *testing.T) {
				err := context.db.View(func(tx *bolt.Tx) error {
					totalIdentifiers := 0

					idPage, err := printdb.GetIdsFromBucket(tx.Bucket([]byte(testBucket)), nil, test.pageSize)
					for {
						if err != nil {
							t.Error(err)
							return nil
						}

						if test.pageSize != nil {
							if test.identifiersToCreate-totalIdentifiers > *test.pageSize {
								assert.Equal(t, *test.pageSize, len(idPage.IDs))
							} else {
								assert.Equal(t, test.identifiersToCreate-totalIdentifiers, len(idPage.IDs))
							}
						} else {
							assert.Equal(t, test.identifiersToCreate, len(idPage.IDs))
						}

						totalIdentifiers = totalIdentifiers + len(idPage.IDs)
						if idPage.NextPage == nil {
							break
						}
						idPage, err = printdb.GetIdsFromBucket(tx.Bucket([]byte(testBucket)), idPage.NextPage, test.pageSize)
					}

					assert.Equal(t, test.identifiersToCreate, totalIdentifiers)

					return nil
				})

				if err != nil {
					t.Error(err)
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
