package components

import (
	"testing"
)

/**
 * Ensure we can read back a value written in memory.
 */
func TestReadBack(t *testing.T) {
	bus := NewBus("main")
	defer bus.Close()

	memory := NewMemory(0xFF, 0x0, bus)
	defer memory.Close()

	bus.Write(0x05, 0x02)
	bus.Write(0x06, 0x10)

	if bus.Read(0x05) != 0x02 {
		t.Fail()
	}
}

/**
 * Ensure reading outside the address space of the memory returns zero.
 */
func TestReadOutsideBounds(t *testing.T) {
	bus := NewBus("main")
	defer bus.Close()

	memory := NewMemory(0x20, 0x0, bus)
	defer memory.Close()

	bus.Write(0x05, 0x02)

	if bus.Read(0x15) != 0x00 {
		t.Log(bus)
		t.Fail()
	}
}

/**
 * Ensure memories on different memory blocks don't interfere.
 */
func TestMultipleMemories(t *testing.T) {
	bus := NewBus("main")
	defer bus.Close()

	memoryA := NewMemory(0x20, 0x0, bus)
	defer memoryA.Close()

	memoryB := NewMemory(0x20, 0x20, bus)
	defer memoryB.Close()

	bus.Write(0x05, 0x02)
	bus.Write(0x25, 0x22)

	if bus.Read(0x05) != 0x02 {
		t.Log(bus)
		t.Fail()
	}

	if bus.Read(0x25) != 0x22 {
		t.Log(bus)
		t.Fail()
	}
}
