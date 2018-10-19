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
		0x00, 0x01, 0xFF,
		0x00, 0x01, 0xFF,
		0x00, 0x01, 0xFF,
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

	for address := 0; address < 0x8; address++ {
		val := bus.Read(Bus.AddressLinesType(address))
		fmt.Println(val)
	}
}
