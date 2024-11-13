package optional_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/micronull/optional"
)

func TestType_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type some struct {
		Field optional.Type[string] `json:"f"`
	}

	tests := [...]struct {
		name  string
		input []byte
		want  string
	}{
		{"nil", nil, ""},
		{"empty", []byte(``), ""},
		{"empty", []byte(`null`), ""},
		{"empty", []byte(`{}`), ""},
		{"empty", []byte(`{"f":""}`), ""},
		{"null", []byte(`{"f":null}`), ""},
		{"has", []byte(`{"f":"some"}`), "some"},
		{"invalid json", []byte(`{"f:"some"`), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got some

			_ = json.Unmarshal(tt.input, &got)

			require.Equal(t, tt.want, got.Field.V)
		})
	}
}

func TestType_UnmarshalJSON_IsSetNull(t *testing.T) {
	t.Parallel()

	type some struct {
		Field optional.Type[string] `json:"f"`
	}

	tests := [...]struct {
		name  string
		input []byte
		want  assert.BoolAssertionFunc
	}{
		{"nil", nil, assert.False},
		{"empty", []byte(``), assert.False},
		{"empty", []byte(`null`), assert.False},
		{"empty", []byte(`{}`), assert.False},
		{"empty", []byte(`{"f":""}`), assert.False},
		{"has", []byte(`{"f":"some"}`), assert.False},

		{"null", []byte(`{"f":null}`), assert.True},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got some

			_ = json.Unmarshal(tt.input, &got)

			tt.want(t, got.Field.IsSetNull())
		})
	}
}

func TestType_UnmarshalJSON_IsSet(t *testing.T) {
	t.Parallel()

	type some struct {
		Field optional.Type[string] `json:"f"`
	}

	tests := [...]struct {
		name  string
		input []byte
		want  assert.BoolAssertionFunc
	}{
		{"nil", nil, assert.False},
		{"empty", []byte(``), assert.False},
		{"empty", []byte(`null`), assert.False},
		{"empty", []byte(`{}`), assert.False},

		{"empty", []byte(`{"f":""}`), assert.True},
		{"null", []byte(`{"f":null}`), assert.True},
		{"has", []byte(`{"f":"some"}`), assert.True},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got some

			_ = json.Unmarshal(tt.input, &got)

			tt.want(t, got.Field.IsSet())
		})
	}
}

func TestType_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name     string
		input    optional.Type[string]
		expected []byte
	}{
		{"normal value", optional.New("test", false), []byte(`"test"`)},
		{"null value", optional.New("", true), []byte(`null`)},
		{"null value", optional.New("some", true), []byte(`null`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestType_ChangeMarshal(t *testing.T) {
	tests := [...]struct {
		name     string
		input    optional.Type[string]
		expected []byte
	}{
		{"normal value", optional.New("test", false), []byte(`"custom marshal"`)},
		{"null value", optional.New("", true), []byte(`null`)},
		{"null value", optional.New("some", true), []byte(`null`)},
	}

	// Custom marshaling function, with fixed value.
	customMarshal := func(v any) ([]byte, error) {
		return []byte(`"custom marshal"`), nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optional.ChangeMarshal(customMarshal)
			result, err := tt.input.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)

			// Reset to default marshal
			optional.ChangeMarshal(json.Marshal)
		})
	}
}

func TestType_ChangeUnmarshal(t *testing.T) {
	type some struct {
		Field optional.Type[string] `json:"f"`
	}

	tests := [...]struct {
		name     string
		input    []byte
		expected string
	}{
		{"normal value", []byte(`{"f":"some value"}`), "custom unmarshal"},
	}

	// Custom unmarshal function to convert strings to uppercase
	customUnmarshal := func(data []byte, v any) error {
		assert.Equal(t, []byte(`"some value"`), data)

		v = "custom unmarshal"

		return nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optional.ChangeUnmarshal(customUnmarshal)

			got := some{}

			err := json.Unmarshal(tt.input, &got)
			require.NoError(t, err)

			// Reset to default unmarshal
			optional.ChangeUnmarshal(json.Unmarshal)
		})
	}
}

func TestType_ChangeUnmarshal_Error(t *testing.T) {
	type some struct {
		Field optional.Type[string] `json:"f"`
	}

	errExpect := errors.New("some error")

	customUnmarshal := func([]byte, any) error {
		return errExpect
	}

	optional.ChangeUnmarshal(customUnmarshal)

	got := some{}

	err := json.Unmarshal([]byte(`{"f":"some"}`), &got)
	require.ErrorIs(t, err, errExpect)

	// Reset to default unmarshal
	optional.ChangeUnmarshal(json.Unmarshal)
}

func TestType_UnmarshalJSON_ChangeNull(t *testing.T) {
	t.Parallel()

	type some struct {
		Field optional.Type[string] `json:"f"`
	}

	got := some{}

	jsString := `{"f":"some"}`
	jsNull := `{"f":null}`
	jsEmpty := `{}`

	_ = json.Unmarshal([]byte(jsString), &got)

	assert.Equal(t, "some", got.Field.V)
	assert.True(t, got.Field.IsSet())
	assert.False(t, got.Field.IsSetNull())

	_ = json.Unmarshal([]byte(jsNull), &got)

	assert.Equal(t, "", got.Field.V)
	assert.True(t, got.Field.IsSet())
	assert.True(t, got.Field.IsSetNull())

	_ = json.Unmarshal([]byte(jsString), &got)

	assert.Equal(t, "some", got.Field.V)
	assert.True(t, got.Field.IsSet())
	assert.False(t, got.Field.IsSetNull())

	_ = json.Unmarshal([]byte(jsEmpty), &got)

	assert.Equal(t, "some", got.Field.V)
	assert.True(t, got.Field.IsSet())
	assert.False(t, got.Field.IsSetNull())
}
