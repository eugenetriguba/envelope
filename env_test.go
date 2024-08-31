package envelope

import (
	"os"
	"testing"
	"time"

	"github.com/eugenetriguba/checkmate/assert"
	"github.com/eugenetriguba/checkmate/check"
)

type TestStruct struct {
	String     string  `env:"TEST_STRING"`
	Int        int     `env:"TEST_INT"`
	Uint       uint    `env:"TEST_UINT"`
	Bool       bool    `env:"TEST_BOOL"`
	Float      float64 `env:"TEST_FLOAT"`
	Untagged   string
	Unexported string `env:"TEST_UNEXPORTED"`
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("TEST_STRING", "hello")
	os.Setenv("TEST_INT", "-42")
	os.Setenv("TEST_UINT", "42")
	os.Setenv("TEST_BOOL", "true")
	os.Setenv("TEST_FLOAT", "3.14")

	defer func() {
		os.Unsetenv("TEST_STRING")
		os.Unsetenv("TEST_INT")
		os.Unsetenv("TEST_UINT")
		os.Unsetenv("TEST_BOOL")
		os.Unsetenv("TEST_FLOAT")
		os.Unsetenv("TEST_UNEXPORTED")
	}()

	var ts TestStruct
	err := LoadFromEnv(&ts)

	check.Nil(t, err)
	check.Equal(t, ts.String, "hello")
	check.Equal(t, ts.Int, -42)
	check.Equal(t, ts.Uint, uint(42))
	check.True(t, ts.Bool)
	check.Equal(t, ts.Float, 3.14)
	check.Equal(t, ts.Untagged, "")
	check.Equal(t, ts.Unexported, "")
}

func TestLoadFromEnvErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		envVars     map[string]string
		expectedErr string
	}{
		{
			name:        "Non-pointer input",
			input:       TestStruct{},
			expectedErr: "ptr must be a pointer to a struct",
		},
		{
			name:        "Pointer to non-struct",
			input:       new(string),
			expectedErr: "ptr must be a pointer to a struct",
		},
		{
			name: "Invalid int",
			input: &struct {
				Int int `env:"TEST_INT"`
			}{},
			envVars:     map[string]string{"TEST_INT": "not an int"},
			expectedErr: "error while setting field Int from environment variable TEST_INT",
		},
		{
			name: "Invalid uint",
			input: &struct {
				Uint uint `env:"TEST_UINT"`
			}{},
			envVars:     map[string]string{"TEST_UINT": "-42"},
			expectedErr: "error while setting field Uint from environment variable TEST_UINT",
		},
		{
			name: "Invalid bool",
			input: &struct {
				Bool bool `env:"TEST_BOOL"`
			}{},
			envVars:     map[string]string{"TEST_BOOL": "not a bool"},
			expectedErr: "error while setting field Bool from environment variable TEST_BOOL",
		},
		{
			name: "Invalid float",
			input: &struct {
				Float float64 `env:"TEST_FLOAT"`
			}{},
			envVars:     map[string]string{"TEST_FLOAT": "not a float"},
			expectedErr: "error while setting field Float from environment variable TEST_FLOAT",
		},
		{
			name: "Unsupported type",
			input: &struct {
				Unsupported time.Time `env:"TEST_UNSUPPORTED"`
			}{},
			envVars:     map[string]string{"TEST_UNSUPPORTED": "2023-08-31T12:34:56Z"},
			expectedErr: "error while setting field Unsupported from environment variable TEST_UNSUPPORTED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			err := LoadFromEnv(tt.input)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
