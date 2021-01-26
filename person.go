package room_allocation

import (
	"encoding/json"
	"errors"
	"strconv"
)

var PersonDoesntExistError = errors.New("The person does not exist")
var PersonNotConnectedError = errors.New("The person is not found as a connection")

type Connection struct {
	Count   Score
	PerLink *Person
}

func (c Connection) Is(p Person) bool {
	return c.PerLink.Is(p)
}

type Person struct {
	Name        string
	Connections []Connection
}

func (c Connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Count   Score
		PerLink string
	}{
		c.Count,
		c.PerLink.Name})
}

type PersonNumber int

func NewPerson(name string, l int) *Person {
	p := new(Person)
	p.Name = name
	p.Connections = make([]Connection, 0, l)
	return p
}
func (p Person) MaxScore() Score {
	maxScore := Score(0)
	for _, c := range p.Connections {
		if c.Count > maxScore {
			maxScore = c.Count
		}
	}
	return maxScore
}

func (p Person) ListConnections() string {
	retStr := ""
	sep := ""
	for _, c := range p.Connections {
		retStr += sep + c.PerLink.String() + ":" + strconv.Itoa(int(c.Count))
		sep = ", "
	}
	return p.Name + " is connected to:[" + retStr + "]"
}

func (p Person) String() string {
	return p.Name
}

func (p *Person) AddConnection(r *Person) {
	p.Connections = append(p.Connections, Connection{Count: 0, PerLink: r})
}

func (p Person) Is(r Person) bool {
	return p.Name == r.Name
}

// Return the connection to the specified Person
func (p Person) getConnection(r Person) (*Connection, error) {
	// Try the connection in the first direction
	for i, m := range p.Connections {
		if m.PerLink.Is(r) {
			return &p.Connections[i], nil
		}
	}

	return nil, PersonNotConnectedError
}
func (p Person) GetConnection(r Person) (*Connection, error) {
	return r.getConnection(p)
}

// And increment the value of the connection betwen two people
func (p Person) addToConnection(r Person) error {
	connection, err := p.getConnection(r)
	if err != nil || connection == nil {
		return err
	}
	connection.Count++
	return nil
}
func (p Person) AddToConnection(r Person) error {
	err := p.addToConnection(r)
	if err != nil {
		return err
	}
	return r.addToConnection(p)
}
