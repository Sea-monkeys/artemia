package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"

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
	//rand.Seed(time.Now().UnixNano())
}

var firstNames = []string{
	"Alice", "Bob", "Charlie", "David", "Emma", "Frank", "Grace", "Henry", "Ivy", "Jack",
	"Kate", "Liam", "Mia", "Noah", "Olivia", "Paul", "Quinn", "Rachel", "Sam", "Tina",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
	"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
}

var cities = []string{
	"Paris", "Lyon", "Marseille", "Toulouse", "Nice", "Nantes", "Strasbourg", "Montpellier", "Bordeaux", "Lille",
	"Rennes", "Reims", "Saint-Etienne", "Toulon", "Le Havre", "Grenoble", "Dijon", "Angers", "Nîmes", "Villeurbanne",
}

func randomName() string {
	return firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
}

func randomAge() int {
	return rand.Intn(60) + 18 // Ages between 18 and 77
}

func randomCity() string {
	return cities[rand.Intn(len(cities))]
}

func main() {
	pl, err := artemia.NewPrevalenceLayer("../data.gob")
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes for the Human structure
	pl.CreateIndex(reflect.TypeOf(Human{}), "City")
	pl.CreateIndex(reflect.TypeOf(Human{}), "Age")

	// Generate and add 500 Humans
	for i := 1; i <= 500; i++ {
		fullName := randomName()
		human := Human{
			FirstName: fullName[:strings.Index(fullName, " ")],
			LastName:  fullName[strings.Index(fullName, " ")+1:],
			Age:       randomAge(),
			City:      randomCity(),
		}
		key := fmt.Sprintf("%d", i)
		err := pl.Set(key, human)
		if err != nil {
			log.Printf("Error setting human %s: %v", key, err)
		} else {
			fmt.Printf("%s Added human: %+v\n", key, human)
		}
	}

}
