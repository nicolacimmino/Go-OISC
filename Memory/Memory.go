package Memory

import (
	"Go-OISC/Bus"
)

type Memory struct {
	data        []Bus.DataLinesType
	baseAddress Bus.AddressLinesType
	size        Bus.AddressLinesType
	bus         *Bus.Bus
	busReceiver chan Bus.BusCycle
}

func NewMemory(size Bus.AddressLinesType, baseAddress Bus.AddressLinesType, bus *Bus.Bus) *Memory {
	busReceiver := make(chan Bus.BusCycle)
	bus.Register(busReceiver)

	m := Memory{make([]Bus.DataLinesType, size), baseAddress, size, bus, busReceiver}

	go m.process()

	return &m
}

func (Memory *Memory) Close() {
	close(Memory.busReceiver)
}

func (memory *Memory) process() {
	for {
		busCycle, ok := <-memory.busReceiver
		if !ok {
			return
		}

		if busCycle.Clock {
			if busCycle.Address < memory.baseAddress || busCycle.Address > memory.baseAddress+memory.size {
				continue
			}

			if busCycle.RW {
				memory.bus.Data = memory.data[busCycle.Address-memory.baseAddress]
				continue
			}

			memory.data[busCycle.Address-memory.baseAddress] = busCycle.Data
		}
	}
}
