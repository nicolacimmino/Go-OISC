package Memory

import "testing"

func TestSetGet(t *testing.T) {
	m := NewMemory(0x100)

	m.Set(5, 12)

	if m.Get(5) != 12 {
		t.Fail()
	}
}

func TestWriteAddressWraps(t *testing.T) {
	m := NewMemory(0x100)

	m.Set(0x101, 12)

	if m.Get(1) != 12 {
		t.Fail()
	}
}

func TestReadAddressWraps(t *testing.T) {
	m := NewMemory(0x100)

	m.Set(0x1, 12)

	if m.Get(0x101) != 12 {
		t.Fail()
	}
}
