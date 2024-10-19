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

	youngPeople := pl.Query(func(item interface{}) bool {
		if p, ok := item.(Human); ok {
			return p.Age < 30
		}
		return false
	})
	for _, p := range youngPeople {
		fmt.Printf("%+v\n", p)
	}

	// Query using Index
	peopleFromParis = pl.QueryByIndex(reflect.TypeOf(Human{}), "City", "Paris")
	for _, p := range peopleFromParis {
		fmt.Printf("Person from Paris: %+v\n", p.(Human))
	}

}
