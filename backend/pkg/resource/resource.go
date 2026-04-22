package resource

import "reflect"

// Transformer defines the interface for resource transformation.
type Transformer interface {
	Transform(model interface{}) interface{}
}

// MapFunc is a helper for functional transformations.
type MapFunc func(model interface{}) interface{}

func (f MapFunc) Transform(model interface{}) interface{} {
	return f(model)
}

// NewItem transforms a single model.
func NewItem(model interface{}, transformer Transformer) interface{} {
	return transformer.Transform(model)
}

// NewCollection transforms a slice of models.
func NewCollection(models interface{}, transformer Transformer) []interface{} {
	v := reflect.ValueOf(models)
	if v.Kind() != reflect.Slice {
		return nil
	}

	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = transformer.Transform(v.Index(i).Interface())
	}
	return result
}
