package Components

import "testing"

/**
 * Ensure the bus is notifying a single subscriber.
 */
func TestSubscribe(t *testing.T) {
	bus := NewBus("main")
	defer bus.Close()

	subscriber := make(BusSubscriber)
	bus.SubscribeToClock(subscriber)

	bus.Lock()
	bus.Address = 0x10
	bus.Data = 0xFF
	bus.Clock <- true

	_ = <-subscriber

	bus.Lock()
	defer bus.Unlock()

	if bus.Address != 0x10 || bus.Data != 0xFF {
		t.Fail()
	}
}

/**
 * Ensure the bus is fanning out events to multiple clockSubscribers.
 */
func TestSubscribeMultiple(t *testing.T) {
	const subscribersCount = 10

	bus := NewBus("main")
	defer bus.Close()

	listener := func(subscriber BusSubscriber, oks *[subscribersCount]bool, ix int, bus Bus) {
		for {
			clock := <-subscriber
			if clock {
				oks[ix] = true
			}
		}
	}

	subscribers := [10]BusSubscriber{}
	oks := [subscribersCount]bool{}

	for ix := 0; ix < subscribersCount; ix++ {
		subscribers[ix] = make(BusSubscriber)
		bus.SubscribeToClock(subscribers[ix])
		go listener(subscribers[ix], &oks, ix, *bus)
	}

	bus.Write(0x10, 0xFF)

	for ix := 0; ix < subscribersCount; ix++ {
		if !oks[ix] {
			t.Fail()
		}
	}
}
