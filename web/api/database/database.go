package database

import (
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

var subControllers = []Controller{&DeviceController{}, &VoteController{}}

// Controller interface provides a list of methods that a controller should implement.
type Controller interface {
	Initialize(*bolt.DB) error

	// Generic get methods.
	Get(*bolt.DB, interface{}) (interface{}, error)
	GetAll(*bolt.DB) ([]interface{}, error)
	GetName() string

	// Generic add/update methods
	Add(*bolt.DB, []byte) error
	AddMultiple(*bolt.DB, []byte) error
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
		m[c.GetName()] = c
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
