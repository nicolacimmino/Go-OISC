package Memory

type Memory struct {
	data []byte
	size uint16
}

func NewMemory(size uint16) *Memory {
	m := Memory{make([]byte, size), size};
	return &m
}

func (m Memory) Set(address uint16, value byte) {
	address = address % m.size
	m.data[address] = value
}

func (m Memory) Get(address uint16) byte {
	address = address % m.size
	return m.data[address]
}
