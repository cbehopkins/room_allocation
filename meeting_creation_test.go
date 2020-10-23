package room_allocation

import (
	"fmt"
	"strconv"
	"testing"
)

func generatePeopleList(cnt int) People {
	retList := make([]string, cnt)
	if cnt > 52 {
		return nil
	}
	for i := 0; i < cnt; i++ {
		j := i
		if i > 25 {
			j -= 58
		}
		retList[i] = string('\u0041' + 32 + j)
	}
	return NewPeople(retList)
}

type MeetingTestData struct {
	p                 int
	numberRooms       int
	targetConnections int
	minConnections    int
	maxConnections    int
}

func TestMeetingAutoMeet(t *testing.T) {
	// First of all declare some people
	testData := []MeetingTestData{
		MeetingTestData{4, 2, 1, 1, 0},
		MeetingTestData{4, 2, 2, 2, 0},
		MeetingTestData{8, 2, 1, 1, 0},
		MeetingTestData{8, 3, 1, 1, 0},
		MeetingTestData{8, 3, 2, 2, 0},
		MeetingTestData{8, 4, 1, 1, 0},
		MeetingTestData{16, 3, 1, 1, 0},
		MeetingTestData{16, 4, 1, 1, 0},
	}
	for i, td := range testData {
		tMeetingAutoMeet(td, i, t)
	}
}

func tMeetingAutoMeet(td MeetingTestData, cnt int, t *testing.T) {
	samplePeople := generatePeopleList(td.p)
	numberRooms := td.numberRooms
	targetConnections := td.targetConnections
	minConnections := Score(td.minConnections)
	maxConnections := Score(td.maxConnections)
	roomsSchedule, err := samplePeople.ToMeeting().AutoMeet(numberRooms, targetConnections)
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
			if (maxConnections > 0) && (connection.Count > maxConnections) {
				t.Error(cnt, "We should have maximum", maxConnections, " connection between ", person0, "and ", connection.PerLink, " have:", connection.Count)
			}
		}
	}
}

func TestMeetingOptimal(t *testing.T) {
	// First of all declare some people
	testData := []MeetingTestData{
		// {	p, numberRooms, targetConnections, minConnections, maxConnections}
		MeetingTestData{4, 2, 1, 1, 2},
		MeetingTestData{4, 2, 2, 2, 2},
		MeetingTestData{8, 2, 1, 1, 4},
		MeetingTestData{8, 3, 1, 1, 3},
		MeetingTestData{8, 3, 2, 2, 3},
		MeetingTestData{8, 4, 1, 1, 3},
		MeetingTestData{16, 2, 1, 1, 7},
		MeetingTestData{16, 3, 1, 1, 7},
		MeetingTestData{16, 4, 1, 1, 4}, //8
		MeetingTestData{16, 5, 1, 1, 4},
		MeetingTestData{16, 6, 1, 1, 3},
		MeetingTestData{16, 7, 1, 1, 5},
		MeetingTestData{16, 8, 1, 1, 2},
	}
	for i, td := range testData {
		tMeetingOptimal(td, strconv.Itoa(i), t)
	}
}
func TestMeetingOptimalGen(t *testing.T) {
	resultArray := make([]string, 0, 64)
	for numPeople := 4; numPeople < 33; numPeople++ {
		for numberRooms := 2; numberRooms < 8; numberRooms++ {
			if numberRooms > (numPeople / 2) {
				continue
			}
			for targetConnections := 1; targetConnections < 2; targetConnections++ {
				minConnections := targetConnections
				maxConnections := numPeople - 1
				td := MeetingTestData{numPeople, numberRooms, targetConnections, minConnections, maxConnections}
				lenRooms, maxCon := tMeetingOptimal(td, fmt.Sprint("numPeople:", numPeople, " numberRooms:", numberRooms, " targetConnections:", targetConnections), t)
				resultArray = append(resultArray, fmt.Sprint("numPeople:", numPeople, " numberRooms:", numberRooms, " NumSessions:", lenRooms, " MaxRemeet:", maxCon))
			}
		}
	}
	for _, s := range resultArray {
		t.Log(s)
	}
}
func tMeetingOptimal(td MeetingTestData, cnt string, t *testing.T) (int, int) {
	samplePeople := generatePeopleList(td.p)
	numberRooms := td.numberRooms
	targetConnections := td.targetConnections
	minConnections := Score(td.minConnections)
	maxConnections := Score(td.maxConnections)
	roomsSchedule, err := samplePeople.ToMeeting().OptimalMeet(numberRooms, targetConnections, 200)
	if err != nil {
		t.Error(cnt, err)
	}
	if roomsSchedule == nil {
		t.Error("No room schedule retrurned")
	}

	maxConn := 0
	for _, person0 := range samplePeople {
		for _, connection := range person0.Connections {
			if connection.Count < minConnections {
				t.Error(cnt, "We should have minimum", minConnections, " connection between ", person0, "and ", connection.PerLink, " have:", connection.Count)
				return 0, 0

			}
			if connection.Count > maxConnections {
				t.Error(cnt, "We should have maximum", maxConnections, " connection between ", person0, "and ", connection.PerLink, " have:", connection.Count)
				return 0, 0
			}
			if int(connection.Count) > maxConn {
				maxConn = int(connection.Count)
			}
		}
	}
	return len(roomsSchedule), maxConn
}
