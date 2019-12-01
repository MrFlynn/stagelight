package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	bolt "go.etcd.io/bbolt"
)

// VoteResult contains the results from voting in the database.
type VoteResult struct {
	Positive int `json:"positive"`
	Negative int `json:"negative"`
}

// VoteController is a struct containing information on the Votes bucket.
type VoteController struct {
	tableName string `default:"Votes"`
}

// Initialize sets the fields of the passed VoteController and create the required bucket.
func (vc *VoteController) Initialize(db *bolt.DB) error {
	vc.tableName = vc.GetName()

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(vc.GetName()))
		if err != nil {
			return fmt.Errorf("Could not create bucket with name %s", vc.GetName())
		}

		return nil
	})

	// Initialize table with 0 values.
	b, _ := json.Marshal(VoteResult{0, 0})
	err = vc.Add(db, b)
	if err != nil {
		return fmt.Errorf("Could not set intial values for bucket")
	}

	return err
}

// Get method gets all of the votes from the database.
func (vc VoteController) Get(db *bolt.DB, _ interface{}) (interface{}, error) {
	votes := VoteResult{0, 0}

	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(vc.GetName())).Cursor()

		_, pos := c.First()
		if pos == nil {
			log.Println("Could not get positive votes")
			return errors.New("Could not get positive votes")
		}

		_, neg := c.Last()
		if neg == nil {
			log.Println("Could not get negative votes")
			return errors.New("Could not get negative votes")
		}

		votes.Positive = int(pos[0])
		votes.Negative = int(neg[0])

		return nil
	})

	return reflect.ValueOf(votes).Interface(), err
}

// GetAll method is an alias for Get method except the output is packaged into a slice.
func (vc VoteController) GetAll(db *bolt.DB) ([]interface{}, error) {
	votes, err := vc.Get(db, nil)

	votesSlice := make([]interface{}, 1)
	votesSlice[0] = votes

	return votesSlice, err
}

// GetName gets the default name of the struct.
func (vc VoteController) GetName() string {
	name := vc.tableName
	if name == "" {
		t := reflect.TypeOf(vc)
		f, _ := t.FieldByName("tableName")
		name = f.Tag.Get("default")
	}

	return name
}

// Add method updates the vote count from the input of a JSON byte stream.
func (vc VoteController) Add(db *bolt.DB, votes []byte) error {
	results := VoteResult{}

	err := json.Unmarshal(votes, &results)
	if err != nil {
		log.Println("Could not unmarshal votes")
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(vc.GetName()))

		err := b.Put([]byte{0}, []byte{uint8(results.Positive)})
		if err != nil {
			log.Println("Could not update positive vote tally")
			return err
		}

		err = b.Put([]byte{1}, []byte{uint8(results.Negative)})
		if err != nil {
			log.Println("Could not update negative vote tally")
			return err
		}

		return nil
	})
}

// AddMultiple is an alias for the Add method.
func (vc VoteController) AddMultiple(db *bolt.DB, votes []byte) error {
	return vc.Add(db, votes)
}
