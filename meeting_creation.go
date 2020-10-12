package room_allocation

import "fmt"

var NotEnoughPeopleError = fmt.Errorf("There are not enough people to split up into that many rooms")
var NotEnoughRoomsError = fmt.Errorf("That os not enough rooms to make sense")
var NotEnoughMeets = fmt.Errorf("Not enough meetings requested, need at least 1")

func (p People) SplitIntoNRooms(n int) (meetingRooms []People, err error) {
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
	targetNumberPeoplePerRoom := (len(p) + n - 1) / n
	//fmt.Println("Targetting ", targetNumberPeoplePerRoom, "from", len(p))
	for i := 0; i < n-1; i++ {
		err := meetingRooms[i].AddBestNPeople(&remainingPool, targetNumberPeoplePerRoom)
		if err != nil {
			return nil, err
		}
	}
	// reuse remainingPool as the final meeting room
	remainingPool.EveryoneHereHasMet()
	meetingRooms[n-1] = remainingPool
	return meetingRooms, nil
}

func (p People) AutoMeet(maxNumRooms, numberOfMeets int) (meetingRoomSeq [][]People, err error) {
	if numberOfMeets < 1 {
		return nil, NotEnoughMeets
	}
	if maxNumRooms > (len(p) / 2) {
		maxNumRooms = len(p) / 2
	}
	for i := 0; p.MinConnectionScore() < Score(numberOfMeets); i++ {
		mRs, err := p.SplitIntoNRooms(maxNumRooms)
		if err != nil {
			return nil, err
		}
		meetingRoomSeq = append(meetingRoomSeq, mRs)
	}
	return
}
