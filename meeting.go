package room_allocation

import "strconv"

type Meeting struct {People}

type MeetingSet []Meeting
type MeetingSchedule []MeetingSet


func (m MeetingSet) String() string {
	retStr := ""
	spacerT := "    "
	spacer := spacerT
	for i, set := range m {
		retStr += spacer + "    " + strconv.Itoa(i) + ":" + set.String()
		spacer = "\n" + spacerT
	}
	return retStr
}

func (m MeetingSchedule) String() string {
	retStr := ""
	for i, schedule := range m {
		retStr += "Session:" + strconv.Itoa(i) + ":\n" + schedule.String() + "\n"

	}
	return retStr + "\n"
}
func (m MeetingSet) Len() int {
	return len(m)
}
func (m MeetingSchedule) Len () int {
	return len(m)
}
func roundDownIntDivide(a, b int) int {
	return (a + b - 1) / b
}
func (p Meeting) targetNumberOfPeoplePerRoom(numRooms int, roomsAllocated int) int {
	initial := roundDownIntDivide(p.Len(), numRooms)
	return roundDownIntDivide(p.Len()-(initial*roomsAllocated), numRooms-roomsAllocated)
}
