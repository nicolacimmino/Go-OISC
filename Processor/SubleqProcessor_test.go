package Processor

import (
	"Go-OISC/Bus"
	"Go-OISC/Memory"
	"testing"
	"time"
)

func TestRunAnd(t *testing.T) {
	bus := Bus.NewBus("main")
	Memory.NewMemory(0xFF, 0x00, bus)
	NewSubleqProcessor(bus)

	testVectors := [][]Bus.DataLinesType{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
		{1, 1, 1},
	}

	for _, testVector := range testVectors {
		A := testVector[0]
		B := testVector[1]

		for address, instruction := range []Bus.DataLinesType{
			0x00, A, 0xFF,
			0x00, B, 0xFF,
			0x00, 0x01, 0xFF,
		} {
			bus.Write(Bus.AddressLinesType(address), instruction)
		}

		bus.ResetBus()
		for {
			if bus.Halt {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		if bus.Read(0x06) != testVector[2] {
			t.Fail()
		}
	}
}
