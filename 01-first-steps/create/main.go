package main

import (
	"encoding/gob"
	"log"

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

	bob := Human{
		FirstName: "Bob",
		LastName:  "Morane",
		Age:       42,
		City:      "Lyon",
	}
	if err := pl.Set("bob", bob); err != nil {
		log.Fatal(err)
	}
}
