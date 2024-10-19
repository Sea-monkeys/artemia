package main

import (
	"encoding/gob"
	"fmt"
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

	if val, ok := pl.Get("bob"); ok {
		if utilisateur, ok := val.(Human); ok {
			fmt.Printf("Human: %+v\n", utilisateur)
		}
	}
}
