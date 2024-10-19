package artemia

/*
These unit tests cover the main functionalities of the artemia library:

1. Creating a new PrevalenceLayer instance
2. Set and Get operations
3. Deleting keys
4. Creating indexes
5. Querying by index
6. Generic queries
7. Transactions
8. Saving and loading data

Some points to note:

1. These tests create temporary files to store the data. The files are removed after each test using `defer os.Remove(filename)`.
2. The `TestStruct` structure is used to test the storage and retrieval of complex objects.
3. The `TestSaveAndLoad` test verifies that both data and indexes are correctly saved and reloaded.
4. The tests use simple assertions with `t.Fatal` to report failures. In a production environment, you might consider using a more robust assertion library like `testify`.
5. The `TestNewPrevalenceLayer` test ensures that a new PrevalenceLayer can be created without errors.
6. `TestSetAndGet` verifies that values can be stored and retrieved correctly.
7. `TestDelete` checks if keys can be properly deleted from the storage.
8. `TestCreateIndex` ensures that indexes can be created on specific fields of stored objects.
9. `TestQueryByIndex` tests the ability to query data using the created indexes.
10. `TestQuery` verifies that generic queries using custom filters work as expected.
11. `TestTransaction` checks if multiple operations can be performed atomically within a transaction.
12. `TestSaveAndLoad` is a comprehensive test that verifies if all data and indexes can be correctly persisted to disk and reloaded.
*/

import (
	"encoding/gob"
	"os"
	"reflect"
	"testing"
)

type TestStruct struct {
	ID   int
	Name string
}

func init() {
	gob.Register(TestStruct{})
}

func TestNewPrevalenceLayer(t *testing.T) {
	filename := "test_artemia.dat"
	defer os.Remove(filename)
	defer os.Remove(filename + ".index")

	pl, err := NewPrevalenceLayer(filename)
	if err != nil {
		t.Fatalf("Error creating PrevalenceLayer: %v", err)
	}
	if pl == nil {
		t.Fatal("PrevalenceLayer is nil")
	}
}

func TestSetAndGet(t *testing.T) {
	filename := "test_artemia.dat"
	defer os.Remove(filename)
	defer os.Remove(filename + ".index")

	pl, _ := NewPrevalenceLayer(filename)

	err := pl.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Error setting value: %v", err)
	}

	value, ok := pl.Get("key1")
	if !ok {
		t.Fatal("Key not found")
	}
	if value != "value1" {
		t.Fatalf("Expected 'value1', got '%v'", value)
	}
}

func TestDelete(t *testing.T) {
	filename := "test_artemia.dat"
	defer os.Remove(filename)
	defer os.Remove(filename + ".index")

	pl, _ := NewPrevalenceLayer(filename)

	pl.Set("key1", "value1")
	err := pl.Delete("key1")
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}

	_, ok := pl.Get("key1")
	if ok {
		t.Fatal("Key should have been deleted")
	}
}

func TestCreateIndex(t *testing.T) {
	filename := "test_artemia.dat"
	defer os.Remove(filename)
	defer os.Remove(filename + ".index")

	pl, _ := NewPrevalenceLayer(filename)

	pl.CreateIndex(reflect.TypeOf(TestStruct{}), "Name")

	typeInfo := getTypeInfo(reflect.TypeOf(TestStruct{}))
	if _, ok := pl.indexes[typeInfo]["Name"]; !ok {
		t.Fatal("Index not created")
	}
}

func TestQueryByIndex(t *testing.T) {
	filename := "test_artemia.dat"
	defer os.Remove(filename)
	defer os.Remove(filename + ".index")

	pl, _ := NewPrevalenceLayer(filename)

	pl.CreateIndex(reflect.TypeOf(TestStruct{}), "Name")

	pl.Set("1", TestStruct{ID: 1, Name: "Alice"})
	pl.Set("2", TestStruct{ID: 2, Name: "Bob"})
	pl.Set("3", TestStruct{ID: 3, Name: "Alice"})

	results := pl.QueryByIndex(reflect.TypeOf(TestStruct{}), "Name", "Alice")
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}

func TestQuery(t *testing.T) {
	filename := "test_artemia.dat"
	defer os.Remove(filename)
	defer os.Remove(filename + ".index")

	pl, _ := NewPrevalenceLayer(filename)

	pl.Set("1", TestStruct{ID: 1, Name: "Alice"})
	pl.Set("2", TestStruct{ID: 2, Name: "Bob"})
	pl.Set("3", TestStruct{ID: 3, Name: "Charlie"})

	filter := CreateFieldFilter("ID", 2)
	results := pl.Query(filter)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].(TestStruct).Name != "Bob" {
		t.Fatalf("Expected 'Bob', got '%s'", results[0].(TestStruct).Name)
	}
}

func TestTransaction(t *testing.T) {
	filename := "test_artemia.dat"
	defer os.Remove(filename)
	defer os.Remove(filename + ".index")

	pl, _ := NewPrevalenceLayer(filename)

	tx := pl.BeginTransaction()
	tx.Set(pl, "key1", "value1")
	tx.Set(pl, "key2", "value2")
	err := pl.Commit(tx)

	if err != nil {
		t.Fatalf("Error committing transaction: %v", err)
	}

	value1, ok1 := pl.Get("key1")
	value2, ok2 := pl.Get("key2")

	if !ok1 || value1 != "value1" || !ok2 || value2 != "value2" {
		t.Fatal("Transaction did not commit correctly")
	}
}

func TestSaveAndLoad(t *testing.T) {
	filename := "test_artemia.dat"
	defer os.Remove(filename)
	defer os.Remove(filename + ".index")

	pl, _ := NewPrevalenceLayer(filename)

	pl.Set("key1", "value1")
	pl.Set("key2", TestStruct{ID: 1, Name: "Alice"})

	pl.CreateIndex(reflect.TypeOf(TestStruct{}), "Name")

	// Force save
	err := pl.save()
	if err != nil {
		t.Fatalf("Error saving PrevalenceLayer: %v", err)
	}

	// Create a new PrevalenceLayer to load the saved data
	pl2, err := NewPrevalenceLayer(filename)
	if err != nil {
		t.Fatalf("Error loading PrevalenceLayer: %v", err)
	}

	value1, ok1 := pl2.Get("key1")
	if !ok1 || value1 != "value1" {
		t.Fatal("Failed to load simple value")
	}

	value2, ok2 := pl2.Get("key2")
	if !ok2 {
		t.Fatal("Failed to load struct value")
	}

	testStruct, ok := value2.(TestStruct)
	if !ok {
		t.Fatalf("Loaded value is not a TestStruct. Got type: %T", value2)
	}

	if testStruct.Name != "Alice" || testStruct.ID != 1 {
		t.Fatalf("Loaded struct does not match. Got: %+v", testStruct)
	}

	// Check if index was loaded
	typeInfo := getTypeInfo(reflect.TypeOf(TestStruct{}))
	if _, ok := pl2.indexes[typeInfo]["Name"]; !ok {
		t.Fatal("Index not loaded")
	}
}
