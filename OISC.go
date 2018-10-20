package main

import (
	"Go-OISC/components"
	"fmt"
)

func main() {
	bus := components.NewBus("main")
	components.NewMemory(0xFF, 0x00, bus)
	components.NewSubleqProcessor(bus)

	for address, instruction := range []components.DataLinesType{
		0x0F, 0x0C, 0xFF, // 0x00	SUBLEQ 0C 09 FF
		0x10, 0X0D, 0xFF, // 0x03	SUBLEQ 0D 09 FF
		0x0E, 0x11, 0xFF, // 0x06	SUBLEQ 0E 11 FF
		0x11, 0x11, 0xFF, // 0x09   SUBLEQ 11 11 FF
		1, // 0x0C	Operand A
		1, // 0x0D	Operand B
		0, // 0x0E	Result
		0, // 0x0F	Swap 1
		0, // 0x10	Swap 2
		1, // 0x11   Swap 3
	} {
		bus.Write(components.AddressLinesType(address), instruction)
	}

	bus.ResetBus()
	<-bus.Brk

	fmt.Print("result:")
	fmt.Println(bus.Read(components.AddressLinesType(0x0E)))
}
