package room_allocation

import (
	"log"
	"testing"
)

func TestPerson0(t *testing.T) {
	samplePeople := NewPeople([]string{"bob", "fred", "Jane", "Lisa"})
	log.Println("Min connections is:", samplePeople.MinConnections())
}
