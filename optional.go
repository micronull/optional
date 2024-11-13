// Package optional provides a generic type that can represent values which may or may not be set,
// including the concept of null. This allows handling scenarios where a value might be missing
// or explicitly set to null in JSON.
package optional

import "encoding/json"

var (
	marshaller   = json.Marshal
	unmarshaller = json.Unmarshal
)

// ChangeMarshal allows you to change the function used for marshalling.
// By default, it uses [json.Marshal]. You can provide an alternative implementation,
// such as from a library like https://pkg.go.dev/github.com/json-iterator/go.
func ChangeMarshal(m func(v any) ([]byte, error)) {
	marshaller = m
}

// ChangeUnmarshal allows you to change the function used for unmarshalling.
// By default, it uses [json.Unmarshal]. You can provide an alternative implementation,
// such as from a library like https://pkg.go.dev/github.com/json-iterator/go.
func ChangeUnmarshal(u func(data []byte, v any) error) {
	unmarshaller = u
}

// Type represents a generic value that may or may not be set and could also be null.
type Type[T any] struct {
	V T    // V holds the actual value of type T.
	n bool // n indicates if the value is explicitly null.
	s bool // s indicates if the value has been set (either to a non-null value or explicitly to null).
}

// New creates a new instance of [Type] with the specified value and null status.
func New[T any](value T, null bool) Type[T] {
	return Type[T]{
		V: value,
		n: null,
	}
}

// IsSetNull checks if the value is explicitly set to null.
func (t Type[T]) IsSetNull() bool {
	return t.n
}

// IsSet checks if the value has been set, either to a non-null value or explicitly to null.
func (t Type[T]) IsSet() bool {
	return t.s
}

var (
	_ json.Unmarshaler = (*Type[any])(nil)
	_ json.Marshaler   = (*Type[any])(nil)
)

// UnmarshalJSON implements the [json.Unmarshaler] interface for [Type].
// It handles unmarshalling JSON data into a [Type] instance, distinguishing between unset values,
// null values, and actual non-null values.
func (t *Type[T]) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil // Treat empty input as not setting the value
	}

	var zero T

	t.V = zero  // Reset value
	t.s = true  // Mark as set since we're processing data
	t.n = false // Reset null flag

	if string(bytes) == "null" {
		t.n = true // Explicitly null case

		return nil
	}

	// Otherwise, unmarshal into the actual value
	return unmarshaller(bytes, &t.V)
}

// MarshalJSON implements the [json.Marshaler] interface for [Type].
// It handles marshalling a [Type] instance to JSON, correctly representing unset values as empty,
// null values as `null`, and non-null values using the specified marshaller.
func (t Type[T]) MarshalJSON() ([]byte, error) {
	if t.n {
		return []byte(`null`), nil // Explicitly return 'null' if set to null
	}

	// Use the current marshaller for non-null values
	return marshaller(t.V)
}
