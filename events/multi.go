package Events

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

type listenData struct {
	index int
	param interface{}
}

// CreateMultipleEvent : create multiple event channel
func CreateMultipleEvent(chSize int) []chan interface{} {
	chSlice := make([]chan interface{}, chSize)
	for i := 0; i < chSize; i++ {
		chSlice[i] = make(chan interface{}, 1)
	}
	return chSlice
}

// CloseMultipleEvent : close multiple event channel
func CloseMultipleEvent(slChannel []chan interface{}) {
	defer func() { recover() }()

	once.Do(func() {
		for _, v := range slChannel {
			close(v)
			v = nil
		}
		slChannel = nil
	})
}

// WaitForMultipleObject
func WaitForMultipleObject(slChannel []chan interface{}, timeout uint64, bWaitAll bool) (idx int, retv []interface{}) {
	if !bWaitAll {
		return waitSingle(slChannel, timeout)
	} else {
		return waitAll(slChannel, timeout)
	}
}

// bWaitAll = false
func waitSingle(slChannel []chan interface{}, timeout uint64) (idx int, retv []interface{}) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			idx = -2
			retv = make([]interface{}, 0)
			retv = append(retv, errors.New(fmt.Sprint(r)))
		}
	}()

	// prevent timeout overflow
	if timeout == 0 {
		timeout = 0x00FFFFFFFFFFFFFF
	}

	// create local channels
	eventChannel := make(chan interface{}, len(slChannel))
	interrupt := make(chan bool, len(slChannel))
	defer close(eventChannel)
	defer close(interrupt)

	// add waitgroup
	wg.Add(len(slChannel))

	// listen channel threads
	for i, c := range slChannel {
		go func(i int, c chan interface{}, interrupt chan bool) {
			//defer func() { recover() }()
			select {
			case v := <-c:
				eventChannel <- listenData{param: v, index: i}
				break

			case <-interrupt:
				break
			}
			close(c)
			wg.Done()
		}(i, c, interrupt)
	}

	// wait event from threads
	select {
	case agg := <-eventChannel:
		idx = agg.(listenData).index
		retv = make([]interface{}, len(slChannel))
		retv[idx] = agg.(listenData).param
		break

	case <-time.After(time.Millisecond * time.Duration(timeout)):
		idx = -1
		retv = make([]interface{}, 0)
		retv = append(retv, errors.New(EVENT_TIMEOUT))
		break
	}

	terminateThread(len(slChannel), interrupt)
	wg.Wait()
	return idx, retv
}

// bWaitAll = true
func waitAll(slChannel []chan interface{}, timeout uint64) (idx int, retv []interface{}) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			idx = -2
			retv = make([]interface{}, 0)
			retv = append(retv, errors.New(fmt.Sprint(r)))
		}
	}()

	// prevent timeout overflow
	if timeout == 0 {
		timeout = 0x00FFFFFFFFFFFFFF
	}

	// init retv
	retv = make([]interface{}, len(slChannel))

	// create default return values
	idx = len(slChannel)
	result := make([]interface{}, len(slChannel))

	// create local channels
	eventChannel := make(chan interface{})
	interrupt := make(chan bool, len(slChannel))
	defer close(eventChannel)
	defer close(interrupt)

	// add waitgroup
	wg.Add(len(slChannel))

	// threads
	for i, c := range slChannel {
		go func(i int, c chan interface{}) {
			select {
			case v := <-c:
				eventChannel <- listenData{param: v, index: i}
				result[i] = v
				break

			case <-interrupt:
				break
			}
			close(c)
			wg.Done()
		}(i, c)
	}

	// wait event wait
	nCount := 0
_EVENT_WAIT:
	select {
	case v := <-eventChannel:
		result[v.(listenData).index] = v.(listenData).param
		nCount = nCount + 1
		if nCount < len(slChannel) {
			goto _EVENT_WAIT
		} else {
			retv = result
			break
		}

	case <-time.After(time.Millisecond * time.Duration(timeout)):
		idx = -1
		retv = make([]interface{}, 0)
		retv = append(retv, EVENT_TIMEOUT)
		break
	}

	terminateThread(len(slChannel), interrupt)
	wg.Wait()
	return idx, retv
}

// terminateThread : close listen threads
func terminateThread(cnt int, interrupt chan bool) {
	defer func() { recover() }()

	for i := 0; i < cnt; i++ {
		interrupt <- true
	}
}
