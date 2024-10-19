# 🦐 Artemia: A Delightful Golang Prevalence Layer

Welcome to Artemia, the Go library that makes data persistence as easy as keeping sea monkeys alive! (Actually, it's much easier than that.)

## What is Artemia?

Artemia is a lightweight, in-memory object prevalence layer for Go applications. It's inspired by the concept of prevalence layers (check out [Prevayler](https://prevayler.org/) and the [Prevalence Layer pattern](https://wiki.c2.com/?PrevalenceLayer) for more info). 

Think of Artemia as a magical aquarium where your data objects live and thrive, with the added bonus of being able to survive power outages (unlike real sea monkeys).

## 🌟 Features

- In-memory data storage for lightning-fast operations
- Transparent persistence to disk
- Support for transactions
- Indexing for speedy queries
- No complex setup or external dependencies

## 🚀 Quick Start

First, install Artemia:

```bash
go get github.com/sea-monkeys/artemia
```

Now, let's dive into some code!

### 🏊‍♂️ Example 1: Adding a Human to Our Data Aquarium

```go
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
	pl, err := artemia.NewPrevalenceLayer("data.gob")
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

	// Bob is now swimming happily in our data aquarium!
}
```

### 🎣 Example 2: Fishing Out Our Human

```go
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
	pl, err := artemia.NewPrevalenceLayer("data.gob")
	if err != nil {
		log.Fatal(err)
	}

	if val, ok := pl.Get("bob"); ok {
		if human, ok := val.(Human); ok {
			fmt.Printf("Caught a wild human: %+v\n", human)
			// Output: Caught a wild human: {FirstName:Bob LastName:Morane Age:42 City:Lyon}
		}
	}
}
```

## 🧜‍♀️ Advanced Features

### Indexing: Finding Nemo (or Bob) Faster

```go
pl.CreateIndex(reflect.TypeOf(Human{}), "City")
lyonnais := pl.QueryByIndex(reflect.TypeOf(Human{}), "City", "Lyon")
fmt.Printf("Humans in Lyon: %+v\n", lyonnais)
```

### Transactions: Atomic Operations for Your Data Aquarium

```go
tx := pl.BeginTransaction()
tx.Set(pl, "alice", Human{FirstName: "Alice", LastName: "Wonder", Age: 28, City: "Paris"})
tx.Set(pl, "charlie", Human{FirstName: "Charlie", LastName: "Brown", Age: 8, City: "New York"})
err := pl.Commit(tx)
// If err is nil, both Alice and Charlie are now part of our data ecosystem!
```

## 🐠 Why Artemia?

1. **Simple as Sea Monkeys**: Easy to set up and use, just add water... err, data!
2. **Fast as a Shark**: In-memory operations for high-speed data access.
3. **Persistent as a Barnacle**: Your data sticks around, even when the power goes out.
4. **Flexible as an Octopus**: Store any Go struct, create indexes, run transactions.

## 🦀 Contributing

We welcome contributions! Feel free to open issues, submit pull requests, or just send us your best fish puns.

## 🐙 License

Artemia is released under the MIT License. See the LICENSE file for details.

Remember: Always practice safe data persistence, and never flush your sea monkeys down the toilet!
