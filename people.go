package room_allocation

type People []*Person

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

type Score int

// The index into a scoreboard is the person's number
type Scoreboard []Score

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
