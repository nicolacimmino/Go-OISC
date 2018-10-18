package Memory

import "Go-OISC/Bus"

/**
 *
 */
type Memory struct {
	data                []Bus.DataLinesType
	baseAddress         Bus.AddressLinesType
	size                Bus.AddressLinesType
	bus                 *Bus.Bus
	busEventsSubscriber Bus.Subscriber
}

/**
 * Constructor.
 */
func NewMemory(size Bus.AddressLinesType, baseAddress Bus.AddressLinesType, bus *Bus.Bus) *Memory {
	busReceiver := make(Bus.Subscriber)
	bus.Subscribe(busReceiver)

	m := Memory{make([]Bus.DataLinesType, size), baseAddress, size, bus, busReceiver}

	go m.process()

	return &m
}

/**
 * Release resources.
 */
func (memory *Memory) Close() {
	close(memory.busEventsSubscriber)
}

/**
 * Wait for events from the bus and process memory commands.
 */
func (memory *Memory) process() {
	for {
		clock, ok := <-memory.busEventsSubscriber
		if !ok {
			return
		}

		if clock {
			if memory.bus.Address < memory.baseAddress || memory.bus.Address > memory.baseAddress+memory.size {
				continue
			}

			if memory.bus.RW {
				memory.bus.Data = memory.data[memory.bus.Address-memory.baseAddress]
				continue
			}

			memory.data[memory.bus.Address-memory.baseAddress] = memory.bus.Data
		}
	}
}
