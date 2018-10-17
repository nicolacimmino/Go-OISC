package Processor

import "Go-OISC/Memory"

type processor struct {
	m *Memory.Memory
}

func NewProcessor() *processor {
	processor := processor{Memory.NewMemory(0x100)}

	return &processor
}

func (processor *processor) Execute() {

}
