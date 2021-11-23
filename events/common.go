package Events

import (
	"errors"
	"fmt"
	"sync"
)

var once sync.Once
var THREAD_TIMEOUT string = "EVENT.TIMEOUT"

// set event
func SetEvent(ch chan interface{}, v interface{}) (err error) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	// send channel
	//  TODO: How to check ch is closed?
	ch <- v
	return err
}
