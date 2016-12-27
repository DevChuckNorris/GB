package main

import (
	"fmt"
)

func main() {
	fmt.Println("GB Emulator v0.1")

	cpu := NewCPU()
	cpu.LoadBootLoader("boot.gb")

	cpu.Run()
}
