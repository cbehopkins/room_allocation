package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cbehopkins/room_allocation"
)

func main() {
	roomCntPtr := flag.Int("rooms", 2, "Number of rooms to use")
	optCntPtr := flag.Int("opt", 8, "How many things to try (2^opt)")
	meetCnt := flag.Int("meets", 1, "What's the minimum number of times people should meet")
	flag.Parse()
	peopleList := flag.Args()
	if len(peopleList) < 4 {
		fmt.Println("Please supply more people, need 4, got:", peopleList)
		os.Exit(1)
	}
	peeps := room_allocation.NewPeople(peopleList)
	roomsSchedule, err := peeps.ToMeeting().OptimalMeet(*roomCntPtr, *meetCnt, 1<<*optCntPtr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}
	fmt.Println(roomsSchedule)
	fmt.Println(peeps.ListConnections())
}
