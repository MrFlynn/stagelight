package database

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	bolt "go.etcd.io/bbolt"
)

// Mode enumerates device modes.
type Mode int

const (
	// Normal indicates that the device is operating in light output mode.
	Normal = iota
	// Vote indicates that the device is currently in a voting mode.
	Vote
)

// Device contains device-specific information about individual bands.
type Device struct {
	ID     uint8    `json:"id"`
	Mode   Mode     `json:"mode"`
	Colors []uint32 `json:"colors"`
}

// MarshallJSON is the interface method for json.Marshall for the Device struct.
func (d Device) MarshallJSON() ([]byte, error) {
	return json.Marshal(d)
}

// UnmarshalJSON is the interface method for json.Unmarshal for the Device struct.
// This method also provides some validation.
func (d *Device) UnmarshalJSON(b []byte) error {
	type Alias Device
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	if !(d.Mode == 0) && !(d.Mode == 1) {
		return fmt.Errorf("Mode %d is not a valid mode. Must be either 0 or 1", d.Mode)
	}

	return nil
}

// DeviceController is a struct containing bucket information as well as required interface methods.
type DeviceController struct {
	tableName string `default:"Devices"`
}

// Initialize creates the bucket required for the "Devices" controller.
func (dc *DeviceController) Initialize(db *bolt.DB) error {
	dc.tableName = "Devices"

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(GetName(dc)))
		if err != nil {
			return fmt.Errorf("Could not create bucket with name %s", dc.tableName)
		}

		return nil
	})

	return err
}

// Get method gets a single device from the Devices bucket that matches the passed integer ID.
func (dc *DeviceController) Get(db *bolt.DB, identifier interface{}) (interface{}, error) {
	device := &Device{}

	id, ok := identifier.(uint8)
	if !ok {
		return device, fmt.Errorf("Could not convert id from interface{} to uint8")
	}

	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(GetName(dc))).Cursor()

		_, v := c.Seek([]byte{id})
		if v == nil {
			log.Printf("Could not find device with ID %d", id)
			return fmt.Errorf("Could not find device with ID %d", id)
		}

		err := json.Unmarshal(v, device)
		if err != nil {
			return err
		}

		return nil
	})

	return reflect.ValueOf(device).Interface(), err
}

// GetAll method gets a list of all devices.
func (dc *DeviceController) GetAll(db *bolt.DB) ([]interface{}, error) {
	i := make([]interface{}, 0)
	devices := &i

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetName(dc)))

		err := b.ForEach(func(k, v []byte) error {
			var dev Device

			err := json.Unmarshal(v, &dev)
			if err != nil {
				log.Printf("Could not get device with ID %d", k[0])
				return err
			}

			*devices = append(*devices, reflect.ValueOf(dev).Interface())
			return nil
		})

		return err
	})

	return *devices, err
}

// Add method adds a single device based on the inputs from a JSON byte stream.
func (dc *DeviceController) Add(db *bolt.DB, dev []byte) error {
	device := Device{}
	err := json.Unmarshal(dev, &device)
	if err != nil {
		return fmt.Errorf("Could not process JSON request. Please check format")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetName(dc)))

		d, err := json.Marshal(&device)
		if err != nil {
			return err
		}

		err = b.Put([]byte{device.ID}, d)
		if err != nil {
			log.Printf("Could not create device with ID %d", device.ID)
			return err
		}

		log.Printf("Succesfully create device with ID %d", device.ID)
		return nil
	})
}

// AddMultiple method adds multiple devices from a JSON byte stream.
func (dc *DeviceController) AddMultiple(db *bolt.DB, devs []byte) error {
	devices := []Device{}
	err := json.Unmarshal(devs, &devices)
	if err != nil {
		return fmt.Errorf("Could not process JSON request. Please check format")
	}

	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetName(dc)))

		for _, device := range devices {
			d, err := json.Marshal(&device)
			if err != nil {
				return err
			}

			err = b.Put([]byte{device.ID}, d)
			if err != nil {
				log.Printf("Unable to update device with ID %d", device.ID)
				return err
			}

			log.Printf("Sucessfully updated device with ID %d", device.ID)
		}

		return nil
	})
}
