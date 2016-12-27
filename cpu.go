package main

import (
	"os"
	"fmt"
)

type Register struct {
	A	byte
	B	byte
	C	byte
	D	byte
	E	byte
	H	byte
	L 	byte
	F	byte

	PC	uint16
	SP	uint16

	M	byte
	T	byte
}

type CPU struct {
	ram		[]byte

	isCB	bool

	Register Register
}

func NewCPU() *CPU {
	cpu := new(CPU)

	cpu.ram = make([]byte, 65535)

	return cpu
}

func (c *CPU) WriteByte(addr uint16, data byte) {
	fmt.Printf("Write 0x%x to 0x%x\n", data, addr)

	c.ram[addr] = data
}

func (c *CPU) ReadByte(addr uint16) byte {
	return c.ram[addr]
}

func (c *CPU) LoadBootLoader(file string) {
	f, err := os.Open("boot.gb")
	if err != nil {
		panic(err)
	}

	n1, err := f.Read(c.ram)
	if err != nil {
		panic(err)
	}

	if n1 != 256 {
		panic(fmt.Errorf("BootLoader is not 256 byte long, is %d long", n1))
	}
}

func (c *CPU) Run() {
	c.Register.PC = 0

	for {
		code := c.ram[c.Register.PC]

		var opcode Opcode
		var ok bool

		if c.isCB {
			opcode, ok = OpcodesCB[code]
			c.isCB = false

			if !ok {
				fmt.Printf("Unknown cb-opcode 0x%x\n", code)
				return
			}
		} else {
			opcode, ok = Opcodes[code]

			if !ok {
				fmt.Printf("Unknown opcode 0x%x\n", code)
				return
			}
		}

		fmt.Printf("Opcode is %s\n", opcode.Mnemonic)

		data := make([]byte, opcode.Length)
		end := c.Register.PC + uint16(opcode.Length)
		i := 0
		for c.Register.PC < end {
			data[i] = c.ram[c.Register.PC]
			c.Register.PC++
			i++
		}

		if opcode.Callback != nil {
			opcode.Callback(c, data)
		} else {
			fmt.Println("Not implemented!")
		}
	}
}

func (c *CPU) ActivateCB() {
	c.isCB = true
}
