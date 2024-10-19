package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

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

/*
GetBytesBody returns the body of an HTTP request as a []byte.
  - It takes a pointer to an http.Request as a parameter.
  - It returns a []byte.
*/
func GetBytesBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}

func main() {
	pl, err := artemia.NewPrevalenceLayer("data.gob")
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes for the Human structure
	pl.CreateIndex(reflect.TypeOf(Human{}), "City")
	pl.CreateIndex(reflect.TypeOf(Human{}), "Age")

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	mux := http.NewServeMux()

	fileServerHtml := http.FileServer(http.Dir("public"))
	mux.Handle("/", fileServerHtml)

	counter := 0

	mux.HandleFunc("POST /api/create/random/human", func(response http.ResponseWriter, request *http.Request) {
		fullName := randomName()
		human := Human{
			FirstName: fullName[:strings.Index(fullName, " ")],
			LastName:  fullName[strings.Index(fullName, " ")+1:],
			Age:       randomAge(),
			City:      randomCity(),
		}
		key := fmt.Sprintf("%d", counter)
		counter += 1

		err := pl.Set(key, human)
		if err != nil {
			result := fmt.Sprintf("Error setting human %s: %v", key, err)
			log.Println("😡", result)
			//response.Write([]byte(result))
			http.Error(response, "Error setting human: "+err.Error(), http.StatusInternalServerError)
			return
		}
		result := fmt.Sprintf("%s Added human: %+v\n", key, human)
		log.Println(result)
		response.Write([]byte(result))

	})

	// Do the same thing but with transactions

	mux.HandleFunc("POST /api/create/random/human/with/transaction", func(response http.ResponseWriter, request *http.Request) {
		fullName := randomName()
		human := Human{
			FirstName: fullName[:strings.Index(fullName, " ")],
			LastName:  fullName[strings.Index(fullName, " ")+1:],
			Age:       randomAge(),
			City:      randomCity(),
		}
		key := fmt.Sprintf("%d", counter)
		counter += 1

		// Create a transaction
		tr := pl.BeginTransaction()
		// Add cmd to the transaction
		tr.Set(pl, key, human)

		// Execute the transaction
		err = pl.Commit(tr)

		if err != nil {
			result := fmt.Sprintf("Error setting human %s: %v", key, err)
			log.Println("😡", result)
			//response.Write([]byte(result))
			http.Error(response, "Error setting human: "+err.Error(), http.StatusInternalServerError)
			return
		}

		result := fmt.Sprintf("%s Added human: %+v\n", key, human)
		log.Println(result)
		response.Write([]byte(result))

	})

	mux.HandleFunc("GET /api/get/humans/from/paris", func(response http.ResponseWriter, request *http.Request) {
		// Query example: find humans in Paris
		parisians := pl.QueryByIndex(reflect.TypeOf(Human{}), "City", "Paris")

		jsonList, err := json.Marshal(&parisians)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}

		response.Header().Add("Content-Type", "application/json; charset=utf-8")
		response.Write(jsonList)
		/*
			fmt.Printf("\nHumans in Paris: %d\n", len(parisians))
			for _, human := range parisians[:5] { // Print first 5 for brevity
				fmt.Printf("%+v\n", human)
			}
		*/
	})

	mux.HandleFunc("GET /api/get/humans/aged/30", func(response http.ResponseWriter, request *http.Request) {
		// Query example: find humans aged 30
		age30 := pl.QueryByIndex(reflect.TypeOf(Human{}), "Age", 30)

		jsonList, err := json.Marshal(&age30)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}

		response.Header().Add("Content-Type", "application/json; charset=utf-8")
		response.Write(jsonList)

		/*
			fmt.Printf("\nHumans aged 30: %d\n", len(age30))
			for _, human := range age30[:5] { // Print first 5 for brevity
				fmt.Printf("%+v\n", human)
			}
		*/
	})

	server := &http.Server{
		Addr:         ":" + httpPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("🌍 http server is listening on: " + httpPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", httpPort, err)
	}

	/*
		var errListening error
		log.Println("🌍 http server is listening on: " + httpPort)
		errListening = http.ListenAndServe(":"+httpPort, mux)

		log.Fatal(errListening)
	*/

}
