package Processor

import (
	"Go-OISC/Bus"
	"time"
)

type subleqProcessor struct {
	bus *Bus.Bus
	PC  Bus.AddressLinesType
}

func NewSubleqProcessor(bus Bus.Bus) *subleqProcessor {
	processor := subleqProcessor{&bus, 0}

	go processor.process()

	return &processor
}

func (processor *subleqProcessor) process() {
	for {
		if processor.bus.Halt {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if processor.PC == 0xFF {
			processor.bus.Halt = true
		}

		opA := processor.bus.Read(processor.PC)
		opB := processor.bus.Read(processor.PC + 1)
		result := opB - opA
		processor.bus.Write(processor.PC, result)
		if result <= 0 {
			processor.PC = Bus.AddressLinesType(processor.bus.Read(processor.PC + 2))
			continue
		}

		processor.PC += 3
	}
}
