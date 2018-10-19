package main

import (
	"Go-OISC/Bus"
	"Go-OISC/Memory"
	"Go-OISC/Processor"
	"fmt"
	"time"
)

func main() {
	bus := Bus.NewBus("main")
	Memory.NewMemory(0xFF, 0x00, bus)
	Processor.NewSubleqProcessor(bus)

	for address, instruction := range []Bus.DataLinesType{
		0x0F, 0x0C, 0xFF, // 0x00	SUBLEQ 0C 09 FF
		0x10, 0X0D, 0xFF, // 0x03	SUBLEQ 0D 09 FF
		0x0E, 0x11, 0xFF, // 0x06	SUBLEQ 0E 11 FF
		0x11, 0x11, 0xFF, // 0x09   SUBLEQ 11 11 FF
		1,    // 0x0C	Operand A
		1,    // 0x0D	Operand B
		0x00, // 0x0E	Result
		0x00, // 0x0F	Swap 1
		0x00, // 0x10	Swap 2
		0x01, // 0x11   Swap 3
	} {
		bus.Write(Bus.AddressLinesType(address), instruction)
	}

	bus.ResetBus()
	for {
		if bus.Halt {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	for address := 0; address < 0xB; address++ {
		val := bus.Read(Bus.AddressLinesType(address))
		fmt.Println(val)
	}

	fmt.Print("result:")
	fmt.Println(bus.Read(Bus.AddressLinesType(0x0E)))
}
