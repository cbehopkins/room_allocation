package room_allocation

import "fmt"

type Score int

// The index into a scoreboard is the person's number
type Scoreboard []Score

type People []*Person

var PersonNotFoundError = fmt.Errorf("Person Not found")
var NoneSuitableFoundError = fmt.Errorf("No suitable people found")
var TooLongError = fmt.Errorf("Returned room too long")
var NotEnoughPeopleError = fmt.Errorf("There are not enough people to split up into that many rooms")
var NotEnoughRoomsError = fmt.Errorf("That os not enough rooms to make sense")
var NotEnoughMeets = fmt.Errorf("Not enough meetings requested, need at least 1")

func NewPeople(l []string) People {
	p := make(People, len(l))
	for i := range l {
		p[i] = NewPerson(l[i], len(l))
	}
	for i := range l {
		// Go to every person
		for j := range l {
			// And add a link to every other person
			if j != i {
				p[i].AddConnection(p[j])
			}
		}
	}
	fmt.Println("Created People:", p)
	return p
}
func (p People) Copy() People {
	ra := make(People, len(p))
	copy(ra, p)
	return ra
}
func (p People) ListConnections() string {
	peopleStr := ""
	spacer := ""
	for _, r := range p {
		peopleStr += spacer + r.ListConnections()
		spacer = ", "
	}
	return "[" + peopleStr + "]"
}
func (p People) GetPersonByName(name string) (int, *Person) {
	for i, m := range p {
		if m.Name == name {
			return i, m
		}
	}
	return -1, nil
}
func (p People) GetPeopleByName(names []string) People {
	np := make(People, 0, len(names))

	for _, name := range names {
		i, _ := p.GetPersonByName(name)
		np = append(np, p[i])
	}
	return np
}

// SelectOptimumOverlap
// Return the People who have the minimum connection score with external People
// I.e. The people who have had the least connections
// between the two groups
func (p People) SelectOptimumOverlap(externalPeople People) People {
	// For each person in the external group
	// What's the minimum connection score with internal group
	minimumScore := p.generateMinimumsScoreboard(externalPeople)

	resultScoreboard := make([]People, len(p))
	for i, m := range p {
		fmt.Print("Looking at:", m)
		// Collect a list of people whi have this minimum score
		resultScoreboard[i] = externalPeople.GetPeopleWithScore(*m, minimumScore.MinValue())
	}
	return maxScResult(resultScoreboard)
}

func (p People) generateMinimumsScoreboard(externalPeople People) Scoreboard {
	externalPeopleMinScoreboard := make(Scoreboard, len(p))
	for i, person := range p {
		externalPeopleMinScoreboard[i] = externalPeople.GetMinScoreWith(*person)
	}
	return externalPeopleMinScoreboard
}

// GetMinimumN
func (p People) GetMinimum() People {
	return p
}

func (p People) GetPeopleWithScore(q Person, s Score) People {
	var retArray People
	for _, m := range p {
		// Get one person from this group of people
		connection, err := m.GetConnection(q)
		if err != nil {
			// It's fine they're not in the list
			continue
		}
		// And get their score to the target person
		if connection.Count == s {
			retArray = append(retArray, m)
		}
	}
	return retArray
}
func (p People) NewRoomWithoutPerson(person Person) People {
	// TODO - this could be so more efficient
	return p.NewRoomWithoutName(person.Name)
}
func (p People) NewRoomWithoutName(name string) People {
	// Note we do not modify the origional room because we will go back to that
	// But basically we want to select everyone who isn't name
	// into a new slicve

	// Locate the position of this person in the string
	loc, _ := p.GetPersonByName(name)

	// And reslice
	return append(append(People{}, p[0:loc]...), p[loc+1:]...)
}
func (p People) NewRoomWithoutNames(names []string) People {
	l := len(names)
	if l == 0 {
		return nil
	}
	if l == 1 {
		return p.NewRoomWithoutName(names[0])
	}
	return p.NewRoomWithoutName(names[0]).NewRoomWithoutNames(names[1:])

}
func (p People) FindPersonIndex(r *Person) (int, error) {
	for i, m := range p {
		if m == r {
			return i, nil
		}
	}
	return -1, PersonNotFoundError
}

func (p *People) AddToAnotherRoomByName(name string, r People) error {
	// Add A person to this room, from another
	i, _ := r.GetPersonByName(name)
	if i < 0 {
		return PersonNotFoundError
	}
	*p = append(*p, r[i])
	return nil
}

func (p *People) AddPersonToMeeting(r *Person) {
	for _, m := range *p {
		// Give these people a connection number bump
		m.AddToConnection(*r)
	}
	*p = append(*p, r)
}
func (p *People) AddPeopleToMeeting(r People) {
	for _, s := range r {
		p.AddPersonToMeeting(s)
	}
}
func (p *People) RemovePerson(r *Person) error {
	indexToRemove, err := p.FindPersonIndex(r)
	if err != nil {
		return err
	}
	(*p)[indexToRemove] = (*p)[len(*p)-1]
	(*p)[len(*p)-1] = nil
	*p = (*p)[:len(*p)-1]
	return nil
}
func (p *People) RemovePeople(r People) error {
	for _, s := range r {
		err := p.RemovePerson(s)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *People) AddBestNPeople(sourceRoom *People, n int) error {
	targetLen := len(*p) + n
	for len(*p) < targetLen {
		err, numAdded := p.addUpToNBestPeople(sourceRoom, n)
		if err != nil {
			return err
		}
		n -= numAdded
	}
	if len(*p) > targetLen {
		return TooLongError
	}
	return nil
}

func (p *People) addUpToNBestPeople(sourceRoom *People, n int) (error, int) {
	if len(*p) == 0 {
		minimumPerson := sourceRoom.MinConnectionPerson()
		p.AddPersonToMeeting(minimumPerson)
		sourceRoom.RemovePerson(minimumPerson)
		n -= 1
	}
	fmt.Println("Meeting Room:", p, "Source Room:", sourceRoom)
	overlapGroup := p.SelectOptimumOverlap(*sourceRoom)
	if len(overlapGroup) == 0 {
		return NoneSuitableFoundError, 0
	}
	if len(overlapGroup) < n {
		n = len(overlapGroup)
	}
	p.AddPeopleToMeeting(overlapGroup[:n])
	return sourceRoom.RemovePeople(overlapGroup[:n]), n
}

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
	remainingPool := p.Copy() // Is this needed?

	// Now let's get into the business logic!
	targetNumberPeoplePerRoom := len(p) / n
	for i := 0; i < n; i++ {
		err := meetingRooms[i].AddBestNPeople(&remainingPool, targetNumberPeoplePerRoom)
		if err != nil {
			return nil, err
		}
	}
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
		mRs, err := p.SplitIntoNRooms(2)
		if err != nil {
			return nil, err
		}
		meetingRoomSeq = append(meetingRoomSeq, mRs)
	}
	return
}
