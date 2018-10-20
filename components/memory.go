package components

/**
 *
 */
type Memory struct {
	data               []DataLinesType
	baseAddress        AddressLinesType
	size               AddressLinesType
	bus                *Bus
	busClockSubscriber BusSubscriber
}

/**
 * Constructor.
 */
func NewMemory(size AddressLinesType, baseAddress AddressLinesType, bus *Bus) *Memory {
	busClockSubscriber := make(BusSubscriber)
	bus.SubscribeToClock(busClockSubscriber)

	m := Memory{make([]DataLinesType, size), baseAddress, size, bus, busClockSubscriber}

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
