package room_allocation

import (
	"log"
	"testing"
)

func TestMeeting0(t *testing.T) {
	// First of all declare some people
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa", "Steve", "James"})
	samplePeople0 := samplePeople.Copy()
	// Select the a person with the lowest score
	minimumPerson := samplePeople0.MinConnectionPerson()
	meetingRoom := People{minimumPerson}
	remainingPool := samplePeople.NewRoomWithoutPerson(*minimumPerson)

	err := meetingRoom.AddBestNPeople(&remainingPool, 1)
	if err != nil {
		t.Error("Adding people failed:", err)
	}
	t.Log("We now have:", meetingRoom)

}
func TestMeeting1(t *testing.T) {
	// First of all declare some people
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa", "Steve", "James"})
	meetingRoom := People{}
	remainingPool := samplePeople.Copy()
	err := meetingRoom.AddBestNPeople(&remainingPool, 3)
	if err != nil {
		t.Error("Adding people failed:", err)
	}
	t.Log("We now have:", meetingRoom)
}
func TestMeeting2(t *testing.T) {
	// First of all declare some people
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa", "Steve", "James", "Lucy"})
	for i := 0; i < 3; i++ {
		mRs, err := samplePeople.SplitIntoNRooms(2)
		if err != nil {
			t.Error("Creating Meeting Rooms failed:", err)
		}
		log.Println("MRs:", mRs, "Min Connections:", samplePeople.MinConnectionScore())
	}
}
func TestMeeting3(t *testing.T) {
	// First of all declare some people
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa", "Steve"})
	roomsSchedule, err := samplePeople.AutoMeet(2, 1)
	if err != nil {
		t.Error(t)
	}

	for i, rooms := range roomsSchedule {
		t.Log("Session ", i, "Rooms:", rooms)
	}
}
func TestMeeting4(t *testing.T) {
	// First of all declare some people
	samplePeople := NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"})
	roomsSchedule, err := samplePeople.AutoMeet(2, 1)
	if err != nil {
		t.Error(t)
	}

	for i, rooms := range roomsSchedule {
		t.Log("Session ", i, "Rooms:", rooms)
	}
	t.Log(samplePeople.ListConnections())
}
func TestMeeting5(t *testing.T) {
	// First of all declare some people
	samplePeople := NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"})
	roomsSchedule, err := samplePeople.AutoMeet(4, 1)
	if err != nil {
		t.Error(t)
	}

	for i, rooms := range roomsSchedule {
		t.Log("Session ", i, "Rooms:", rooms)
	}
	t.Log(samplePeople.ListConnections())
}
func TestMeeting6(t *testing.T) {
	// First of all declare some people
	samplePeople := NewPeople([]string{"a", "b", "c", "d"})
	roomsSchedule, err := samplePeople.AutoMeet(2, 1)
	if err != nil {
		t.Error(t)
	}

	for i, rooms := range roomsSchedule {
		t.Log("Session ", i, "Rooms:", rooms)
	}
	t.Log(samplePeople.ListConnections())
}
