package whiskey

import (
	"testing"
)

func TestDataStore(t *testing.T) {
	ds := NewDataStore()

	// Test Set and GetString
	ds.Set("key1", "value1")
	if val, ok := ds.GetString("key1"); !ok || val != "value1" {
		t.Errorf("Expected 'value1', got '%v'", val)
	}

	// Test Set and GetInt
	ds.Set("key2", 42)
	if val, ok := ds.GetInt("key2"); !ok || val != 42 {
		t.Errorf("Expected 42, got '%v'", val)
	}

	// Test Set and GetFloat
	ds.Set("key3", 3.14)
	if val, ok := ds.GetFloat("key3"); !ok || val != 3.14 {
		t.Errorf("Expected 3.14, got '%v'", val)
	}

	// Test Set and GetBool
	ds.Set("key4", true)
	if val, ok := ds.GetBool("key4"); !ok || val != true {
		t.Errorf("Expected true, got '%v'", val)
	}

	// Test Get with wrong type
	if _, ok := ds.GetInt("key1"); ok {
		t.Errorf("Expected GetInt to fail for key1")
	}

	// Test Get for non-existent key
	if _, ok := ds.Get("nonexistent"); ok {
		t.Errorf("Expected Get to fail for nonexistent key")
	}

	// Test Delete
	ds.Delete("key1")
	if _, ok := ds.Get("key1"); ok {
		t.Errorf("Expected key1 to be deleted")
	}

	// Test Clear
	ds.Set("key5", "value5")
	ds.Clear()
	if _, ok := ds.Get("key5"); ok {
		t.Errorf("Expected all keys to be cleared")
	}

	// Test GetAll
	ds.Set("key6", "value6")
	ds.Set("key7", 7)
	allData := ds.GetAll()
	if len(allData) != 2 {
		t.Errorf("Expected 2 items in the store, got %d", len(allData))
	}
	if allData["key6"].value != "value6" || allData["key7"].value != 7 {
		t.Errorf("GetAll returned incorrect data")
	}

	// Test overwriting existing key
	ds.Set("key6", "newValue")
	if val, ok := ds.GetString("key6"); !ok || val != "newValue" {
		t.Errorf("Expected 'newValue', got '%v'", val)
	}
}
