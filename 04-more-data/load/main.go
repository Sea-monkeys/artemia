package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"reflect"

	"github.com/sea-monkeys/artemia"
)

type Human struct {
	FirstName string
	LastName  string
	Age       int
	City      string
}

func init() {
	gob.Register(Human{})
}

func main() {
	pl, err := artemia.NewPrevalenceLayer("../data.gob")
	if err != nil {
		log.Fatal(err)
	}

	peopleFromParis := pl.Query(artemia.CreateFieldFilter("City", "Paris"))
	for _, p := range peopleFromParis {
		fmt.Printf("%+v\n", p)
	}

	// Create indexes for the Human structure
	pl.CreateIndex(reflect.TypeOf(Human{}), "City")
	pl.CreateIndex(reflect.TypeOf(Human{}), "Age")

	// Query example: find humans in Paris
	parisians := pl.QueryByIndex(reflect.TypeOf(Human{}), "City", "Paris")
	fmt.Printf("\nHumans in Paris: %d\n", len(parisians))
	for _, human := range parisians[:5] { // Print first 5 for brevity
		fmt.Printf("%+v\n", human)
	}

	// Query example: find humans aged 30
	age30 := pl.QueryByIndex(reflect.TypeOf(Human{}), "Age", 30)
	fmt.Printf("\nHumans aged 30: %d\n", len(age30))
	for _, human := range age30[:5] { // Print first 5 for brevity
		fmt.Printf("%+v\n", human)
	}
}
