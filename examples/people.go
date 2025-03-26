// Code generated within main for Person. DO NOT EDIT.
package main

import "time"

const (
	PersonPerson1ID = "person-1"
	PersonPerson2ID = "person-2"
)

var PersonPerson1 = Person{
	Address: Address{
		City:    "Boston",
		Country: "USA",
		State:   "MA",
		Street:  "123 Main St",
		ZipCode: "02108",
	},
	BirthDate: time.Date(1985, time.June, 15, 0, 0, 0, 0, time.UTC),
	Email:     "john.doe@example.com",
	FirstName: "John",
	ID:        "person-1",
	LastName:  "Doe",
}
var PersonPerson2 = Person{
	Address: Address{
		City:    "San Francisco",
		Country: "USA",
		State:   "CA",
		Street:  "456 Oak Ave",
		ZipCode: "94107",
	},
	BirthDate: time.Date(1990, time.August, 22, 0, 0, 0, 0, time.UTC),
	Email:     "jane.smith@example.com",
	FirstName: "Jane",
	ID:        "person-2",
	LastName:  "Smith",
}
var AllPersons = []Person{PersonPerson1, PersonPerson2}
