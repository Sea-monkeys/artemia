package main

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/sea-monkeys/artemia"
)

type UserSession struct {
	UserID    string
	LoginTime time.Time
	ExpiresAt time.Time
}

func init() {
	gob.Register(UserSession{})
}

func main() {
	pl, err := artemia.NewPrevalenceLayer("sessions.gob")
	if err != nil {
		panic(err)
	}

	// Function to create a new session
	createSession := func(userID string) error {
		session := UserSession{
			UserID:    userID,
			LoginTime: time.Now(),
			ExpiresAt: time.Now().Add(3 * time.Second), // The session expires in 3 seconds
		}
		return pl.Set(userID, session)
	}

	// Function to clean up expired sessions
	cleanExpiredSessions := func() {
		now := time.Now()
		sessions := pl.Query(func(item interface{}) bool {
			session, ok := item.(UserSession)
			return ok && session.ExpiresAt.Before(now)
		})

		for _, session := range sessions {
			s := session.(UserSession)
			fmt.Printf("Remove the session of %s\n", s.UserID)
			pl.Delete(s.UserID)
		}
	}

	// Create some sessions
	createSession("user1")
	createSession("user2")
	createSession("user3")

	time.Sleep(5 * time.Second) // Wait for 5 seconds

	// Create some more sessions
	createSession("user4")
	createSession("user5")

	//Clean up expired sessions
	cleanExpiredSessions()

	// Check for remaining sessions
	remainingSessions := pl.Query(func(item interface{}) bool {
		_, ok := item.(UserSession)
		return ok
	})

	fmt.Println("Remaining sessions after cleaning:")
	for _, session := range remainingSessions {
		s := session.(UserSession)
		fmt.Printf("UserID: %s, Expires: %s\n", s.UserID, s.ExpiresAt)
	}
}
