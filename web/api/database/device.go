package database

import (
	"encoding/json"
	"fmt"
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
