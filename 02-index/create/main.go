package main

import (
	"encoding/gob"
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

	// create indexes for the Human structure
	pl.CreateIndex(reflect.TypeOf(Human{}), "City")
	pl.CreateIndex(reflect.TypeOf(Human{}), "Age")

	pl.Set("1", Human{FirstName: "Alice", LastName: "Cooper", Age: 30, City: "Paris"})
	pl.Set("2", Human{FirstName: "Bob", LastName: "Morane", Age: 25, City: "Lyon"})
	pl.Set("3", Human{FirstName: "Charlie", LastName: "Brown", Age: 35, City: "Paris"})

}
