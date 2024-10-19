package main

import (
	"encoding/gob"
	"fmt"
	"reflect"

	"github.com/sea-monkeys/artemia"
)

type User struct {
	ID   string
	Name string
	Age  int
}

func init() {
	gob.Register(User{})
}

func main() {
	// Initialize the prevalence layer
	pl, err := artemia.NewPrevalenceLayer("users.gob")
	if err != nil {
		panic(err)
	}

	// Add some users
	users := []User{
		{ID: "1", Name: "Alice", Age: 30},
		{ID: "2", Name: "Bob", Age: 25},
		{ID: "3", Name: "Charlie", Age: 35},
		{ID: "4", Name: "David", Age: 28},
	}

	for _, user := range users {
		if err := pl.Set(user.ID, user); err != nil {
			panic(err)
		}
	}

	// Create an index on the Age field
	pl.CreateIndex(reflect.TypeOf(User{}), "Age")

	// Using pl.Query to display all users
	fmt.Println("All users (pl.Query):")
	allUsers := pl.Query(func(item interface{}) bool {
		_, ok := item.(User)
		return ok
	})

	for _, u := range allUsers {
		user := u.(User)
		fmt.Printf("ID: %s, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}

	// Using pl.QueryByIndex to display users of a specific age
	fmt.Println("\nUsers aged 28 (pl.QueryByIndex):")
	usersAge28 := pl.QueryByIndex(reflect.TypeOf(User{}), "Age", 28)

	for _, u := range usersAge28 {
		user := u.(User)
		fmt.Printf("ID: %s, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}

	// Example of using pl.Query with a custom filter
	fmt.Println("\nUsers over 30 years old (pl.Query with custom filter):")
	olderUsers := pl.Query(func(item interface{}) bool {
		user, ok := item.(User)
		return ok && user.Age > 30
	})

	for _, u := range olderUsers {
		user := u.(User)
		fmt.Printf("ID: %s, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}
