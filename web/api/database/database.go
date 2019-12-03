package database

import (
	"fmt"
	"log"
	"os"
	"reflect"

	bolt "go.etcd.io/bbolt"
)

var subControllers = []Controller{&DeviceController{}, &VoteController{}, &ColorController{}}

// Controller interface provides a list of methods that a controller should implement.
type Controller interface {
	Initialize(*bolt.DB) error

	// Generic get methods.
	Get(*bolt.DB, interface{}) (interface{}, error)
	GetAll(*bolt.DB) ([]interface{}, error)

	// Generic add/update methods
	Add(*bolt.DB, []byte) error
	AddMultiple(*bolt.DB, []byte) error
}

// GetName gets the default bucket name (or tableName) of the given struct.
func GetName(c Controller) string {
	v := reflect.ValueOf(c)
	f := v.Elem().FieldByName("tableName")

	if f.IsValid() {
		tableName := f.String()
		if tableName == "" {
			tf, _ := reflect.TypeOf(c).Elem().FieldByName("tableName")
			return tf.Tag.Get("default")
		}

		return tableName
	}

	return ""
}

// DBHandler acts as an interface for bbolt.
type DBHandler struct {
	db          *bolt.DB
	controllers map[string]Controller
	path        string
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

	m := map[string]Controller{}
	for _, c := range subControllers {
		m[GetName(c)] = c
	}

	handler := DBHandler{
		db:          db,
		controllers: m,
		path:        path,
	}

	if !exists {
		for _, c := range handler.controllers {
			c.Initialize(handler.db)
		}
	}

	return handler
}

// Get runs the Get interface method on the controller with the given name.
func (h *DBHandler) Get(controller string, id interface{}) (interface{}, error) {
	if c, ok := h.controllers[controller]; ok {
		return c.Get(h.db, id)
	}

	return nil, fmt.Errorf("Could not find controller with name %s", controller)
}

// GetAll runs the GetAll interface method on the controller with the given name.
func (h *DBHandler) GetAll(controller string) ([]interface{}, error) {
	if c, ok := h.controllers[controller]; ok {
		return c.GetAll(h.db)
	}

	return nil, fmt.Errorf("Could not find controller with name %s", controller)
}

// Add runs the Add interface method on the controller with the given name.
func (h *DBHandler) Add(controller string, singleton []byte) error {
	if c, ok := h.controllers[controller]; ok {
		return c.Add(h.db, singleton)
	}

	return fmt.Errorf("Could not find controller with name %s", controller)
}

// AddMultiple runs the AddMultiple interface method on the controller with the given name.
func (h *DBHandler) AddMultiple(controller string, set []byte) error {
	if c, ok := h.controllers[controller]; ok {
		return c.AddMultiple(h.db, set)
	}

	return fmt.Errorf("Could not find controller with name %s", controller)
}

// BridgePayload contains the information required by the bridge
// to send to each of the bands.
type BridgePayload struct {
	ID            uint8
	Mode          Mode
	ColorSequence []uint32
}

// CreatePayload creates a communication payload for the bridge.
func (h *DBHandler) CreatePayload(devices []Device) []BridgePayload {
	completePayload := []BridgePayload{}

	for _, device := range devices {
		color, err := h.Get("Colors", device.ColorScheme)
		if err != nil {
			log.Printf("Could not get color with ID %d", device.ColorScheme)
		}

		scheme := color.(*Color)
		payload := BridgePayload{
			ID:            device.ID,
			Mode:          device.Mode,
			ColorSequence: scheme.Sequence,
		}

		completePayload = append(completePayload, payload)
	}

	return completePayload
}
