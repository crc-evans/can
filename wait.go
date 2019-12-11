package can

import (
	"time"
)

// var ErrTimeout = fmt.Errorf("Timeout")

type ErrTimeout struct{}

func (err *ErrTimeout) Error() string { return "timeout" }

// A WaitResponse encapsulates the response of waiting for a frame.
type WaitResponse struct {
	Frame Frame
	Err   error
}

type waiter struct {
	id     uint32
	wait   chan WaitResponse
	bus    *Bus
	filter Handler
}

// Wait returns a channel, which receives a frame or an error, if the
// frame with the expected id didn't arrive on time.
func Wait(bus *Bus, id uint32, timeout time.Duration) <-chan WaitResponse {
	waiter := waiter{
		id:   id,
		wait: make(chan WaitResponse),
		bus:  bus,
	}

	ch := make(chan WaitResponse)

	go func() {
		select {
		case resp := <-waiter.wait:
			ch <- resp
		case <-time.After(timeout):
			err := &ErrTimeout{}
			ch <- WaitResponse{Frame{}, err}
		}
	}()

	waiter.filter = newFilter(id, &waiter)
	bus.Subscribe(waiter.filter)

	return ch
}

// WaitFunc returns a channel, which receives a frame or an error, if the
// frame with the expected id didn't arrive on time.
func WaitFunc(bus *Bus, filter func(Frame) bool, timeout time.Duration) <-chan WaitResponse {
	waiter := waiter{
		wait: make(chan WaitResponse),
		bus:  bus,
	}

	ch := make(chan WaitResponse)

	waiter.filter = newFuncFilter(filter, &waiter)

	go func() {
		select {
		case resp := <-waiter.wait:
			ch <- resp
		case <-time.After(timeout):
			bus.Unsubscribe(waiter.filter) // must unsubscribe so handler does not match future messages
			err := &ErrTimeout{}
			ch <- WaitResponse{Frame{}, err}
		}
	}()

	bus.Subscribe(waiter.filter)

	return ch
}

func (w *waiter) Handle(frame Frame) {
	w.bus.Unsubscribe(w.filter)
	w.wait <- WaitResponse{frame, nil}
}
