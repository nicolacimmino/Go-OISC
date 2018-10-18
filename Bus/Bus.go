package Bus

import (
	"sync"
)

type AddressLinesType uint8
type DataLinesType uint8

type Bus struct {
	Address     AddressLinesType
	Data        DataLinesType
	Clock       chan bool
	RW          bool
	subscribers []chan BusCycle
	name        string
	Halt        bool
	sync.RWMutex
}

type BusCycle struct {
	Address AddressLinesType
	Data    DataLinesType
	RW      bool
	Clock   bool
	Bus     *Bus
}

var instances = make(map[string]*Bus)

func NewBus(busName string) *Bus {

	existingBus, exixsts := instances[busName]

	if exixsts {
		return existingBus
	}

	bus := Bus{0, 0, make(chan bool), false, make([]chan BusCycle, 0), busName, true, sync.RWMutex{}}

	go bus.process()

	instances[busName] = &bus

	return &bus
}

func (bus *Bus) Close() {
	delete(instances, bus.name)
	close(bus.Clock)
}

func (bus *Bus) Register(subscriber chan BusCycle) {
	bus.subscribers = append(bus.subscribers, subscriber)
}

func (bus *Bus) Write(address AddressLinesType, data DataLinesType) {
	bus.Data = data
	bus.Address = address
	bus.RW = false
	bus.Lock()
	bus.Clock <- true
	bus.Lock()
	bus.Clock <- false

	bus.Lock()
	defer bus.Unlock()
}

func (bus *Bus) Read(address AddressLinesType) DataLinesType {
	//bus.Lock()
	bus.Address = address
	bus.RW = true
	bus.Lock()
	bus.Clock <- true
	bus.Lock()
	bus.Clock <- false

	bus.Lock()
	defer bus.Unlock()

	return bus.Data
}

func (bus *Bus) process() {
	for {
		clock, ok := <-bus.Clock
		if !ok {
			return
		}

		cycle := BusCycle{bus.Address, bus.Data, bus.RW, clock, bus}

		for _, subscriber := range bus.subscribers {
			subscriber <- cycle
		}
		bus.Unlock()
	}
}
