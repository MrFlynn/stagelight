package database

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
	"reflect"
)

// Color contains color information for the bands.
type Color struct {
	ID       uint8    `json:"id"`
	Name     string   `json:"name"`
	Sequence []uint32 `json:"sequence"`
}

var defaultColor = Color{
	ID:       0,
	Name:     "Default (Off)",
	Sequence: []uint32{0},
}

func (c Color) size() (int, bool) {
	// This gets the size of the array in bytes.
	// Since only 24 bits of each uint32 is used to store color
	// information (R, G, B are 1 byte each) the number of bytes is 3/4 the
	// total number of bytes.
	calculatedSize := len(c.Sequence) * 3

	return calculatedSize, calculatedSize < 129 // Allows for 43 colors in the sequence.
	// The limitation of this is to prevent overflow of the packet radio buffer. Each
	// transmission can only hold up to 221 bytes including the frame header. This gives plenty
	// of room for future functionality.
}

// MarshalJSON implements the JSON interface for Color.
func (c *Color) MarshalJSON() ([]byte, error) {
	type Alias Color
	aux := struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	return json.Marshal(aux)
}

// UnmarshalJSON implements the JSON interface for Color.
func (c *Color) UnmarshalJSON(b []byte) error {
	type Alias Color
	aux := struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	_, ok := c.size()
	if !ok {
		return fmt.Errorf("The total number of colors allowed in a sequence has been exceeded")
	}

	return nil
}

// ColorController is a struct containing information on the Colors bucket.
type ColorController struct {
	tableName string `default:"Colors"`
}

// Initialize creates the Colors bucket required for this controller.
func (cc *ColorController) Initialize(db *bolt.DB) error {
	cc.tableName = GetName(cc)

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(GetName(cc)))
		if err != nil {
			return fmt.Errorf("Could not create bucked with name %s", GetName(cc))
		}

		return nil
	})

	b, _ := json.Marshal(defaultColor) // This should not error.
	err = cc.Add(db, b)
	if err != nil {
		return fmt.Errorf("Could not initialize database with default color scheme")
	}

	return err
}

// Get gets a single Color from the database based on given ID.
func (cc *ColorController) Get(db *bolt.DB, identifier interface{}) (interface{}, error) {
	color := &Color{}

	id, ok := identifier.(uint8)
	if !ok {
		return color, fmt.Errorf("Could not convert id from interface to uint8")
	}

	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(GetName(cc))).Cursor()

		_, v := c.Seek([]byte{id})
		if v == nil {
			log.Printf("Could not find color with ID %d", id)
			return fmt.Errorf("Could not find color with ID %d", id)
		}

		return json.Unmarshal(v, color)
	})

	color.ID = id
	return reflect.ValueOf(color).Interface(), err
}

// GetAll method gets a list of all color sequences from the database.
func (cc *ColorController) GetAll(db *bolt.DB) ([]interface{}, error) {
	i := make([]interface{}, 0)
	colors := &i

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetName(cc)))

		err := b.ForEach(func(k, v []byte) error {
			var color Color

			err := json.Unmarshal(v, &color)
			if err != nil {
				log.Printf("Could not get color with ID %d", k[0])
				return err
			}

			*colors = append(*colors, reflect.ValueOf(color).Interface())
			return nil
		})

		return err
	})

	return *colors, err
}

// Add method adds a single color to the database.
func (cc *ColorController) Add(db *bolt.DB, color []byte) error {
	c := Color{}

	err := json.Unmarshal(color, &c)
	if err != nil {
		return fmt.Errorf("Could not process JSON request. Please check request format")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetName(cc)))

		err := b.Put([]byte{c.ID}, color)
		if err != nil {
			log.Printf("Could not update color with ID %d and sequence %v", c.ID, c.Sequence)
			return err
		}

		log.Printf("Successfully updated color with ID %d", c.ID)
		return nil
	})
}

// AddMultiple adds multiple color sequences to the database.
func (cc *ColorController) AddMultiple(db *bolt.DB, colors []byte) error {
	c := []Color{}

	err := json.Unmarshal(colors, &c)
	if err != nil {
		return fmt.Errorf("Could not process JSON request. Please check request format")
	}

	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetName(cc)))

		for _, color := range c {
			colorAt, err := json.Marshal(&color)
			if err != nil {
				return err
			}

			err = b.Put([]byte{color.ID}, colorAt)
			if err != nil {
				log.Printf("Unable to update color with ID %d and sequence %v", color.ID, color.Sequence)
				return err
			}

			log.Printf("Successfully update color with ID %d", color.ID)
		}

		return nil
	})
}
