package whiskey

import "reflect"

type Type int

const (
	TypeString Type = iota
	TypeInt
	TypeFloat
	TypeBool
)

type Value struct {
	value    any
	dataType Type
}

// A simple key-value store used for storing data for every request
type DataStore struct {
	data map[string]Value
}

func NewDataStore() *DataStore {
	return &DataStore{
		data: make(map[string]Value),
	}
}
func (ds *DataStore) Set(key string, value any) {
	dataType := reflect.TypeOf(value)
	switch dataType.Kind() {
	case reflect.String:
		ds.data[key] = Value{value: value, dataType: TypeString}
	case reflect.Int:
		ds.data[key] = Value{value: value, dataType: TypeInt}
	case reflect.Float64:
		ds.data[key] = Value{value: value, dataType: TypeFloat}
	case reflect.Bool:
		ds.data[key] = Value{value: value, dataType: TypeBool}
	}
}

func (ds *DataStore) GetString(key string) (string, bool) {
	if value, ok := ds.data[key]; ok && value.dataType == TypeString {
		return value.value.(string), true
	}
	return "", false
}

func (ds *DataStore) GetInt(key string) (int, bool) {
	if value, ok := ds.data[key]; ok && value.dataType == TypeInt {
		return value.value.(int), true
	}
	return 0, false
}

func (ds *DataStore) GetFloat(key string) (float64, bool) {
	if value, ok := ds.data[key]; ok && value.dataType == TypeFloat {
		return value.value.(float64), true
	}
	return 0, false
}

func (ds *DataStore) GetBool(key string) (bool, bool) {
	if value, ok := ds.data[key]; ok && value.dataType == TypeBool {
		return value.value.(bool), true
	}
	return false, false
}

func (ds *DataStore) Get(key string) (any, bool) {
	if value, ok := ds.data[key]; ok {
		return value.value, true
	}
	return nil, false
}

func (ds *DataStore) Delete(key string) {
	delete(ds.data, key)
}

func (ds *DataStore) Clear() {
	ds.data = make(map[string]Value)
}

func (ds *DataStore) GetAll() map[string]Value {
	return ds.data
}
