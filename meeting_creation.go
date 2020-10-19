package room_allocation

import (
	"fmt"
	"strconv"
)

var NotEnoughPeopleError = fmt.Errorf("There are not enough people to split up into that many rooms")
var NotEnoughRoomsError = fmt.Errorf("That os not enough rooms to make sense")
var NotEnoughMeets = fmt.Errorf("Not enough meetings requested, need at least 1")

type MeetingSet []People
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

func roundDownIntDivide(a, b int) int {
	return (a + b - 1) / b
}
func (p People) targetNumberOfPeoplePerRoom(numRooms int, roomsAllocated int) int {
	initial := roundDownIntDivide(len(p), numRooms)
	return roundDownIntDivide(len(p)-(initial*roomsAllocated), numRooms-roomsAllocated)
}

func (p People) SplitIntoNRooms(n int) (meetingRooms MeetingSet, err error) {
	if len(p) < n {
		return nil, NotEnoughPeopleError
	}
	if len(p) <= 2 {
		return nil, NotEnoughPeopleError
	}
	if n < 2 {
		return nil, NotEnoughRoomsError
	}
	meetingRooms = make([]People, n)
	for i := 0; i < n; i++ {
		meetingRooms[i] = People{}
	}
	remainingPool := p.Copy()

	// Now let's get into the business logic!
	//fmt.Println("Targetting ", targetNumberPeoplePerRoom, "from", len(p))
	for i := 0; i < n-1; i++ {
		targetNumberPeoplePerRoom := p.targetNumberOfPeoplePerRoom(n, i)
		if len(remainingPool) < 2 {
			fmt.Println("We'll have a long person in the meeting room, bacause:", targetNumberPeoplePerRoom, len(p), p, n)
		}
		// Populate each room with the best people to go in them
		err := meetingRooms[i].AddBestNPeople(&remainingPool, targetNumberPeoplePerRoom)
		if err != nil {
			return nil, err
		}
	}
	// reuse remainingPool as the final meeting room
	meetingRooms[n-1] = remainingPool.EveryoneHereHasMet()
	return meetingRooms, nil
}

func (p People) AutoMeet(maxNumRooms, numberOfMeets int) (meetingRoomSeq MeetingSchedule, err error) {
	if numberOfMeets < 1 {
		return nil, NotEnoughMeets
	}
	if maxNumRooms > (len(p) / 2) {
		maxNumRooms = len(p) / 2
	}
	MaxItterations := 100
	for i := 0; (p.MinConnectionScore() < Score(numberOfMeets)) && (i < MaxItterations); i++ {
		//fmt.Println("AutoMeet:", i, p.MinConnectionScore(), numberOfMeets, p.ListConnections())
		mRs, err := p.SplitIntoNRooms(maxNumRooms)
		if err != nil {
			return nil, err
		}
		//fmt.Println("Have some meetings:", mRs)
		meetingRoomSeq = append(meetingRoomSeq, mRs)
	}
	return
}

type OptFunc func(People, MeetingSchedule) []int

func (p People) OptimalMeet(maxNumRooms, numberOfMeets, itterations int) (MeetingSchedule, error) {
	optFunc := func(po People, ml MeetingSchedule) []int {
		// Optimise first to have the fewest number of meetings to go to
		numberOfMeetingsNeeded := len(ml)
		// Second to have the minimum score
		// i.e. being in meetings with the same person several times, ideally you are less peaky!
		maxScore := int(po.MaxScore())
		return []int{numberOfMeetingsNeeded, maxScore}
	}
	return p.meetOptimiser(maxNumRooms, numberOfMeets, itterations, optFunc)
}
func (p People) RunMeetings(meetingSchedule MeetingSchedule) {
	for _, session := range meetingSchedule {
		for _, meeting := range session {
			p.RunMeeting(meeting)
		}
	}
}
func (p People) meetOptimiser(maxNumRooms, numberOfMeets, itterations int, optFunc OptFunc) (MeetingSchedule, error) {
	var meetingRoomSeq MeetingSchedule
	minVal := []int{MaxInt, MaxInt} // FIXME
	minimiser := func(tv []int) bool {
		for j := 0; j < len(tv); j++ {
			if tv[j] < minVal[j] {
				return true
			}
		}
		return false
	}

	for i := 0; i < itterations; i++ {
		pc := p.CopyBlank()
		meetingRoomSeqTemp, err := pc.AutoMeet(maxNumRooms, numberOfMeets)
		if err != nil {
			fmt.Println("Something went wrong!", err)
			return nil, err
		}
		if len(meetingRoomSeqTemp) == 0 {
			fmt.Println("Something went wrong with the length!")
			continue
		}
		tv := optFunc(pc, meetingRoomSeqTemp)
		if minimiser(tv) {
			minVal = tv
			meetingRoomSeq = meetingRoomSeqTemp
		}
	}
	p.RunMeetings(meetingRoomSeq)
	return meetingRoomSeq, nil
}
