package Bus

type addressLinesType uint16
type dataLinesType uint16

type bus struct {
	Address     addressLinesType
	Data        dataLinesType
	Clock       chan bool
	subscribers []chan busCycle
}

type busCycle struct {
	Address addressLinesType
	Data    dataLinesType
	Clock   bool
}

func NewBus() *bus {
	bus := bus{0, 0, make(chan bool), make([]chan busCycle, 0)}

	go bus.process()

	return &bus
}

func (b *bus) Register(subscriber chan busCycle) {
	b.subscribers = append(b.subscribers, subscriber)
}

func (b *bus) process() {
	for {
		clock := <-b.Clock
		cycle := busCycle{b.Address, b.Data, clock}

		for _, subscriber := range b.subscribers {
			subscriber <- cycle
		}

	}
}
