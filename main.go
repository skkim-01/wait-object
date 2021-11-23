package main

import (
	"fmt"
	"time"

	Events "github.com/skkim-01/wait-object/events"
	//Event "./event"
)

var isFin bool = false

func main() {
	fmt.Println("[#Main] Start Test Main")

	// # Create Event
	chEvent := Events.CreateSingleEvent()

	// # Test-Case #1
	_testInvoke(chEvent)

	for !isFin {
		// # Wait a second
		time.Sleep(time.Second)
	}

	// # Close event
	Events.CloseSingleEvent(chEvent)
	fmt.Println("[#Main] Fin Process")
}

/* Test-Case #1 */
func _testInvoke(chEvent chan interface{}) {
	// # Invoke Thread
	go Thread_1st(chEvent)

	// # Wait a 3-second
	time.Sleep(time.Second * 3)

	// # Invoke event with value: if comment this line, check timeout
	Events.SetEvent(chEvent, "invoke event from #main")
}
func Thread_1st(chEvent chan interface{}) {
	fmt.Println("[#Thread #1] Start")

	var ui64MilliSec uint64 = 5000

	// Wait event for five seconds
	retv := Events.WaitForSingleObject(chEvent, ui64MilliSec)
	fmt.Println("[#Thread #1] WaitForSingleObject.result:", retv)

	isFin = true
	fmt.Println("[#Thread #1] Fin")
}
