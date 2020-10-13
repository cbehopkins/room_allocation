package room_allocation

import "fmt"

type Score int

// The index into a scoreboard is the person's number
type Scoreboard []Score

type People []*Person

var PersonNotFoundError = fmt.Errorf("Person Not found")

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
	return p
}
func (p People) Names() []string {
	retStr := make([]string, len(p))
	for i, q := range p {
		retStr[i] = q.Name
	}
	return retStr
}
func (p People) Copy() People {
	ra := make(People, len(p))
	copy(ra, p)
	return ra
}
func (p People) CopyBlank() People {
	return NewPeople(p.Names())
}
func (p People) MaxScore() Score {
	maxFound := Score(0)
	for _, s := range p {
		m := s.MaxScore()
		if m > maxFound {
			maxFound = m
		}
	}
	return maxFound
}

func (p People) ListConnections() string {
	peopleStr := ""
	for _, r := range p {
		peopleStr += "\n" + r.ListConnections()
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

func (p People) generateMinimumsScoreboard(externalPeople People) Scoreboard {
	externalPeopleMinScoreboard := make(Scoreboard, len(p))
	for i, person := range p {
		externalPeopleMinScoreboard[i] = externalPeople.GetMinScoreWith(*person)
	}
	return externalPeopleMinScoreboard
}

func (p People) GetPeopleWithScore(q Person, s Score) (retArray People) {
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
func (p People) EveryoneHereHasMet() {
	for i, _ := range p {
		for j, m := range p {
			if j > i {
				m.AddToConnection(*p[i])
			}
		}
	}
}
func (p People) RunMeeting(q People) {
	room := p.GetPeopleByName(q.Names())
	room.EveryoneHereHasMet()
	//fmt.Println("RunMeeting: People", p, p.ListConnections(), "after", q)
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
