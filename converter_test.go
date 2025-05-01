package whiskey

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Admin bool    `json:"admin"`
	Score float64 `json:"score"`
}

func TestMapToStruct(t *testing.T) {
	tests := []struct {
		name      string
		inputMap  map[string]string
		output    any
		expectErr bool
	}{
		{
			name: "Valid input map",
			inputMap: map[string]string{
				"name":  "John",
				"age":   "30",
				"admin": "true",
				"score": "95.5",
			},
			output:    &TestStruct{},
			expectErr: false,
		},
		{
			name: "Invalid int conversion",
			inputMap: map[string]string{
				"name":  "John",
				"age":   "invalid",
				"admin": "true",
				"score": "95.5",
			},
			output:    &TestStruct{},
			expectErr: true,
		},
		{
			name: "Invalid bool conversion",
			inputMap: map[string]string{
				"name":  "John",
				"age":   "30",
				"admin": "not_bool",
				"score": "95.5",
			},
			output:    &TestStruct{},
			expectErr: true,
		},
		{
			name: "Invalid float conversion",
			inputMap: map[string]string{
				"name":  "John",
				"age":   "30",
				"admin": "true",
				"score": "not_float",
			},
			output:    &TestStruct{},
			expectErr: true,
		},
		{
			name: "Missing json tag",
			inputMap: map[string]string{
				"name":  "John",
				"age":   "30",
				"admin": "true",
				"score": "95.5",
			},
			output:    &struct{ NoTagField string }{},
			expectErr: false,
		},
		{
			name:      "Non-pointer struct",
			inputMap:  map[string]string{},
			output:    TestStruct{},
			expectErr: true,
		},
		{
			name:      "Non-struct pointer",
			inputMap:  map[string]string{},
			output:    new(int),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapToStruct(tt.inputMap, tt.output)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}

			if !tt.expectErr {
				expected := reflect.ValueOf(tt.output).Elem()
				for key, value := range tt.inputMap {
					field := expected.FieldByNameFunc(func(name string) bool {
						if field, ok := reflect.TypeOf(tt.output).Elem().FieldByName(name); ok {
							return field.Tag.Get("json") == key
						}
						return false
					})
					if field.IsValid() {
						if field.Kind() == reflect.String && field.String() != value {
							t.Errorf("Expected %s, got %s", value, field.String())
						}
					}
				}
			}
		})
	}
}
