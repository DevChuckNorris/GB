package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("GB Emulator v0.1")

	args := os.Args[1:]

	cpu := NewCPU()
	cpu.LoadBootLoader("boot.gb")
	cpu.LoadROM(args[0])

	cpu.Run()
}
