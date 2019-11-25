package database

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
)

// DBHandler acts as an interface for bbolt.
type DBHandler struct {
	db   *bolt.DB
	path string
}

// New creates a new handler interface for bbolt.
func New(path string) DBHandler {
	exists := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exists = false
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatalf("Could not open database: %s", path)
	}

	handler := DBHandler{
		db:   db,
		path: path,
	}

	if !exists {
		err := handler.initialize()
		if err != nil {
			log.Fatal(err)
		}
	}

	return handler
}

func (handler *DBHandler) initialize() error {
	err := handler.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Devices"))
		if err != nil {
			return fmt.Errorf("Could not create bucket 'Devices'\n%s", err)
		}

		return nil
	})

	return err
}

// GetDevice gets information about specific device from the database.
func (handler *DBHandler) GetDevice(id uint8) (Device, error) {
	var device *Device

	err := handler.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("Devices")).Cursor()

		_, v := c.Seek([]byte{id})
		if v == nil {
			return fmt.Errorf("Could not find device with ID %d", id)
		}

		dev, err := create(id, v)
		if err != nil {
			return err
		}

		*device = *dev
		return nil
	})

	return *device, err
}

// GetAllDevices gets a list of all devices in the database.
func (handler *DBHandler) GetAllDevices() ([]Device, error) {
	var devices *[]Device

	err := handler.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Devices"))

		err := b.ForEach(func(k, v []byte) error {
			dev, err := create(k[0], v)
			if err != nil {
				return err
			}

			*devices = append(*devices, *dev)
			return nil
		})

		return err
	})

	return *devices, err
}

// AddDevice adds a single empty device to the database.
func (handler *DBHandler) AddDevice(id uint8) error {
	return handler.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Devices"))

		emptyDevice := Device{Id: id, Mode: Normal}
		d, err := json.Marshal(&emptyDevice)
		if err != nil {
			return err
		}

		err = b.Put([]byte{id}, d)
		return err
	})
}

// UpdateDevices updates a list of devices in the database.
func (handler *DBHandler) UpdateDevices(devices []Device) error {
	return handler.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Devices"))

		for _, device := range devices {
			d, err := json.Marshal(&device)
			if err != nil {
				return err
			}

			err = b.Put([]byte{device.Id}, d)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
