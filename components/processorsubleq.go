package components

import (
	"fmt"
)

/**
 *
 */
type subleqProcessor struct {
	PC AddressLinesType // The Program Counter
	A  DataLinesType    // The Accumulator
	Z  bool             // The Z flag, set if the last operation result was zero
	N  bool             // The N flag, set if the last operation result was negative
	processorAbstract
}

/**
 * Constructor.
 */
func NewSubleqProcessor(bus *Bus) *subleqProcessor {
	processor := subleqProcessor{}
	processor.initProcessor(bus, processor.onReset)

	return &processor
}

func (processor *subleqProcessor) onReset() {
	processor.PC = 0
	processor.A = 0
	processor.Z = false
	processor.N = false

	go processor.process()
}

/**
 * Executes the actual program.
 */
func (processor *subleqProcessor) process() {
	for {
		// If we jump to 0xFF this is a break, this go routine ends here.
		if processor.PC == 0xFF {
			processor.bus.Brk <- true
			return
		}

		// Fetch the three operands addresses
		opAAddress := AddressLinesType(processor.bus.Read(processor.PC))
		opBAddress := AddressLinesType(processor.bus.Read(processor.PC + 1))
		branchAddress := AddressLinesType(processor.bus.Read(processor.PC + 2))
		processor.PC += 3

		// Fetch the A/B values
		opA := processor.bus.Read(opAAddress)
		opB := processor.bus.Read(opBAddress)

		// Compute result
		processor.A = opB - opA
		processor.Z = opB == opA
		processor.N = opB < opA

		// Write the result back to opA
		processor.bus.Write(opAAddress, processor.A)

		// Branch is less or equal
		if processor.Z || processor.N {
			processor.PC = branchAddress
		}
	}
}

/**
 * String representation of the processor status.
 */
func (processor *subleqProcessor) String() string {
	return fmt.Sprintf("PC:%#X A:%#X Z:%t N:%t", processor.PC, processor.A, processor.Z, processor.N)
}
