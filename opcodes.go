package main

import "fmt"

type OpcodeFunction func(*CPU, []byte)

type Opcode struct {
	Mnemonic 	string
	Length 		uint8
	Duration	uint8
	Callback	OpcodeFunction
}

var Opcodes = map[uint8]Opcode {
	0x0e: {Mnemonic: "LD C,d8",		Length: 2, Duration: 8,		Callback: ld_c_n},
	0x20: {Mnemonic: "JR NZ,r8",	Length: 2, Duration: 12/8,	Callback: jr_nz_n},
	0x21: {Mnemonic: "LD HL,d16", 	Length: 3, Duration: 12,	Callback: ld_hl_nn},
	0x31: {Mnemonic: "LD SP,d16", 	Length: 3, Duration: 12,	Callback: ld_sp_nn},
	0x32: {Mnemonic: "LD (HL-),A",	Length: 1, Duration: 8,		Callback: ld_hld_a},
	0x3e: {Mnemonic: "LD A,d8",		Length: 2, Duration: 8,		Callback: ld_a_n},
	0x77: {Mnemonic: "LD (HL),A",	Length: 1, Duration: 8, 	Callback: ld_hl_a},
	//0x80: {Mnemonic: "ADD A,B",		Length: 1, Duration: 4, 	Callback: add_a_b},
	0xaf: {Mnemonic: "XOR A", 		Length: 1, Duration: 4,		Callback: xor_a},
	0xcb: {Mnemonic: "PREFIX CB",	Length: 1, Duration: 4,		Callback: prefixCB},
	//0xe2: {Mnemonic: "LD (C),A",	Length: 2, Duration: 8,		Callback: ld_c_a},
}

var OpcodesCB = map[uint8]Opcode {
	0x7c: {Mnemonic: "BIT 7,H",		Length: 1, Duration: 8,		Callback: bit_7_h},
}

func ld_c_n(cpu *CPU, data []byte) {
	cpu.Register.C = data[1]
	cpu.Register.M = 2
}

func jr_nz_n(cpu *CPU, data []byte) {
	i := int32(data[1])
	if i > 127 {
		i =-((^i+1)&255)
	}
	cpu.Register.M = 2
	if cpu.Register.F & 0x80 == 0x00 {
		if i < 0 {
			cpu.Register.PC -= uint16(-i)
		} else {
			cpu.Register.PC += uint16(i)
		}

		fmt.Printf("Jumped to 0x%x\n", cpu.Register.PC)

		cpu.Register.M++
	}
}

func ld_hl_nn(cpu *CPU, data []byte) {
	cpu.Register.L = data[1]
	cpu.Register.H = data[2]
	cpu.Register.M = 3
}

func ld_sp_nn(cpu *CPU, data []byte) {
	cpu.Register.SP = uint16(data[1]) | (uint16(data[2]) << 8)
	cpu.Register.M = 3
}

func ld_hld_a(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H)<<8)+uint16(cpu.Register.L)

	cpu.WriteByte(addr, cpu.Register.A)
	cpu.Register.L = (cpu.Register.L-1)&0xff
	if cpu.Register.L == 255 {
		cpu.Register.H = (cpu.Register.H-1)&0xff
	}
	cpu.Register.M = 2
}

func ld_a_n(cpu *CPU, data []byte) {
	cpu.Register.A = data[1]
	cpu.Register.M = 2
}

func ld_hl_a(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H) << 8)+uint16(cpu.Register.L)
	cpu.WriteByte(addr, cpu.Register.A)
	cpu.Register.M = 2
}

func add_a_b(cpu *CPU, data []byte) {
	a := cpu.Register.A
	cpu.Register.A += cpu.Register.B
	if uint16(a) + uint16(cpu.Register.B) > 255 {
		cpu.Register.F = 0x10
	} else {
		cpu.Register.F = 0x00
	}
	cpu.Register.A &= 0xff

	if cpu.Register.A == 0 {
		cpu.Register.F |= 0x80
	}

	if cpu.Register.A^cpu.Register.B^a & 0x10 == 0x10 {
		cpu.Register.F |= 0x20
	}

	cpu.Register.M = 1
}

func xor_a(cpu *CPU, data []byte) {
	cpu.Register.A ^= cpu.Register.A
	cpu.Register.A &= 255
	if cpu.Register.A == 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}
	cpu.Register.M = 1
}

func prefixCB(cpu *CPU, data []byte) {
	cpu.ActivateCB()
}

func ld_c_a(cpu *CPU, data []byte) {

}

func bit_7_h(cpu *CPU, data []byte) {
	cpu.Register.F &= 0x1f
	cpu.Register.F |= 0x20
	if cpu.Register.H & 0x80 == 0x80 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}
	cpu.Register.M = 2
}