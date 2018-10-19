package Memory

import "Go-OISC/Bus"

/**
 *
 */
type Memory struct {
	data               []Bus.DataLinesType
	baseAddress        Bus.AddressLinesType
	size               Bus.AddressLinesType
	bus                *Bus.Bus
	busClockSubscriber Bus.Subscriber
}

/**
 * Constructor.
 */
func NewMemory(size Bus.AddressLinesType, baseAddress Bus.AddressLinesType, bus *Bus.Bus) *Memory {
	busClockSubscriber := make(Bus.Subscriber)
	bus.SubscribeToClock(busClockSubscriber)

	m := Memory{make([]Bus.DataLinesType, size), baseAddress, size, bus, busClockSubscriber}

	go m.process()

	return &m
}

/**
 * Release resources.
 */
func (memory *Memory) Close() {
	close(memory.busClockSubscriber)
}

/**
 * Wait for events from the bus and process memory commands.
 */
func (memory *Memory) process() {
	for {
		clock, ok := <-memory.busClockSubscriber
		if !ok {
			break
		}

		if clock {
			if memory.bus.Address < memory.baseAddress || memory.bus.Address >= memory.baseAddress+memory.size {
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
