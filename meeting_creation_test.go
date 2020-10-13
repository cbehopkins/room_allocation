package room_allocation

import (
	"log"
	"testing"
)

func TestMeeting0(t *testing.T) {
	// First of all declare some people
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa", "Steve", "James"})
	samplePeople0 := samplePeople.CopyBlank()
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
	remainingPool := samplePeople.CopyBlank()
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

type MeetingTestData struct {
	p                 People
	numberRooms       int
	targetConnections int
	minConnections    int
	maxConnections    int
}

func TestMeetingAutoMeet(t *testing.T) {
	// First of all declare some people
	testData := []MeetingTestData{
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d"}), 2, 1, 1, 2},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d"}), 2, 2, 2, 3},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 2, 1, 1, 5},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 3, 1, 1, 4},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 3, 2, 2, 5},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 4, 1, 1, 2},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h",
			"i", "j", "k", "l", "m", "n", "p", "q"}), 3, 1, 1, 9},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h",
			"i", "j", "k", "l", "m", "n", "p", "q"}), 4, 1, 1, 5},
	}
	for i, td := range testData {
		tMeetingAutoMeet(td, i, t)
	}
}

func tMeetingAutoMeet(td MeetingTestData, cnt int, t *testing.T) {
	samplePeople := td.p
	numberRooms := td.numberRooms
	targetConnections := td.targetConnections
	minConnections := Score(td.minConnections)
	maxConnections := Score(td.maxConnections)
	roomsSchedule, err := samplePeople.AutoMeet(numberRooms, targetConnections)
	if err != nil {
		t.Error(cnt, err)
	}

	for i, rooms := range roomsSchedule {
		t.Log("Session ", i, "Rooms:", rooms)
	}
	t.Log(samplePeople.ListConnections())
	for _, person0 := range samplePeople {
		for _, connection := range person0.Connections {
			if connection.Count < minConnections {
				t.Error(cnt, "We should have minimum", minConnections, " connection between ", person0, "and ", connection.PerLink, " have:", connection.Count)
			}
			if connection.Count > maxConnections {
				t.Error(cnt, "We should have maximum", maxConnections, " connection between ", person0, "and ", connection.PerLink, " have:", connection.Count)
			}
		}
	}
}

func TestMeetingOptimal(t *testing.T) {
	// First of all declare some people
	testData := []MeetingTestData{
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d"}), 2, 1, 1, 2},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d"}), 2, 2, 2, 2},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 2, 1, 1, 4},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 3, 1, 1, 3},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 3, 2, 2, 7},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 4, 1, 1, 3},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h",
			"i", "j", "k", "l", "m", "n", "p", "q"}), 3, 1, 1, 8},
		MeetingTestData{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h",
			"i", "j", "k", "l", "m", "n", "p", "q"}), 4, 1, 1, 6},
	}
	for i, td := range testData {
		tMeetingOptimal(td, i, t)
	}
}

func tMeetingOptimal(td MeetingTestData, cnt int, t *testing.T) {
	samplePeople := td.p
	numberRooms := td.numberRooms
	targetConnections := td.targetConnections
	minConnections := Score(td.minConnections)
	maxConnections := Score(td.maxConnections)
	roomsSchedule, err := samplePeople.OptimalMeet(numberRooms, targetConnections, 200)
	if err != nil {
		t.Error(cnt, err)
	}
	if roomsSchedule == nil {
		t.Error("No room schedule retrurned")
	}
	t.Log("Results for tsetcase:", cnt)
	for i, rooms := range roomsSchedule {
		t.Log("Session ", i, "Rooms:", rooms)
	}
	t.Log(samplePeople.ListConnections())
	for _, person0 := range samplePeople {
		for _, connection := range person0.Connections {
			if connection.Count < minConnections {
				t.Error(cnt, "We should have minimum", minConnections, " connection between ", person0, "and ", connection.PerLink, " have:", connection.Count)
			}
			if connection.Count > maxConnections {
				t.Error(cnt, "We should have maximum", maxConnections, " connection between ", person0, "and ", connection.PerLink, " have:", connection.Count)
			}
		}
	}
}
