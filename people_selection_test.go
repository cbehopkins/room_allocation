package room_allocation

import "testing"

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
