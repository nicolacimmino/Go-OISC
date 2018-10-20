package Components

import (
	"fmt"
	"sync"
)

type AddressLinesType uint8

type DataLinesType uint8

type BusSubscriber chan bool

/**
 * Map of named instances of this multiton.
 */
var instances = make(map[string]*Bus)

/**
 *
 */
type Bus struct {
	Address          AddressLinesType
	Data             DataLinesType
	Clock            chan bool
	Reset            chan bool
	Brk              chan bool
	RW               bool
	clockSubscribers []BusSubscriber
	resetSubscribers []BusSubscriber
	name             string
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

	bus := Bus{0, 0,
		make(chan bool),
		make(chan bool),
		make(chan bool),
		false,
		make([]BusSubscriber, 0),
		make([]BusSubscriber, 0),
		busName,
		sync.Mutex{},
	}

	instances[busName] = &bus

	go bus.processClockLine()
	go bus.processResetLine()

	return &bus
}

/**
 * Subscribe to clock bus events.
 */
func (bus *Bus) SubscribeToClock(subscriber BusSubscriber) {
	bus.clockSubscribers = append(bus.clockSubscribers, subscriber)
}

/**
 * Subscribe to reset bus events.
 */
func (bus *Bus) SubscribeToReset(subscriber BusSubscriber) {
	bus.resetSubscribers = append(bus.resetSubscribers, subscriber)
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
	close(bus.Reset)
}

/**
 * Clock devices on the bus.
 */
func (bus *Bus) clockBus(state bool) {
	bus.Lock()
	bus.Clock <- state
}

/**
 * Reset devices on the bus.
 */
func (bus *Bus) ResetBus() {
	bus.Lock()
	bus.Reset <- true
	bus.Lock()
	bus.Reset <- false
	bus.Lock()
	bus.Unlock()
}

/**
 * Monitor the bus Clock line and fan it out to all clockSubscribers.
 */
func (bus *Bus) processClockLine() {
	for {
		clock, ok := <-bus.Clock
		if !ok {
			return
		}

		for _, subscriber := range bus.clockSubscribers {
			subscriber <- clock
		}
		bus.Unlock()
	}
}

/**
 * Monitor the bus Reset line and fan it out to all resetSubscribers.
 */
func (bus *Bus) processResetLine() {
	for {
		reset, ok := <-bus.Reset
		if !ok {
			return
		}

		for _, subscriber := range bus.resetSubscribers {
			subscriber <- reset
		}
		bus.Unlock()
	}
}

/**
 * String representation of the bus status.
 */
func (bus *Bus) String() string {
	return fmt.Sprintf("A:%#X D:%#X", bus.Address, bus.Data)
}
