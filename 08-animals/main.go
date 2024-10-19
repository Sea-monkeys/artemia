package main

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/sea-monkeys/artemia"
)

// Animal is an interface that both Dog and Cat will implement
type Animal interface {
	GetID() string
	GetName() string
	GetSpecies() string
}

// Dog represents a dog
type Dog struct {
	ID    string
	Name  string
	Breed string
}

func (d Dog) GetID() string      { return d.ID }
func (d Dog) GetName() string    { return d.Name }
func (d Dog) GetSpecies() string { return "Dog" }

// Cat represents a cat
type Cat struct {
	ID    string
	Name  string
	Color string
}

func (c Cat) GetID() string      { return c.ID }
func (c Cat) GetName() string    { return c.Name }
func (c Cat) GetSpecies() string { return "Cat" }

func init() {
	// Register types for gob encoding
	gob.Register(Dog{})
	gob.Register(Cat{})
}

func main() {
	// Initialize the prevalence layer
	pl, err := artemia.NewPrevalenceLayer("animals.gob")
	if err != nil {
		panic(err)
	}

	// Create some dogs and cats
	animals := []Animal{
		Dog{ID: "d1", Name: "Buddy", Breed: "Labrador"},
		Dog{ID: "d2", Name: "Max", Breed: "German Shepherd"},
		Cat{ID: "c1", Name: "Whiskers", Color: "Orange"},
		Cat{ID: "c2", Name: "Luna", Color: "Black"},
		Dog{ID: "d3", Name: "Charlie", Breed: "Beagle"},
		Cat{ID: "c3", Name: "Milo", Color: "White"},
	}

	// Store animals in the prevalence layer
	for _, animal := range animals {
		if err := pl.Set(animal.GetID(), animal); err != nil {
			panic(err)
		}
	}

	// Query all dogs
	fmt.Println("All Dogs:")
	dogs := queryAnimalsBySpecies(pl, "Dog")
	printAnimals(dogs)

	// Query all cats
	fmt.Println("\nAll Cats:")
	cats := queryAnimalsBySpecies(pl, "Cat")
	printAnimals(cats)

	// Demonstrate querying by other criteria (e.g., name containing 'a')
	fmt.Println("\nAnimals with 'a' in their name:")
	animalsWithA := queryAnimalsByNameContaining(pl, "a")
	printAnimals(animalsWithA)
}

func queryAnimalsBySpecies(pl *artemia.PrevalenceLayer, species string) []Animal {
	results := pl.Query(func(item interface{}) bool {
		if animal, ok := item.(Animal); ok {
			return animal.GetSpecies() == species
		}
		return false
	})

	animals := make([]Animal, len(results))
	for i, result := range results {
		animals[i] = result.(Animal)
	}
	return animals
}

func queryAnimalsByNameContaining(pl *artemia.PrevalenceLayer, substring string) []Animal {
	results := pl.Query(func(item interface{}) bool {
		if animal, ok := item.(Animal); ok {
			return containsIgnoreCase(animal.GetName(), substring)
		}
		return false
	})

	animals := make([]Animal, len(results))
	for i, result := range results {
		animals[i] = result.(Animal)
	}
	return animals
}

func printAnimals(animals []Animal) {
	for _, animal := range animals {
		switch a := animal.(type) {
		case Dog:
			fmt.Printf("Dog: ID=%s, Name=%s, Breed=%s\n", a.ID, a.Name, a.Breed)
		case Cat:
			fmt.Printf("Cat: ID=%s, Name=%s, Color=%s\n", a.ID, a.Name, a.Color)
		}
	}
}

func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
