package Memory

import (
	"Go-OISC/Bus"
	"testing"
)

func TestReadBack(t *testing.T) {
	bus := Bus.NewBus("main")
	defer bus.Close()

	memory := NewMemory(0xFF, 0x0, bus)
	defer memory.Close()

	bus.Write(0x05, 0x02)
	bus.Write(0x06, 0x10)
	bus.Read(0x05)

	if bus.Data != 0x02 {
		t.Fail()
	}
}

func TestReadOutsideBounds(t *testing.T) {
	bus := Bus.NewBus("main")
	defer bus.Close()

	memory := NewMemory(0x20, 0x0, bus)
	defer memory.Close()

	bus.Write(0x05, 0x02)
	bus.Read(0x15)

	if bus.Data != 0x00 {
		t.Log(bus)
		t.Fail()
	}
}
