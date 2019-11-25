package database

import "encoding/json"

import "fmt"

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
	Id     uint8   `json:"-"`
	Mode   Mode    `json:"mode"`
	Colors []uint8 `json:"colors"`
}

func create(id uint8, buf []byte) (*Device, error) {
	var deviceInfo Device

	err := json.Unmarshal(buf, &deviceInfo)
	if err != nil {
		return &Device{}, fmt.Errorf("Could not unmarshall device info for device with ID:%d", id)
	}

	deviceInfo.id = id
	return &deviceInfo, nil
}
