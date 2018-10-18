package Bus

import (
	"fmt"
	"sync"
)

type AddressLinesType uint8

type DataLinesType uint8

type Subscriber chan bool

/**
 * Map of named instances of this multiton.
 */
var instances = make(map[string]*Bus)

/**
 *
 */
type Bus struct {
	Address     AddressLinesType
	Data        DataLinesType
	Clock       chan bool
	RW          bool
	subscribers []Subscriber
	name        string
	Halt        bool
	sync.Mutex
}

/**
 * Constructor. Bus is a multiton, named instances can be created/retrieved with NewBus.
 */
func NewBus(busName string) *Bus {

	existingBus, exists := instances[busName]

	if exists {
		return existingBus
	}

	bus := Bus{0, 0, make(chan bool), false, make([]Subscriber, 0), busName, true, sync.Mutex{}}
	instances[busName] = &bus

	go bus.process()

	return &bus
}

/**
 * Subscribe to bus events.
 */
func (bus *Bus) Subscribe(subscriber Subscriber) {
	bus.subscribers = append(bus.subscribers, subscriber)
}

/**
 * Perform a write operation.
 */
func (bus *Bus) Write(address AddressLinesType, data DataLinesType) {
	bus.Data = data
	bus.Address = address
	bus.RW = false
	bus.clockBus(true)
	bus.clockBus(false)

	bus.Lock()
	bus.Unlock()
}

/**
 * Perform a read operation.
 */
func (bus *Bus) Read(address AddressLinesType) DataLinesType {
	bus.Address = address
	bus.RW = true
	bus.clockBus(true)
	bus.clockBus(false)

	bus.Lock()
	defer bus.Unlock()

	return bus.Data
}

/**
 * Close the bus and release resources.
 */
func (bus *Bus) Close() {
	delete(instances, bus.name)
	close(bus.Clock)
}

/**
 * Clock the bus.
 */
func (bus *Bus) clockBus(state bool) {
	bus.Lock()
	bus.Clock <- state
}

/**
 * Monitor the bus Clock and fan it out to all subscribers.
 */
func (bus *Bus) process() {
	for {
		clock, ok := <-bus.Clock
		if !ok {
			return
		}

		for _, subscriber := range bus.subscribers {
			subscriber <- clock
		}
		bus.Unlock()
	}
}

func (bus *Bus) String() string {
	return fmt.Sprintf("A:%#X D:%#X", bus.Address, bus.Data)
}
