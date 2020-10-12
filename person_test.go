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
	err := bob.addToConnection(*fred)
	if err != nil {
		t.Error("Got an error trying to form connection")
	}
	err = fred.addToConnection(*bob)

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
		t.Log("Have a look:", samplePeople.ListConnections())
	}

	t.Log("Now let's check the reverse")
	// The reverse connection should also have been populated
	connection, err = fred.GetConnection(*bob)
	if err != nil {
		t.Error("Couldn't get fred's connection to bob")
	}
	if connection.Count != 1 {
		t.Error("This connection should have a score of 1, Got:", connection)
	}
	connectionScore = samplePeople.MinConnectionScore()
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
	minimumPerson := samplePeople.MinConnectionPerson()
	if minimumPerson != lisa {
		t.Error("We should have selected Lisa:", minimumPerson)
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
	roomWithoutBob := samplePeople.NewRoomWithoutName("bob")
	log.Println("Room without bob:", roomWithoutBob)
	roomWithBob := People{}
	err := roomWithBob.AddToAnotherRoomByName("bob", samplePeople)
	if err != nil {
		t.Error("got an error back:", err)
	}
	log.Println("With bob:", roomWithBob)

}

func TestSelectCandidate0(t *testing.T) {
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa"})

	var bob *Person

	_, bob = samplePeople.GetPersonByName("bob")
	_, fred := samplePeople.GetPersonByName("fred")
	_, lisa := samplePeople.GetPersonByName("Lisa")
	bob.AddToConnection(*fred)
	bob.AddToConnection(*lisa)
	log.Println(bob.ListConnections())

	connection_bob_fred, err := bob.GetConnection(*fred)
	if err != nil {
		t.Error("Couldn't get bob's connection to fred")
	}
	if connection_bob_fred.Count != 1 {
		t.Error("This connection should have a score of 1, Got:", connection_bob_fred)
	}

	candidates := samplePeople.GetPeopleWithScore(*bob, 0)
	if len(candidates) != 0 {
		t.Error("No-one should have a score of 0, we got:", candidates)
	}
	candidates = samplePeople.GetPeopleWithScore(*bob, 1)
	if len(candidates) != 2 {
		t.Error("There should be 2 candidates, we got:", candidates)
	}
}

func TestSelectCandidate1(t *testing.T) {
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa"})

	var bob *Person

	_, bob = samplePeople.GetPersonByName("bob")
	_, fred := samplePeople.GetPersonByName("fred")
	bob.AddToConnection(*fred)
	log.Println(bob.ListConnections())

	connection_bob_fred, err := bob.GetConnection(*fred)
	if err != nil {
		t.Error("Couldn't get bob's connection to fred")
	}
	if connection_bob_fred.Count != 1 {
		t.Error("This connection should have a score of 1, Got:", connection_bob_fred)
	}

	candidates := samplePeople.GetPeopleWithScore(*bob, 0)
	if len(candidates) != 1 {
		t.Error("Lisa should have a score of 0, we got:", candidates)
	}
	candidates = samplePeople.GetPeopleWithScore(*bob, 1)
	if len(candidates) != 1 {
		t.Error("There should be 1 candidates, we got:", candidates)
	}
}

func TestSelectGroup0(t *testing.T) {
	samplePeople0 := NewPeople([]string{"bob", "fred", "Lisa", "Steve"})
	samplePeople1 := samplePeople0.NewRoomWithoutNames([]string{"bob", "fred"})
	if len(samplePeople1) != 2 {
		t.Error("We should have extracted 2 people, instead we got:", samplePeople1)
	}

	_, bob := samplePeople0.GetPersonByName("bob")
	_, fred := samplePeople0.GetPersonByName("fred")
	_, lisa := samplePeople0.GetPersonByName("Lisa")
	_, steve := samplePeople0.GetPersonByName("Steve")

	bob.AddToConnection(*steve)

	minScore := samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 0 {
		t.Error("There should be some people who aren't connected to everyone 0")
	}
	bob.AddToConnection(*fred)
	bob.AddToConnection(*lisa)
	minScore = samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 0 {
		t.Error("There should be some people who aren't connected to everyone 1")
	}
	fred.AddToConnection(*lisa)
	fred.AddToConnection(*steve)
	lisa.AddToConnection(*steve)

	minScore = samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 1 {
		t.Error("Everyone should now be connected to everyone")
	}
}

func TestSelectGroup1(t *testing.T) {
	samplePeople0 := NewPeople([]string{"bob", "fred", "Lisa", "Steve"})
	samplePeople1 := samplePeople0.NewRoomWithoutNames([]string{"bob", "fred"})
	if len(samplePeople1) != 2 {
		t.Error("We should have extracted 2 people, instead we got:", samplePeople1)
	}

	_, bob := samplePeople0.GetPersonByName("bob")
	_, fred := samplePeople0.GetPersonByName("fred")
	_, lisa := samplePeople0.GetPersonByName("Lisa")
	_, steve := samplePeople0.GetPersonByName("Steve")

	bob.AddToConnection(*steve)

	minScore := samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 0 {
		t.Error("There should be some people who aren't connected to everyone 0")
	}
	bob.AddToConnection(*fred)
	bob.AddToConnection(*fred)
	bob.AddToConnection(*fred)
	bob.AddToConnection(*lisa)
	fred.AddToConnection(*lisa)
	fred.AddToConnection(*lisa)
	minScore = samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 0 {
		t.Error("There should be some people who aren't connected to everyone 1")
	}
	fred.AddToConnection(*lisa)
	fred.AddToConnection(*steve)
	lisa.AddToConnection(*steve)

	minScore = samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 1 {
		t.Error("Everyone should now be connected to everyone")
	}
}

func TestSelectGroup2(t *testing.T) {
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa", "Steve"})

	samplePeople0 := samplePeople.GetPeopleByName([]string{"bob", "fred"})
	if len(samplePeople0) != 2 {
		t.Error("We should have extracted 2 people, instead we got:", samplePeople0)
	}
	// Note this selects those without
	samplePeople1 := samplePeople.NewRoomWithoutNames([]string{"bob", "fred"})
	if len(samplePeople1) != 2 {
		t.Error("We should have left 2 people, instead we got:", samplePeople1)
	}

	_, bob := samplePeople.GetPersonByName("bob")
	_, fred := samplePeople.GetPersonByName("fred")
	_, lisa := samplePeople.GetPersonByName("Lisa")
	_, steve := samplePeople.GetPersonByName("Steve")

	bob.AddToConnection(*fred)
	bob.AddToConnection(*lisa)
	fred.AddToConnection(*lisa)
	fred.AddToConnection(*steve)
	lisa.AddToConnection(*steve)

	minScore := samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 0 {
		t.Error("There should be some people who aren't connected to everyone 0")
	}

	// Steve and Bob have not met, will the tool tell us this?
	people2 := samplePeople0.SelectOptimumOverlap(samplePeople1)
	if len(people2) != 1 {
		t.Error("Damn, there's not 1 person in:", people2)
	}
	log.Println("We should select any of:", people2)

	bob.AddToConnection(*steve)

	minScore = samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 1 {
		t.Error("Everyone should now be connected to everyone")
	}
}

func TestSelectGroup3(t *testing.T) {
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa", "Steve"})

	samplePeople0 := samplePeople.GetPeopleByName([]string{"bob", "fred"})
	if len(samplePeople0) != 2 {
		t.Error("We should have extracted 2 people, instead we got:", samplePeople0)
	}
	samplePeople1 := samplePeople.NewRoomWithoutNames([]string{"bob", "fred"})
	if len(samplePeople1) != 2 {
		t.Error("We should have extracted 2 people, instead we got:", samplePeople1)
	}

	_, bob := samplePeople.GetPersonByName("bob")
	_, fred := samplePeople.GetPersonByName("fred")
	_, lisa := samplePeople.GetPersonByName("Lisa")
	_, steve := samplePeople.GetPersonByName("Steve")

	bob.AddToConnection(*fred)
	bob.AddToConnection(*lisa)
	fred.AddToConnection(*lisa)
	fred.AddToConnection(*steve)
	lisa.AddToConnection(*steve)

	// Steve and Bob have not met, will the tool tell us this?
	people2 := samplePeople0.SelectOptimumOverlap(samplePeople1)
	if len(people2) != 1 {
		t.Error("Damn, there's not 1 person in:", people2)
	}
	t.Log("Steve's connection before adding to the meeting are:", steve.ListConnections())
	samplePeople0.AddPersonToMeeting(people2[0])
	t.Log("Steve's connection after adding to the meeting are:", steve.ListConnections())
	log.Println("The room should now have Steve in:", samplePeople0.ListConnections())
	minScore := samplePeople0.generateMinimumsScoreboard(samplePeople1)
	if minScore.MinValue() != 1 {
		t.Error("Everyone should now be connected to everyone")
	}
}
func TestMinimumPeople0(t *testing.T) {
	samplePeople := NewPeople([]string{"bob", "fred", "Lisa", "Steve"})
	minimumPerson := samplePeople.MinConnectionPerson()
	t.Log("We have selected:", minimumPerson)
	_, bob := samplePeople.GetPersonByName("bob")
	_, fred := samplePeople.GetPersonByName("fred")
	_, lisa := samplePeople.GetPersonByName("Lisa")
	_, steve := samplePeople.GetPersonByName("Steve")
	bob.AddToConnection(*fred)
	bob.AddToConnection(*lisa)
	bob.AddToConnection(*steve)
	fred.AddToConnection(*steve)
	lisa.AddToConnection(*steve)
	minimumPerson = samplePeople.MinConnectionPerson()
	if minimumPerson != lisa {
		t.Error("We should have selected Lisa:", minimumPerson)
	}
	t.Log("We selected Lisa:", minimumPerson)

	// Now make sure everyone has met with everyone and it doesn't fail
	lisa.AddToConnection(*fred)
	minimumPerson = samplePeople.MinConnectionPerson()
	t.Log("We selected:", minimumPerson)

	bob.AddToConnection(*fred)
	bob.AddToConnection(*lisa)
	bob.AddToConnection(*steve)
	fred.AddToConnection(*steve)
	lisa.AddToConnection(*steve)
	minimumPerson = samplePeople.MinConnectionPerson()
	if minimumPerson != lisa {
		t.Error("We should have selected Lisa:", minimumPerson)
	}
	t.Log("We selected Lisa:", minimumPerson)

}
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
func TestSelectOptimumOverlap0(t *testing.T) {
	samplePeople := NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"})
	remainingPool := samplePeople.Copy()
	meetingRoom := People{}
	meetingRoom.AddBestNPeople(&remainingPool, 4)
	t.Log("MR:", meetingRoom, "Remaining", remainingPool)
	t.Log("*********Resetting for round 2!")
	remainingPool = samplePeople.Copy()
	meetingRoom = People{}
	meetingRoom.AddBestNPeople(&remainingPool, 3)
	t.Log("MR:", meetingRoom, "Remaining", remainingPool)
	t.Log(samplePeople.ListConnections())

	res := meetingRoom.SelectOptimumOverlap(remainingPool)
	t.Log("Optimimum is apparently", res)
}
