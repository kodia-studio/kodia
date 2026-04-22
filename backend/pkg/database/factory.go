package database

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// Factory defines the structure for generating model instances.
type Factory struct {
	db        *gorm.DB
	blueprints map[string]func() interface{}
}

var globalFactory *Factory

// InitFactory initializes the global factory system.
func InitFactory(db *gorm.DB) {
	globalFactory = &Factory{
		db:        db,
		blueprints: make(map[string]func() interface{}),
	}
}

// Define registers a new blueprint for a model.
func Define(model interface{}, blueprint func() interface{}) {
	if globalFactory == nil {
		return
	}
	name := reflect.TypeOf(model).String()
	globalFactory.blueprints[name] = blueprint
}

// Create generates and persists a model instance.
func Create(model interface{}) (interface{}, error) {
	if globalFactory == nil {
		return nil, fmt.Errorf("factory not initialized")
	}

	name := reflect.TypeOf(model).String()
	blueprint, ok := globalFactory.blueprints[name]
	if !ok {
		return nil, fmt.Errorf("no blueprint defined for model %s", name)
	}

	instance := blueprint()
	
	// Persist to DB
	if err := globalFactory.db.Create(instance).Error; err != nil {
		return nil, err
	}

	return instance, nil
}

// CreateMany generates and persists multiple model instances.
func CreateMany(model interface{}, count int) ([]interface{}, error) {
	var results []interface{}
	for i := 0; i < count; i++ {
		res, err := Create(model)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}


