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

	// # Create Single Event
	//chEvent := Events.CreateSingleEvent()
	// # Test-Case #1
	//_testInvoke(chEvent)

	// # Create Multi Event
	var nEvent int = 3
	slEvent := Events.CreateMultipleEvent(nEvent)
	go Thread_MultiWait(slEvent)
	time.Sleep(time.Second)
	Events.SetEvent(slEvent[0], "event#1")
	Events.SetEvent(slEvent[1], "event#2")
	Events.SetEvent(slEvent[2], "event#3")

	for !isFin {
		// # Wait a second
		time.Sleep(time.Second)
	}

	// # Close event
	//Events.CloseSingleEvent(chEvent)
	Events.CloseMultipleEvent(slEvent)
	fmt.Println("[#Main] Fin Process")
}

/* Test-Case #1 */
func _testInvoke(chEvent chan interface{}) {
	// # Invoke Thread
	go Thread_SingleWait(chEvent)

	// # Wait a 3-second
	time.Sleep(time.Second * 3)

	// # Invoke event with value: if comment this line, check timeout
	Events.SetEvent(chEvent, "invoke event from #main")
}

func Thread_SingleWait(chEvent chan interface{}) {
	fmt.Println("[#Thread #1] Start")

	var ui64MilliSec uint64 = 5000

	// Wait event for five seconds
	retv := Events.WaitForSingleObject(chEvent, ui64MilliSec)
	fmt.Println("[#Thread #1] WaitForSingleObject.result:", retv)

	isFin = true
	fmt.Println("[#Thread #1] Fin")
}

func Thread_MultiWait(slEvent []chan interface{}) {
	fmt.Println("[#Thread #1] Start")

	var ui64MilliSec uint64 = 5000

	// Wait event for five seconds
	idx, retv := Events.WaitForMultipleObject(slEvent, ui64MilliSec, false)
	fmt.Println("[#Thread #1] WaitForMultiple.idx:", idx, ", data:", retv)
	fmt.Println("[#Thread #1] :", retv[0])
	fmt.Println("[#Thread #1] :", retv[1])
	fmt.Println("[#Thread #1] :", retv[2])

	isFin = true
	fmt.Println("[#Thread #1] Fin")
}
