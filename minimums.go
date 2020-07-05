package room_allocation

import "errors"

const MaxUint = ^uint(0)
const MinInt = 0
const MaxInt = int(MaxUint >> 1)

var EndListError = errors.New("End of list reached")

func FindMinimum(valueGetter func(location int) (int, error)) int {
	minFound := MaxInt
	for i := 0; i < MaxInt; i++ {
		v, err := valueGetter(i)
		if err != nil {
			return minFound
		}
		if v <= minFound {
			minFound = v
		}
	}
	return MaxInt
}

func (p People) MinConnections() Score {
	getter := func(location int) (int, error) {
		if location >= len(p) {
			return MaxInt, EndListError
		}
		return int(p[location].MinConnections()), nil
	}

	return Score(FindMinimum(getter))
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
	return Score(FindMinimum(getter))
}

func (p Person) MinConnections() Score {
	getter := func(location int) (int, error) {
		if location >= len(p.Connections) {
			return MaxInt, EndListError
		}
		return int(p.Connections[location].Count), nil
	}
	return Score(FindMinimum(getter))
}
func (s Scoreboard) MinValue() Score {
	getter := func(location int) (int, error) {
		if location >= len(s) {
			return MaxInt, EndListError
		}
		return int(s[location]), nil
	}
	return Score(FindMinimum(getter))

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
