package main

import (
	"Go-OISC/Bus"
	"Go-OISC/Memory"
	"fmt"
	"time"
)

func main() {
	bus := Bus.NewBus("main")
	Memory.NewMemory(0xFF, 0x00, bus)
	//Processor.NewSubleqProcessor(*bus)

	bus.Write(0x00, 10)
	bus.Write(0x01, 35)
	bus.Write(0x02, 03)
	bus.Write(0x03, 40)
	bus.Write(0x04, 20)
	bus.Write(0x05, 0xFF)

	//bus.Halt = false
	time.Sleep(1 * time.Second)
	//bus.Halt = true

	for address := 0; address < 0x8; address++ {
		val := bus.Read(Bus.AddressLinesType(address))
		fmt.Println(val)
		//time.Sleep(1 * time.Second)
	}
}
