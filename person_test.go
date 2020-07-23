package room_allocation

import (
	"log"
	"testing"
)

func TestPerson0(t *testing.T) {
	samplePeople := NewPeople([]string{"bob", "fred"})
	if samplePeople.MinConnectionScore() != 0 {
		t.Error("There Should be 0 connections to start with")
	}

	_, bob := samplePeople.GetPersonByName("bob")
	_, fred := samplePeople.GetPersonByName("fred")
	bob.AddToConnection(*fred)

	connection, err := bob.GetConnection(*fred)
	if err != nil {
		t.Error("Couldn't get bob's connection to fred")
	}
	if connection.Count != 1 {
		t.Error("This connection should have a score of 1, Got:", connection)
	}
	connectionScore := samplePeople.MinConnectionScore()
	if connectionScore != 1 {
		t.Error("The Minimum Connection Score should be 1:", connectionScore)
	}
}

func TestPerson1(t *testing.T) {
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa"})
	if samplePeople.MinConnectionScore() != 0 {
		t.Error("There Should be 0 connections to start with")
	}

	_, bob := samplePeople.GetPersonByName("bob")
	_, fred := samplePeople.GetPersonByName("fred")
	_, lisa := samplePeople.GetPersonByName("Lisa")
	bob.AddToConnection(*fred)
	bob.AddToConnection(*lisa)

	connection_bob_fred, err := bob.GetConnection(*fred)
	if err != nil {
		t.Error("Couldn't get bob's connection to fred")
	}
	if connection_bob_fred.Count != 1 {
		t.Error("This connection should have a score of 1, Got:", connection_bob_fred)
	}

	connection_bob_lisa, err := bob.GetConnection(*lisa)
	if err != nil {
		t.Error("Couldn't get bob's connection to lisa")
	}

	connectionScore := samplePeople.MinConnectionScore()
	if connectionScore != 0 {
		t.Log(connection_bob_lisa)
		t.Error("The Minimum Connection Score should be 0:", connectionScore)
	}

	fred.AddToConnection(*lisa)
	connectionScore = samplePeople.MinConnectionScore()
	if connectionScore != 1 {
		t.Error("The Minimum Connection Score should be 1:", connectionScore)
	}
}

func TestGroupSelection(t *testing.T) {
	// Start off with our list of people
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa"})
	// Let's just put bob in a room
	roomWithoutBob := samplePeople.TakeOutOfRoomByName("bob")
	log.Println("Room without bob:", roomWithoutBob)
	roomWithBob := People{}
	err := roomWithBob.AddToAnotherRoomByName("bob", samplePeople)
	if err != nil {
		t.Error("got an error back:", err)
	}
	log.Println("With bob:", roomWithBob)

}
