package room_allocation

import (
	"errors"
)

const MaxUint = ^uint(0)
const MinInt = 0
const MaxInt = int(MaxUint >> 1)

var EndListError = errors.New("End of list reached")

func FindMinimum(valueGetter func(location int) (int, error)) (int, int) {
	minFound := MaxInt
	minLoc := -1
	for i := 0; i < MaxInt; i++ {
		v, err := valueGetter(i)
		if err != nil {
			return minFound, minLoc
		}
		if v <= minFound {
			minFound = v
			minLoc = i
		}
	}
	return minFound, minLoc
}
func (p People) MinConnection() (score int, location int) {
	getter := func(location int) (int, error) {
		if location >= len(p) {
			return MaxInt, EndListError
		}
		return int(p[location].MinConnections()), nil
	}
	return FindMinimum(getter)
}
func (p People) MinConnectionScore() Score {
	score, _ := p.MinConnection()
	return Score(score)
}

// MinConnectionPerson
// Within a People go through and return the person who has the lowest score
// (to anyone)
func (p People) MinConnectionPerson() *Person {
	_, location := p.MinConnection()
	return p[location]
}
func (p People) MinNConnections(n int) People {
	retArray := People{}
	for i := 0; i < n; i++ {
		personToAdd := p.MinConnectionPerson()
		retArray = append(retArray, personToAdd)
	}
	return retArray
}
func (p People) GetMinScoreWith(r Person) Score {
	getter := func(location int) (int, error) {
		if location >= len(p) {
			return MaxInt, EndListError
		}
		connection, err := p[location].GetConnection(r)
		if err != nil {
			return MaxInt, err
		}
		return int(connection.Count), nil
	}
	score, _ := FindMinimum(getter)
	return Score(score)
}

func (p Person) MinConnections() Score {
	getter := func(location int) (int, error) {
		if location >= len(p.Connections) {
			return MaxInt, EndListError
		}
		return int(p.Connections[location].Count), nil
	}
	score, _ := FindMinimum(getter)
	return Score(score)
}
func (s Scoreboard) MinValue() Score {
	getter := func(location int) (int, error) {
		if location >= len(s) {
			return MaxInt, EndListError
		}
		return int(s[location]), nil
	}
	score, _ := FindMinimum(getter)
	return Score(score)
}

func maxScResult(pa []People) People {
	maxLenFound := len(pa[0])
	maxEntry := MinInt
	for i, p := range pa {
		if len(p) > maxLenFound {
			maxEntry = i
			maxLenFound = len(p)
		}
	}
	return pa[maxEntry]
}
