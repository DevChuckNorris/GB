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
	0x05: {Mnemonic: "DEC B",		Length: 1, Duration: 4, 	Callback: dec_b},
	0x06: {Mnemonic: "LD B,d8",		Length: 2, Duration: 8,		Callback: ld_b_n},
	0x0c: {Mnemonic: "INC C",		Length: 1, Duration: 4,		Callback: inc_c},
	0x0e: {Mnemonic: "LD C,d8",		Length: 2, Duration: 8,		Callback: ld_c_n},
	0x11: {Mnemonic: "LD DE,d16",	Length: 3, Duration: 12,	Callback: ld_de_nn},
	0x13: {Mnemonic: "INC DE",		Length: 1, Duration: 8, 	Callback: inc_de},
	0x17: {Mnemonic: "RLA",			Length: 1, Duration: 4,		Callback: rla},
	0x1a: {Mnemonic: "LD A,(DE)", 	Length: 1, Duration: 8,		Callback: ld_a_de},
	0x20: {Mnemonic: "JR NZ,r8",	Length: 2, Duration: 12/8,	Callback: jr_nz_n},
	0x21: {Mnemonic: "LD HL,d16", 	Length: 3, Duration: 12,	Callback: ld_hl_nn},
	0x22: {Mnemonic: "LD (HL+),A",	Length: 1, Duration: 8, 	Callback: ld_hli_a},
	0x23: {Mnemonic: "INC HL",		Length: 1, Duration: 8, 	Callback: inc_hl},
	0x31: {Mnemonic: "LD SP,d16", 	Length: 3, Duration: 12,	Callback: ld_sp_nn},
	0x32: {Mnemonic: "LD (HL-),A",	Length: 1, Duration: 8,		Callback: ld_hld_a},
	0x3e: {Mnemonic: "LD A,d8",		Length: 2, Duration: 8,		Callback: ld_a_n},
	0x4f: {Mnemonic: "LD C,A",		Length: 1, Duration: 4,		Callback: ld_c_a},
	0x77: {Mnemonic: "LD (HL),A",	Length: 1, Duration: 8, 	Callback: ld_hl_a},
	0x7b: {Mnemonic: "LD A,E",		Length: 1, Duration: 4,		Callback: ld_a_e},
	//0x80: {Mnemonic: "ADD A,B",		Length: 1, Duration: 4, 	Callback: add_a_b},
	0xaf: {Mnemonic: "XOR A", 		Length: 1, Duration: 4,		Callback: xor_a},
	0xc1: {Mnemonic: "POP BC",		Length: 1, Duration: 12, 	Callback: pop_bc},
	0xc5: {Mnemonic: "PUSH BC",		Length: 1, Duration: 16,	Callback: push_bc},
	0xc9: {Mnemonic: "RET",			Length: 1, Duration: 16, 	Callback: ret},
	0xcb: {Mnemonic: "PREFIX CB",	Length: 1, Duration: 4,		Callback: prefixCB},
	0xcd: {Mnemonic: "CALL a16",	Length: 3, Duration: 24,	Callback: call_nn},
	0xe0: {Mnemonic: "LDH A,(a8)",	Length: 2, Duration: 12,	Callback: ldh_a_n},
	0xe2: {Mnemonic: "LD (C),A",	Length: 1, Duration: 8,		Callback: ld_c_a_2},
	0xfe: {Mnemonic: "CP d8",		Length: 2, Duration: 8,		Callback: cp_n},
}

var OpcodesCB = map[uint8]Opcode {
	0x11: {Mnemonic: "RL C",		Length: 1, Duration: 8,		Callback: rl_c},
	0x7c: {Mnemonic: "BIT 7,H",		Length: 1, Duration: 8,		Callback: bit_7_h},
}

func dec_b(cpu *CPU, data []byte) {
	cpu.Register.B--
	cpu.Register.B &= 0xff
	if cpu.Register.B != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}
}

func ld_b_n(cpu *CPU, data []byte) {
	cpu.Register.B = data[1]
	cpu.Register.M = 2
}

func inc_c(cpu *CPU, data []byte) {
	cpu.Register.C++

	if cpu.Register.C != 0 {
		cpu.Register.F = 0x00
	} else {
		// Overflow
		cpu.Register.F = 0x80
	}
}

func ld_c_n(cpu *CPU, data []byte) {
	cpu.Register.C = data[1]
	cpu.Register.M = 2
}

func ld_de_nn(cpu *CPU, data []byte) {
	cpu.Register.E = data[1]
	cpu.Register.D = data[2]
	cpu.Register.M = 3
}

func inc_de(cpu *CPU, data []byte) {
	cpu.Register.E = (cpu.Register.E+1)&0xFF
	if cpu.Register.E == 0 {
		cpu.Register.D = (cpu.Register.D+1)&0xFF
	}
	cpu.Register.M = 1
}

func rla(cpu *CPU, data []byte) {
	var ci = byte(0)
	if cpu.Register.F & 0x10 == 0x10 {
		ci = 1
	}

	var co = byte(0)
	if cpu.Register.A & 0x80 == 0x80 {
		co = 0x10
	}

	cpu.Register.A = (cpu.Register.A << 1)+ci
	cpu.Register.A &= 0xff

	cpu.Register.F = (cpu.Register.F & 0xEF)+co
	cpu.Register.M = 1
}

func ld_a_de(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.D) << 8) + uint16(cpu.Register.E)
	cpu.Register.A = cpu.ReadByte(addr)
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

func ld_hli_a(cpu *CPU, data []byte) {
	cpu.WriteByte((uint16(cpu.Register.H) << 8)+uint16(cpu.Register.L), cpu.Register.A)
	cpu.Register.L = (cpu.Register.L + 1) & 0xff
	if cpu.Register.L == 0 {
		cpu.Register.H = (cpu.Register.H + 1) & 0xff
	}
	cpu.Register.M = 2
}

func inc_hl(cpu *CPU, data []byte) {
	cpu.Register.L = (cpu.Register.L + 1) & 0xFF
	if cpu.Register.L == 0 {
		cpu.Register.H = (cpu.Register.H + 1) & 0xFF
	}

	cpu.Register.M = 1
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

func ld_c_a(cpu *CPU, data []byte) {
	cpu.Register.C = cpu.Register.A
	cpu.Register.M = 1
}

func ld_hl_a(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H) << 8)+uint16(cpu.Register.L)
	cpu.WriteByte(addr, cpu.Register.A)
	cpu.Register.M = 2
}

func ld_a_e(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.Register.E
	cpu.Register.M = 1
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

func pop_bc(cpu *CPU, data []byte) {
	cpu.Register.C = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.B = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.M = 3
}

func push_bc(cpu *CPU, data []byte) {
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.B)
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.C)
	cpu.Register.M = 3
}

func ret(cpu *CPU, data []byte) {
	cpu.Register.PC = cpu.ReadWord(cpu.Register.SP)
	cpu.Register.SP += 2
	cpu.Register.M = 3
}

func prefixCB(cpu *CPU, data []byte) {
	cpu.ActivateCB()
}

func call_nn(cpu *CPU, data []byte) {
	cpu.Register.SP -= 2
	cpu.WriteWord(cpu.Register.SP, cpu.Register.PC)

	addr := (uint16(data[2]) << 8)+uint16(data[1])
	cpu.Register.PC = addr
	cpu.Register.M = 5
}

func ldh_a_n(cpu *CPU, data []byte) {
	cpu.WriteByte(0xFF00 + uint16(data[1]), cpu.Register.A)
	cpu.Register.M = 3
}

func ld_c_a_2(cpu *CPU, data []byte) {
	cpu.WriteByte(0xFF00+uint16(cpu.Register.C), cpu.Register.A)
	cpu.Register.M = 2
}

func cp_n(cpu *CPU, data []byte) {
	var i = int16(cpu.Register.A)
	var m = int16(data[1])

	i -= m

	if i < 0 {
		cpu.Register.F = 0x50
	} else {
		cpu.Register.F = 0x40
	}

	i&=255

	if i == 0 {
		cpu.Register.F |= 0x80
	}

	if (int16(cpu.Register.A) ^ i ^ m) & 0x10 == 0x10 {
		cpu.Register.F |= 0x20
	}

	cpu.Register.M = 2
}

func rl_c(cpu *CPU, data []byte) {
	var ci = byte(0)
	if cpu.Register.F & 0x10 == 0x10 {
		ci = 1
	}

	var co = byte(0)
	if cpu.Register.C & 0x80 == 0x80 {
		co = 0x10
	}

	cpu.Register.C = (cpu.Register.C << 1)+ci
	cpu.Register.C &= 0xff

	if cpu.Register.C != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}

	cpu.Register.F = (cpu.Register.F & 0xEF) + co
	cpu.Register.M = 2
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