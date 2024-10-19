package artemia

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"sync"
)

type Index struct {
	FieldName string
	Data      map[interface{}][]string
}

type Transaction struct {
	operations []func() error
}

type TypeInfo struct {
	PkgPath string
	Name    string
}

type PrevalenceLayer struct {
	data      map[string]interface{}
	indexes   map[TypeInfo]map[string]*Index
	filename  string
	indexFile string
	mu        sync.RWMutex
	txMu      sync.Mutex
}

func NewPrevalenceLayer(filename string) (*PrevalenceLayer, error) {

	gob.Register(TypeInfo{})

	pl := &PrevalenceLayer{
		data:      make(map[string]interface{}),
		indexes:   make(map[TypeInfo]map[string]*Index),
		filename:  filename,
		indexFile: filename + ".index",
	}

	if err := pl.load(); err != nil {
		return nil, err
	}

	return pl, nil
}

func (pl *PrevalenceLayer) BeginTransaction() *Transaction {
	return &Transaction{}
}

func (tx *Transaction) Set(pl *PrevalenceLayer, key string, value interface{}) *Transaction {
	tx.operations = append(tx.operations, func() error {
		oldValue, exists := pl.data[key]
		if exists {
			pl.removeFromIndexes(key, oldValue)
		}
		pl.data[key] = value
		pl.updateIndexes(key, value)
		return nil
	})
	return tx
}

func (tx *Transaction) Delete(pl *PrevalenceLayer, key string) *Transaction {
	tx.operations = append(tx.operations, func() error {
		if oldValue, exists := pl.data[key]; exists {
			pl.removeFromIndexes(key, oldValue)
			delete(pl.data, key)
		}
		return nil
	})
	return tx
}

func (pl *PrevalenceLayer) Commit(tx *Transaction) error {
	pl.txMu.Lock()
	defer pl.txMu.Unlock()

	pl.mu.Lock()
	defer pl.mu.Unlock()

	for _, op := range tx.operations {
		if err := op(); err != nil {
			return err
		}
	}

	return pl.save()
}

func (pl *PrevalenceLayer) Set(key string, value interface{}) error {
	tx := pl.BeginTransaction()
	tx.Set(pl, key, value)
	return pl.Commit(tx)
}

func (pl *PrevalenceLayer) Get(key string) (interface{}, bool) {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	value, ok := pl.data[key]
	return value, ok
}

func (pl *PrevalenceLayer) Delete(key string) error {
	tx := pl.BeginTransaction()
	tx.Delete(pl, key)
	return pl.Commit(tx)
}

func getTypeInfo(t reflect.Type) TypeInfo {
	return TypeInfo{
		PkgPath: t.PkgPath(),
		Name:    t.Name(),
	}
}

func (pl *PrevalenceLayer) CreateIndex(structType reflect.Type, field string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	typeInfo := getTypeInfo(structType)
	if _, exists := pl.indexes[typeInfo]; !exists {
		pl.indexes[typeInfo] = make(map[string]*Index)
	}

	if _, exists := pl.indexes[typeInfo][field]; !exists {
		pl.indexes[typeInfo][field] = &Index{
			FieldName: field,
			Data:      make(map[interface{}][]string),
		}
		for key, value := range pl.data {
			if reflect.TypeOf(value) == structType {
				pl.addToIndex(typeInfo, field, key, value)
			}
		}
	}
}

func (pl *PrevalenceLayer) addToIndex(typeInfo TypeInfo, field, key string, value interface{}) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if getTypeInfo(v.Type()) != typeInfo {
		return
	}
	f := v.FieldByName(field)
	if !f.IsValid() {
		return
	}
	fieldValue := f.Interface()
	index := pl.indexes[typeInfo][field]
	if keys, exists := index.Data[fieldValue]; exists {
		for _, existingKey := range keys {
			if existingKey == key {
				return // Key already in index, no need to add
			}
		}
		index.Data[fieldValue] = append(keys, key)
	} else {
		index.Data[fieldValue] = []string{key}
	}
}

func (pl *PrevalenceLayer) removeFromIndexes(key string, value interface{}) {
	typeInfo := getTypeInfo(reflect.TypeOf(value))
	if indexes, exists := pl.indexes[typeInfo]; exists {
		for field, index := range indexes {
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			f := v.FieldByName(field)
			if !f.IsValid() {
				continue
			}
			fieldValue := f.Interface()
			if keys, ok := index.Data[fieldValue]; ok {
				newKeys := make([]string, 0, len(keys)-1)
				for _, k := range keys {
					if k != key {
						newKeys = append(newKeys, k)
					}
				}
				if len(newKeys) > 0 {
					index.Data[fieldValue] = newKeys
				} else {
					delete(index.Data, fieldValue)
				}
			}
		}
	}
}

func (pl *PrevalenceLayer) updateIndexes(key string, value interface{}) {
	typeInfo := getTypeInfo(reflect.TypeOf(value))
	if indexes, exists := pl.indexes[typeInfo]; exists {
		for field := range indexes {
			pl.addToIndex(typeInfo, field, key, value)
		}
	}
}

func (pl *PrevalenceLayer) Query(filter func(interface{}) bool) []interface{} {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	var results []interface{}
	for _, value := range pl.data {
		if filter(value) {
			results = append(results, value)
		}
	}
	return results
}

func (pl *PrevalenceLayer) QueryByIndex(structType reflect.Type, field string, value interface{}) []interface{} {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	typeInfo := getTypeInfo(structType)
	var results []interface{}
	if typeIndexes, exists := pl.indexes[typeInfo]; exists {
		if index, found := typeIndexes[field]; found {
			if keys, ok := index.Data[value]; ok {
				for _, key := range keys {
					results = append(results, pl.data[key])
				}
			}
		}
	}
	return results
}

func (pl *PrevalenceLayer) save() error {
	tempDataFile := pl.filename + ".tmp"
	tempIndexFile := pl.indexFile + ".tmp"

	if err := pl.saveData(tempDataFile); err != nil {
		return err
	}
	if err := pl.saveIndexes(tempIndexFile); err != nil {
		os.Remove(tempDataFile)
		return err
	}

	if err := os.Rename(tempDataFile, pl.filename); err != nil {
		return fmt.Errorf("erreur lors du remplacement du fichier de données: %v", err)
	}
	if err := os.Rename(tempIndexFile, pl.indexFile); err != nil {
		return fmt.Errorf("erreur lors du remplacement du fichier d'index: %v", err)
	}

	return nil
}

func (pl *PrevalenceLayer) saveData(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier de données: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(pl.data); err != nil {
		return fmt.Errorf("erreur lors de l'encodage des données: %v", err)
	}

	return nil
}

func (pl *PrevalenceLayer) saveIndexes(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier d'index: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(pl.indexes); err != nil {
		return fmt.Errorf("erreur lors de l'encodage des index: %v", err)
	}

	return nil
}

func (pl *PrevalenceLayer) load() error {
	if err := pl.loadData(); err != nil {
		return err
	}
	if err := pl.loadIndexes(); err != nil {
		pl.rebuildIndexes()
	}
	return nil
}

func (pl *PrevalenceLayer) loadData() error {
	file, err := os.Open(pl.filename)
	if err != nil {
		if os.IsNotExist(err) {
			pl.data = make(map[string]interface{})
			return nil
		}
		return fmt.Errorf("erreur lors de l'ouverture du fichier de données: %v", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&pl.data); err != nil {
		return fmt.Errorf("erreur lors du décodage des données: %v", err)
	}

	return nil
}

func (pl *PrevalenceLayer) loadIndexes() error {
	file, err := os.Open(pl.indexFile)
	if err != nil {
		if os.IsNotExist(err) {
			pl.rebuildIndexes()
			return nil
		}
		return fmt.Errorf("erreur lors de l'ouverture du fichier d'index: %v", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&pl.indexes); err != nil {
		file.Close()
		os.Remove(pl.indexFile)
		return fmt.Errorf("erreur lors du décodage des index: %v", err)
	}

	return nil
}

func (pl *PrevalenceLayer) rebuildIndexes() {
	pl.indexes = make(map[TypeInfo]map[string]*Index)
	for key, value := range pl.data {
		pl.updateIndexes(key, value)
	}
	pl.saveIndexes(pl.indexFile)
}

func CreateFieldFilter(fieldName string, expectedValue interface{}) func(interface{}) bool {
	return func(item interface{}) bool {
		v := reflect.ValueOf(item)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return false
		}
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}
		return reflect.DeepEqual(field.Interface(), expectedValue)
	}
}
