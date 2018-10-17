package Bus

import (
	"testing"
)

func TestSubscribe(t *testing.T) {
	bus := NewBus()
	receiver := make(chan busCycle)
	bus.Register(receiver)

	bus.Address = 0x10
	bus.Data = 0xFF
	bus.Clock <- true

	busOperation := <-receiver

	if busOperation.Address != 0x10 || busOperation.Data != 0xFF {
		t.Fail()
	}
}

func TestSubscribeMultiple(t *testing.T) {
	bus := NewBus()
	receiverA := make(chan busCycle)
	bus.Register(receiverA)

	receiverB := make(chan busCycle)
	bus.Register(receiverB)

	listener := func(receiver chan busCycle, aOK *bool) {
		for {
			cycle := <-receiver
			if cycle.Clock == true {
				*aOK = true
			}
		}
	}

	aOK := false
	go listener(receiverA, &aOK)
	bOK := false
	go listener(receiverB, &bOK)

	bus.Address = 0x10
	bus.Data = 0xFF
	bus.Clock <- true
	bus.Clock <- false

	if !aOK || !bOK {
		t.Fail()
	}
}
