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
			if j > i {
				p[i].AddConnection(p[j])
			}
		}
	}
	return p
}

func (p People) GetPersonByName(name string) (int, *Person) {
	for i, m := range p {
		if m.Name == name {
			return i, m
		}
	}
	return -1, nil
}

// SelectOptimumOverlap
// Return the People who have the minimum connection score with external People
// I.e. The people who have had the least connections
// between the two groups
func (p People) SelectOptimumOverlap(externalPeople People) People {
	// For each person in the external group
	// What's the minimum connection score with internal group
	var externalPeopleMinScoreboard Scoreboard
	for i, m := range p {
		externalPeopleMinScoreboard[i] = externalPeople.GetMinScoreWith(*m)
	}
	minimumScore := externalPeopleMinScoreboard.MinValue()

	var resultScoreboard []People
	for i, m := range p {
		// Collect a list of people whi have this minimum score
		resultScoreboard[i] = externalPeople.GetPeopleWithScore(*m, minimumScore)
	}
	return maxScResult(resultScoreboard)
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
			// Yep, that's a person which has that score
			retArray = append(retArray, m)
		}
	}
	return retArray
}

func (p People) TakeOutOfRoomByName(name string) People {
	// Note we do not modify the origional room because we will go back to that
	// But basically we want to select everyone who isn't name
	// into a new slicve

	// Locate the position of this person in the string
	loc, _ := p.GetPersonByName(name)

	// And reslice
	return append(append(People{}, p[0:loc]...), p[loc+1:]...)
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
