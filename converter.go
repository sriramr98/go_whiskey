package whiskey

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// MapToStruct takes a map and converts it to a struct based on `json` tags
// We aren't converting map to json string and back to json since we also need to do type casting of string to the field type
func mapToStruct(m map[string]string, s interface{}) error {
	// Check if s is a pointer
	value := reflect.ValueOf(s)
	if value.Kind() != reflect.Ptr {
		return errors.New("s must be a pointer")
	}

	// Get the underlying value that the pointer points to
	value = value.Elem()

	// Check if the underlying value is a struct
	if value.Kind() != reflect.Struct {
		return errors.New("s must be a pointer to a struct")
	}

	// Loop through the fields of the struct
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		// Remove anything after a comma in the json tag (handles omitempty, etc.)
		if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
			jsonTag = jsonTag[:commaIdx]
		}

		// Check if the map has the key
		if val, ok := m[jsonTag]; ok {
			srcField := value.Field(i)
			if !srcField.CanSet() {
				continue // Skip if we can't set this field
			}

			srcValue, err := convertValueToType(val, srcField)
			if err != nil {
				return err
			}

			// Set the value in the struct
			if srcValue.IsValid() {
				srcField.Set(srcValue)
			} else {
				return errors.New("invalid value type")
			}

		}
	}

	return nil
}

func convertValueToType(value string, field reflect.Value) (reflect.Value, error) {
	switch field.Type().Kind() {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(intValue), nil
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(boolValue), nil
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(floatValue), nil
	default:
		return reflect.Value{}, errors.New("unsupported type")
	}
}
