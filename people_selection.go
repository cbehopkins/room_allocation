package room_allocation

import "fmt"

var TooLongError = fmt.Errorf("Returned room too long")
var NoneSuitableFoundError = fmt.Errorf("No suitable people found")

// SelectOptimumOverlap
// Return the People who have the minimum connection score with external People
// I.e. The people who have had the least connections
// between the two groups
func (p People) SelectOptimumOverlap(externalPeople People) People {
	fmt.Println("SelectOptimumOverlap get in::", p, " and ", externalPeople)
	// For each person in the external group
	// What's the minimum connection score with internal group
	minimumScore := p.generateMinimumsScoreboard(externalPeople)
	//	for i, person := range p {
	//		fmt.Println("Person:", person, " has score ", minimumScore[i])
	//	}
	resultScoreboard := make([]People, len(p))
	for i, m := range p {
		// Collect a list of people who have this minimum score
		// The scoreboard's index is cogent with the people in p
		// And the list is the people who have that score
		resultScoreboard[i] = externalPeople.GetPeopleWithScore(*m, minimumScore.MinValue())
	}
	return p.determineScoreboardOverlap(resultScoreboard, minimumScore.MinValue())
}
func (p People) determineScoreboardOverlap(resultScoreboard []People, minScore Score) People {
	correlateMap := make(map[string]int)
	lookupMap := make(map[string]*Person)
	for i, person := range p {
		if false {
			fmt.Println("Person:", person, " has score ", minScore, " with ", resultScoreboard[i])
		}
		for _, person1 := range resultScoreboard[i] {
			if val, ok := correlateMap[person1.Name]; ok {
				correlateMap[person1.Name] = val + 1
			} else {
				correlateMap[person1.Name] = 1
				lookupMap[person1.Name] = person1
			}
		}
	}
	return p.returnLowestPeople(correlateMap, lookupMap)
}

func (p People) returnLowestPeople(correlateMap map[string]int, lookupMap map[string]*Person) (retP People) {
	lowestVal := 0
	for _, val := range correlateMap {
		if val > lowestVal {
			lowestVal = val
		}
	}
	for key, val := range correlateMap {
		if val == lowestVal {
			retP = append(retP, lookupMap[key])
		}
	}
	return
}

// AddBestNPeople
// Add in people to this room from the source room
// and delete them from the source room
// Select the "best" people for this
func (p *People) AddBestNPeople(sourceRoom *People, n int) error {
	targetLen := len(*p) + n
	for len(*p) < targetLen {
		//fmt.Println("p is", p, "n", n)
		err, numAdded := p.addUpToNBestPeople(sourceRoom, n)
		if err != nil {
			return err
		}
		//fmt.Println("Num Added", numAdded)
		n -= numAdded
	}
	// fmt.Println("p is now", p, "n", n)
	if len(*p) > targetLen {
		return TooLongError
	}
	return nil
}

func (p *People) addUpToNBestPeople(sourceRoom *People, n int) (error, int) {
	adj := 0
	if len(*p) == 0 {
		minimumPerson := sourceRoom.MinConnectionPerson()
		p.AddPersonToMeeting(minimumPerson)
		sourceRoom.RemovePerson(minimumPerson)
		n -= 1
		adj = 1
		//fmt.Println("Selected minimum Person", minimumPerson)
	}
	overlapGroup := p.SelectOptimumOverlap(*sourceRoom)
	if len(overlapGroup) == 0 {
		return NoneSuitableFoundError, 0
	}
	if len(overlapGroup) < n {
		n = len(overlapGroup)
	}
	p.AddPeopleToMeeting(overlapGroup[:n])
	return sourceRoom.RemovePeople(overlapGroup[:n]), n + adj
}
