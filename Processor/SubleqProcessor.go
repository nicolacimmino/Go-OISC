package Processor

import (
	"Go-OISC/Bus"
	"fmt"
	"time"
)

type subleqProcessor struct {
	bus *Bus.Bus
	PC  Bus.AddressLinesType
	A   Bus.DataLinesType
	Z   bool
	N   bool
}

func NewSubleqProcessor(bus *Bus.Bus) *subleqProcessor {
	processor := subleqProcessor{}
	processor.bus = bus

	resetSubscriber := make(Bus.Subscriber)
	bus.SubscribeToReset(resetSubscriber)

	go processor.process()
	go processor.processResetLine(resetSubscriber)

	return &processor
}

func (processor *subleqProcessor) processResetLine(resetSubscriber Bus.Subscriber) {
	for {
		reset, ok := <-resetSubscriber
		if !ok {
			break
		}

		if reset {
			processor.PC = 0
			processor.bus.Halt = false
		}
	}
}

func (processor *subleqProcessor) process() {
	for {
		if processor.bus.Halt {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if processor.PC == 0xFF {
			processor.bus.Halt = true
			continue
		}

		fmt.Println(processor)

		opA := processor.bus.Read(processor.PC)
		opB := processor.bus.Read(processor.PC + 1)
		processor.A = opB - opA
		processor.bus.Write(processor.PC, processor.A)
		processor.Z = processor.A == 0
		processor.N = opB < opA

		if processor.Z || processor.N {
			processor.PC = Bus.AddressLinesType(processor.bus.Read(processor.PC + 2))
			continue
		}

		processor.PC += 3
	}
}

func (processor *subleqProcessor) String() string {
	return fmt.Sprintf("PC:%#X A:%#X Z:%t N:%t HALT:%t", processor.PC, processor.A, processor.Z, processor.N, processor.bus.Halt)
}
