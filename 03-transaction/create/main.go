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

	// create indexes for the Human structure
	pl.CreateIndex(reflect.TypeOf(Human{}), "City")
	pl.CreateIndex(reflect.TypeOf(Human{}), "Age")

	// Create a transaction
	tr := pl.BeginTransaction()

	// Add cmds to the transaction
	tr.Set(pl, "1", Human{FirstName: "Alice", LastName: "Cooper", Age: 30, City: "Paris"})
	tr.Set(pl, "2", Human{FirstName: "Bob", LastName: "Morane", Age: 25, City: "Lyon"})
	tr.Set(pl, "3", Human{FirstName: "Charlie", LastName: "Brown", Age: 35, City: "Paris"})

	//tr.Delete(pl, "3")

	// Execute the transaction
	err = pl.Commit(tr)
	if err != nil {
		fmt.Println("Error when executing the transaction:", err)
	}

}
