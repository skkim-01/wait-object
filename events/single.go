package Events

import (
	"errors"
	"fmt"
	"time"
)

// CreateSigleEvent : create single event channel
func CreateSingleEvent() chan interface{} {
	ch := make(chan interface{}, 1)
	return ch
}

// CloseSingleEvent : close channel
func CloseSingleEvent(ch chan interface{}) {
	once.Do(func() {
		close(ch)
		ch = nil
	})
}

// WaitForSingleObject
func WaitForSingleObject(ch chan interface{}, timeout uint64) (retv interface{}) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			retv = errors.New(fmt.Sprint(r))
		}
	}()

	// prevent timeout overflow
	if timeout == 0 {
		timeout = 0x00FFFFFFFFFFFFFF
	}

	select {
	case v := <-ch:
		return v

	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return errors.New(THREAD_TIMEOUT)
	}
}
