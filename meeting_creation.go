package room_allocation

import (
	"fmt"
	"runtime"
	"sync"
)

var NotEnoughPeopleError = fmt.Errorf("There are not enough people to split up into that many rooms")
var NotEnoughRoomsError = fmt.Errorf("That os not enough rooms to make sense")
var NotEnoughMeets = fmt.Errorf("Not enough meetings requested, need at least 1")

func (p People) ToMeeting() Meeting {
	return Meeting{p}
}

func (p Meeting) SplitIntoNRooms(n int) (meetingRooms MeetingSet, err error) {
	if p.Len() < n {
		return nil, NotEnoughPeopleError
	}
	if p.Len() <= 2 {
		return nil, NotEnoughPeopleError
	}
	if n < 2 {
		return nil, NotEnoughRoomsError
	}
	meetingRooms = make(MeetingSet, n)
	for i := 0; i < n; i++ {
		meetingRooms[i] = Meeting{}
	}
	remainingPool := p.Copy()

	// Now let's get into the business logic!
	//fmt.Println("Targetting ", targetNumberPeoplePerRoom, "from", len(p))
	for i := 0; i < n-1; i++ {
		targetNumberPeoplePerRoom := p.targetNumberOfPeoplePerRoom(n, i)
		if remainingPool.Len() < 2 {
			fmt.Println("We'll have a long person in the meeting room, bacause:", targetNumberPeoplePerRoom, p.Len(), p, n)
		}
		// Populate each room with the best people to go in them
		err := meetingRooms[i].AddBestNPeople(&remainingPool, targetNumberPeoplePerRoom)
		if err != nil {
			return nil, err
		}
	}
	// reuse remainingPool as the final meeting room
	meetingRooms[n-1] = remainingPool.EveryoneHereHasMet().ToMeeting()
	return meetingRooms, nil
}

func (p Meeting) AutoMeet(maxNumRooms, numberOfMeets int) (meetingRoomSeq MeetingSchedule, err error) {
	if numberOfMeets < 1 {
		return nil, NotEnoughMeets
	}
	if maxNumRooms > (p.Len() / 2) {
		maxNumRooms = p.Len() / 2
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

type OptFunc func(Meeting, MeetingSchedule) []int

func (p Meeting) OptimalMeet(maxNumRooms, numberOfMeets, itterations int) (MeetingSchedule, error) {
	optFunc := func(po Meeting, ml MeetingSchedule) []int {
		// Optimise first to have the fewest number of meetings to go to
		numberOfMeetingsNeeded := ml.Len()
		// Second to have the minimum score
		// i.e. being in meetings with the same person several times, ideally you are less peaky!
		maxScore := int(po.MaxScore())
		return []int{numberOfMeetingsNeeded, maxScore}
	}
	return p.meetOptimiser(maxNumRooms, numberOfMeets, itterations, optFunc)
}
func (p Meeting) RunMeetings(meetingSchedule MeetingSchedule) {
	for _, session := range meetingSchedule {
		for _, meeting := range session {
			p.RunMeeting(meeting)
		}
	}
}

type resultStruct struct {
	p   Meeting
	m   MeetingSchedule
	err error
}

func (p Meeting) optimiserWorker(maxNumRooms, numberOfMeets int, wg *sync.WaitGroup, sourceChan <-chan Meeting, resultChan chan<- resultStruct) {
	for p := range sourceChan {
		pc := p.CopyBlank().ToMeeting()
		meetingRoomSeqTemp, err := pc.AutoMeet(maxNumRooms, numberOfMeets)

		if err != nil {
			resultChan <- resultStruct{pc, nil, err}
			continue
		}

		resultChan <- resultStruct{pc, meetingRoomSeqTemp, nil}
	}
	wg.Done()
}
func (p Meeting) consolidateOptimisation(optFunc OptFunc, resultChan <-chan resultStruct) <-chan MeetingSchedule {
	meetingScheduleChan := make(chan MeetingSchedule)
	go func() {
		defer close(meetingScheduleChan)
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
		for results := range resultChan {
			if results.err != nil {
				meetingScheduleChan <- nil
				return
			}
			tv := optFunc(results.p, results.m)
			if minimiser(tv) {
				minVal = tv
				meetingRoomSeq = results.m
			}
		}
		meetingScheduleChan <- meetingRoomSeq

	}()
	return meetingScheduleChan
}
func (p Meeting) meetOptimiser(maxNumRooms, numberOfMeets, itterations int, optFunc OptFunc) (MeetingSchedule, error) {
	sourceChan := make(chan Meeting)
	resultChan := make(chan resultStruct)

	var wg sync.WaitGroup
	workerCount := runtime.NumCPU()
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go p.optimiserWorker(maxNumRooms, numberOfMeets, &wg, sourceChan, resultChan)
	}
	go func() {
		for i := 0; i < itterations; i++ {
			sourceChan <- p
		}
		close(sourceChan)
	}()
	meetingScheduleChan := p.consolidateOptimisation(optFunc, resultChan)
	wg.Wait()
	close(resultChan)
	meetingRoomSeq, _ := <-meetingScheduleChan
	p.RunMeetings(meetingRoomSeq)
	return meetingRoomSeq, nil
}
