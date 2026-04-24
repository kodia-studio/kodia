package kodia

import "fmt"

// Resolve retrieves a typed dependency from the container.
// Returns (value, ok) — ok is false if key not found or value is not of type T.
func Resolve[T any](app *App, key string) (T, bool) {
	val, ok := app.Get(key)
	if !ok {
		var zero T
		return zero, false
	}
	typed, ok := val.(T)
	return typed, ok
}

// MustResolve retrieves a typed dependency from the container or panics.
// Panics if key not found or value is not of type T with a descriptive message.
func MustResolve[T any](app *App, key string) T {
	val, ok := Resolve[T](app, key)
	if !ok {
		var zero T
		panic(fmt.Sprintf("container: key %q not found or wrong type (expected %T)", key, zero))
	}
	return val
}
