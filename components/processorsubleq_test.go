package components

import (
	"testing"
)

func TestRunAnd(t *testing.T) {
	bus := NewBus("main")
	NewMemory(0xFF, 0x00, bus)
	NewSubleqProcessor(bus)

	testVectors := [][]DataLinesType{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
		{1, 1, 1},
	}

	for _, testVector := range testVectors {
		A := testVector[0]
		B := testVector[1]

		for address, instruction := range []DataLinesType{
			0x0F, 0x0C, 0xFF, // 0x00	SUBLEQ 0C 09 FF
			0x10, 0X0D, 0xFF, // 0x03	SUBLEQ 0D 09 FF
			0x0E, 0x11, 0xFF, // 0x06	SUBLEQ 0E 11 FF
			0x11, 0x11, 0xFF, // 0x09   SUBLEQ 11 11 FF
			A, // 0x0C	Operand A
			B, // 0x0D	Operand B
			0, // 0x0E	Result
			0, // 0x0F	Swap 1
			0, // 0x10	Swap 2
			1, // 0x11   Swap 3
		} {
			bus.Write(AddressLinesType(address), instruction)
		}

		bus.ResetBus()
		<-bus.Brk

		if bus.Read(0x0E) != testVector[2] {
			t.Fail()
		}
	}
}
