package Bus

import (
	"testing"
)

func TestSubscribe(t *testing.T) {
	bus := NewBus("main")
	defer bus.Close()

	receiver := make(chan BusCycle)
	bus.Register(receiver)

	bus.Lock()
	bus.Address = 0x10
	bus.Data = 0xFF
	bus.Clock <- true

	busOperation := <-receiver

	if busOperation.Address != 0x10 || busOperation.Data != 0xFF {
		t.Fail()
	}
}

func TestSubscribeMultiple(t *testing.T) {
	bus := NewBus("main")
	defer bus.Close()

	receiverA := make(chan BusCycle)
	bus.Register(receiverA)

	receiverB := make(chan BusCycle)
	bus.Register(receiverB)

	listener := func(receiver chan BusCycle, aOK *bool) {
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

	bus.Write(0x10, 0xFF)

	if !aOK || !bOK {
		t.Log(aOK, bOK)
		t.Fail()
	}
}
