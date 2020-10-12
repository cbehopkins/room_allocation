package room_allocation

import "testing"

func TestSelectOptimumOverlap0(t *testing.T) {
	testcases := []struct {
		p                    People
		numberPeopleToSelect int
	}{
		{NewPeople([]string{"a", "b", "c", "d"}), 2},
		{NewPeople([]string{"a", "b", "c", "d", "e", "f"}), 3},
		{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 2},
		{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 3},
		{NewPeople([]string{"a", "b", "c", "d", "e", "f", "g", "h"}), 4},
	}
	for _, tc := range testcases {
		tSelectOptimumOverlap(tc.p, t, tc.numberPeopleToSelect)
	}
}

func tSelectOptimumOverlap(samplePeople People, t *testing.T, numberPeopleToSelect int) {
	remainingPool := samplePeople.Copy()
	meetingRoom := People{}
	meetingRoom.AddBestNPeople(&remainingPool, numberPeopleToSelect)
	t.Log("MR:", meetingRoom, "Remaining", remainingPool)
	t.Log("*********Resetting for round 2!")
	remainingPool = samplePeople.Copy()
	meetingRoom = People{}
	meetingRoom.AddBestNPeople(&remainingPool, numberPeopleToSelect-1)
	t.Log("MR:", meetingRoom, "Remaining", remainingPool)
	t.Log(samplePeople.ListConnections())

	res := meetingRoom.SelectOptimumOverlap(remainingPool)
	t.Log("For MR:", meetingRoom, "Optimimum is apparently", res)
	for _, person0 := range res {
		for _, person1 := range meetingRoom {
			connection, err := person0.GetConnection(*person1)
			if err != nil {
				t.Error("Unable to find connection between ", person0, " and ", person1)
			}
			if connection.Count != 0 {
				t.Error("The connection cound should be 0, is:", connection.Count)
			}
		}
	}
}
