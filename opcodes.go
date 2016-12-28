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
	0x00: {Mnemonic: "NOP",			Length: 1, Duration: 4,		Callback: nop},
	0x01: {Mnemonic: "LD BC,d16",	Length: 3, Duration: 12,	Callback: ld_bc_dd},
	0x03: {Mnemonic: "INC BC",		Length: 1, Duration: 8,		Callback: inc_bc},
	0x04: {Mnemonic: "INC B", 		Length: 1, Duration: 4,		Callback: inc_b},
	0x05: {Mnemonic: "DEC B",		Length: 1, Duration: 4, 	Callback: dec_b},
	0x06: {Mnemonic: "LD B,d8",		Length: 2, Duration: 8,		Callback: ld_b_n},
	0x0b: {Mnemonic: "DEC BC",		Length: 1, Duration: 8,		Callback: dec_bc},
	0x0c: {Mnemonic: "INC C",		Length: 1, Duration: 4,		Callback: inc_c},
	0x0d: {Mnemonic: "DEC C",		Length: 1, Duration: 4,		Callback: dec_c},
	0x0e: {Mnemonic: "LD C,d8",		Length: 2, Duration: 8,		Callback: ld_c_n},
	0x11: {Mnemonic: "LD DE,d16",	Length: 3, Duration: 12,	Callback: ld_de_nn},
	0x12: {Mnemonic: "LD (DE),A",	Length: 1, Duration: 8,		Callback: ld_de_a},
	0x13: {Mnemonic: "INC DE",		Length: 1, Duration: 8, 	Callback: inc_de},
	0x15: {Mnemonic: "DEC D",		Length: 1, Duration: 4,		Callback: dec_d},
	0x16: {Mnemonic: "LD D,d8",		Length: 2, Duration: 8,		Callback: ld_d_d},
	0x17: {Mnemonic: "RLA",			Length: 1, Duration: 4,		Callback: rla},
	0x18: {Mnemonic: "JR r8",		Length: 2, Duration: 12,	Callback: jr},
	0x19: {Mnemonic: "ADD HL,DE",	Length: 1, Duration: 8,		Callback: add_hl_de},
	0x1a: {Mnemonic: "LD A,(DE)", 	Length: 1, Duration: 8,		Callback: ld_a_de},
	0x1c: {Mnemonic: "INC E",		Length: 1, Duration: 4,		Callback: inc_e},
	0x1d: {Mnemonic: "DEC E",		Length: 1, Duration: 4,		Callback: dec_e},
	0x1e: {Mnemonic: "LD E,d8",		Length: 2, Duration: 8,		Callback: ld_e_d},
	0x20: {Mnemonic: "JR NZ,r8",	Length: 2, Duration: 12/8,	Callback: jr_nz_n},
	0x21: {Mnemonic: "LD HL,d16", 	Length: 3, Duration: 12,	Callback: ld_hl_nn},
	0x22: {Mnemonic: "LD (HL+),A",	Length: 1, Duration: 8, 	Callback: ld_hli_a},
	0x23: {Mnemonic: "INC HL",		Length: 1, Duration: 8, 	Callback: inc_hl},
	0x24: {Mnemonic: "INC H",		Length: 1, Duration: 4,		Callback: inc_h},
	0x25: {Mnemonic: "DEC H",		Length: 1, Duration: 4,		Callback: dec_h},
	0x28: {Mnemonic: "JR Z,r8",		Length: 2, Duration: 12/8,	Callback: jr_z_r},
	0x2a: {Mnemonic: "LD A,(HL+)",	Length: 1, Duration: 8,		Callback: ld_a_hli},
	0x2c: {Mnemonic: "INC L",		Length: 1, Duration: 4, 	Callback: inc_l},
	0x2e: {Mnemonic: "LD L,d8",		Length: 2, Duration: 8, 	Callback: ld_l_d},
	0x2f: {Mnemonic: "CPL",			Length: 1, Duration: 4,		Callback: cpl},
	0x30: {Mnemonic: "JR NC,r8",	Length: 2, Duration: 12/8,	Callback: jr_nc_r},
	0x31: {Mnemonic: "LD SP,d16", 	Length: 3, Duration: 12,	Callback: ld_sp_nn},
	0x32: {Mnemonic: "LD (HL-),A",	Length: 1, Duration: 8,		Callback: ld_hld_a},
	0x34: {Mnemonic: "INC (HL)", 	Length: 1, Duration: 12,	Callback: inc_r_hl},
	0x35: {Mnemonic: "DEC (HL)",	Length: 1, Duration: 12,	Callback: dec_r_hl},
	0x36: {Mnemonic: "LD (HL),d8", 	Length: 2, Duration: 12,	Callback: ld_hl_d},
	0x3d: {Mnemonic: "DEC A",		Length: 1, Duration: 4,		Callback: dec_a},
	0x3e: {Mnemonic: "LD A,d8",		Length: 2, Duration: 8,		Callback: ld_a_n},
	0x47: {Mnemonic: "LD B,A",		Length: 1, Duration: 4, 	Callback: ld_b_a},
	0x4f: {Mnemonic: "LD C,A",		Length: 1, Duration: 4,		Callback: ld_c_a},
	0x56: {Mnemonic: "LD D,(HL)",	Length: 1, Duration: 8,		Callback: ld_d_hl},
	0x57: {Mnemonic: "LD D,A",		Length: 1, Duration: 4,		Callback: ld_d_a},
	0x5e: {Mnemonic: "LD E,(HL)", 	Length: 1, Duration: 8,		Callback: ld_e_hl},
	0x5f: {Mnemonic: "LD E,A",		Length: 1, Duration: 4, 	Callback: ld_e_a},
	0x67: {Mnemonic: "LD H,A",		Length: 1, Duration: 4,		Callback: ld_h_a},
	0x77: {Mnemonic: "LD (HL),A",	Length: 1, Duration: 8, 	Callback: ld_hl_a},
	0x78: {Mnemonic: "LD A,B",		Length: 1, Duration: 4, 	Callback: ld_a_b},
	0x79: {Mnemonic: "LD A,C",		Length: 1, Duration: 4,		Callback: ld_a_c},
	0x7b: {Mnemonic: "LD A,E",		Length: 1, Duration: 4,		Callback: ld_a_e},
	0x7c: {Mnemonic: "LD A,H",		Length: 1, Duration: 4,		Callback: ld_a_h},
	0x7d: {Mnemonic: "LD A,L",		Length: 1, Duration: 4,		Callback: ld_a_l},
	0x7e: {Mnemonic: "LD A,(HL)",	Length: 1, Duration: 8,		Callback: ld_a_hl},
	0x80: {Mnemonic: "ADD A,B",		Length: 1, Duration: 4, 	Callback: add_a_b},
	0x86: {Mnemonic: "ADD A,(HL)",	Length: 1, Duration: 8,		Callback: add_a_hl},
	0x87: {Mnemonic: "ADD A,A",		Length: 1, Duration: 4, 	Callback: add_a_a},
	0x90: {Mnemonic: "SUB B",		Length: 1, Duration: 4,		Callback: sub_b},
	0xa1: {Mnemonic: "AND C",		Length: 1, Duration: 4,		Callback: and_c},
	0xa7: {Mnemonic: "AND A",		Length: 1, Duration: 4,		Callback: and_a},
	0xa9: {Mnemonic: "XOR C",		Length: 1, Duration: 4,		Callback: xor_c},
	0xaf: {Mnemonic: "XOR A", 		Length: 1, Duration: 4,		Callback: xor_a},
	0xb0: {Mnemonic: "OR B",		Length: 1, Duration: 4, 	Callback: or_b},
	0xb1: {Mnemonic: "OR C",		Length: 1, Duration: 4, 	Callback: or_c},
	0xbe: {Mnemonic: "CP (HL)",		Length: 1, Duration: 8,		Callback: cp_hl},
	0xc1: {Mnemonic: "POP BC",		Length: 1, Duration: 12, 	Callback: pop_bc},
	0xc3: {Mnemonic: "JP a16",		Length: 3, Duration: 16,	Callback: jp_aa},
	0xc5: {Mnemonic: "PUSH BC",		Length: 1, Duration: 16,	Callback: push_bc},
	0xc8: {Mnemonic: "RET Z",		Length: 1, Duration: 20/8,	Callback: ret_z},
	0xc9: {Mnemonic: "RET",			Length: 1, Duration: 16, 	Callback: ret},
	0xca: {Mnemonic: "JP Z,a16",	Length: 3, Duration: 16/12,	Callback: jp_z_aa},
	0xcb: {Mnemonic: "PREFIX CB",	Length: 1, Duration: 4,		Callback: prefixCB},
	0xcd: {Mnemonic: "CALL a16",	Length: 3, Duration: 24,	Callback: call_nn},
	0xd1: {Mnemonic: "POP DE",		Length: 1, Duration: 12,	Callback: pop_de},
	0xd5: {Mnemonic: "PUSH DE", 	Length: 1, Duration: 16,	Callback: push_de},
	0xe0: {Mnemonic: "LDH A,(a8)",	Length: 2, Duration: 12,	Callback: ldh_a_n},
	0xe1: {Mnemonic: "POP HL",		Length: 1, Duration: 12,	Callback: pop_hl},
	0xe2: {Mnemonic: "LD (C),A",	Length: 1, Duration: 8,		Callback: ld_c_a_2},
	0xe5: {Mnemonic: "PUSH HL",		Length: 1, Duration: 16,	Callback: push_hl},
	0xe6: {Mnemonic: "AND d8",		Length: 2, Duration: 8,		Callback: and_d},
	0xe9: {Mnemonic: "JP (HL)",		Length: 1, Duration: 4,		Callback: jp_hl},
	0xea: {Mnemonic: "LD (a16),A", 	Length: 3, Duration: 16,	Callback: ld_aa_a},
	0xef: {Mnemonic: "RST 28H", 	Length: 1, Duration: 16,	Callback: rst_28h},
	0xf0: {Mnemonic: "LDH A,(a8)", 	Length: 2, Duration: 12,	Callback: ldh_a_a},
	0xf1: {Mnemonic: "POP AF",		Length: 1, Duration: 12,	Callback: pop_af},
	0xf3: {Mnemonic: "DI",			Length: 1, Duration: 4,		Callback: di},
	0xf5: {Mnemonic: "PUSH AF",		Length: 1, Duration: 16,	Callback: push_af},
	0xfa: {Mnemonic: "LD A,(a16)",	Length: 3, Duration: 16,	Callback: ld_a_aa},
	0xfb: {Mnemonic: "EI",			Length: 1, Duration: 4, 	Callback: ei},
	0xfe: {Mnemonic: "CP d8",		Length: 2, Duration: 8,		Callback: cp_n},
}

var OpcodesCB = map[uint8]Opcode {
	0x11: {Mnemonic: "RL C",		Length: 1, Duration: 8,		Callback: rl_c},
	0x37: {Mnemonic: "SWAP A",		Length: 1, Duration: 8,		Callback: swap_a},
	0x7c: {Mnemonic: "BIT 7,H",		Length: 1, Duration: 8,		Callback: bit_7_h},
	0x87: {Mnemonic: "RES 0,A",		Length: 1, Duration: 8,		Callback: res_0_a},
	0xfe: {Mnemonic: "SET 7,(HL)",	Length: 1, Duration: 16,	Callback: set_7_hl},
}

func nop(cpu *CPU, data []byte) {
	cpu.Register.M = 1
}

func ld_bc_dd(cpu *CPU, data []byte) {
	cpu.Register.C = data[1]
	cpu.Register.B = data[2]
	cpu.Register.M = 3
}

func inc_bc(cpu *CPU, data []byte) {
	cpu.Register.C = (cpu.Register.C+1)&0xFF
	if cpu.Register.C == 0 {
		cpu.Register.B = (cpu.Register.B+1)&0xFF
	}
	cpu.Register.M = 1
}

func inc_b(cpu *CPU, data []byte) {
	cpu.Register.B++

	if cpu.Register.B != 0 {
		cpu.Register.F = 0x00
	} else {
		// Overflow
		cpu.Register.F = 0x80
	}
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

func dec_bc(cpu *CPU, data []byte) {
	cpu.Register.C = (cpu.Register.C - 1) & 255
	if cpu.Register.C == 255 {
		cpu.Register.B = (cpu.Register.B - 1) & 255
	}
	cpu.Register.M = 1
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

func dec_c(cpu *CPU, data []byte) {
	cpu.Register.C--
	cpu.Register.C &= 0xff
	if cpu.Register.C != 0 {
		cpu.Register.F = 0
	} else {
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

func ld_de_a(cpu *CPU, data []byte) {
	cpu.WriteByte((uint16(cpu.Register.D) << 8) + uint16(cpu.Register.E), cpu.Register.A)
	cpu.Register.M = 2
}

func inc_de(cpu *CPU, data []byte) {
	cpu.Register.E = (cpu.Register.E+1)&0xFF
	if cpu.Register.E == 0 {
		cpu.Register.D = (cpu.Register.D+1)&0xFF
	}
	cpu.Register.M = 1
}

func dec_d(cpu *CPU, data []byte) {
	cpu.Register.D--
	cpu.Register.D &= 0xff
	if cpu.Register.D != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}
}

func ld_d_d(cpu *CPU, data []byte) {
	cpu.Register.D = data[1]
	cpu.Register.M = 2
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

func jr(cpu *CPU, data []byte) {
	i := int32(data[1])
	if i > 127 {
		i =-((^i+1)&255)
	}
	cpu.Register.M = 2

	if i < 0 {
		cpu.Register.PC -= uint16(-i)
	} else {
		cpu.Register.PC += uint16(i)
	}

	//fmt.Printf("Jumped to 0x%x\n", cpu.Register.PC)

	cpu.Register.M++
}

func add_hl_de(cpu *CPU, data []byte) {
	hl := (uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L)
	de := (uint16(cpu.Register.D) << 8) + uint16(cpu.Register.E)

	if uint32(hl) + uint32(de) > 65535 {
		cpu.Register.F |= 0x10
	} else {
		cpu.Register.F &= 0xEF
	}

	hl += de

	cpu.Register.H = uint8(hl >> 8) & 255
	cpu.Register.L = uint8(hl & 255)

	cpu.Register.M = 3
}

func ld_a_de(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.D) << 8) + uint16(cpu.Register.E)
	cpu.Register.A = cpu.ReadByte(addr)
	cpu.Register.M = 2
}

func inc_e(cpu *CPU, data []byte) {
	cpu.Register.E++

	if cpu.Register.E != 0 {
		cpu.Register.F = 0x00
	} else {
		// Overflow
		cpu.Register.F = 0x80
	}
}

func dec_e(cpu *CPU, data []byte) {
	cpu.Register.E--
	cpu.Register.E &= 0xff
	if cpu.Register.E != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}
}

func ld_e_d(cpu *CPU, data []byte) {
	cpu.Register.E = data[1]
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

		//fmt.Printf("Jumped to 0x%x\n", cpu.Register.PC)

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

func inc_h(cpu *CPU, data []byte) {
	cpu.Register.H++

	if cpu.Register.H != 0 {
		cpu.Register.F = 0x00
	} else {
		// Overflow
		cpu.Register.F = 0x80
	}
}

func dec_h(cpu *CPU, data []byte) {
	cpu.Register.H--
	cpu.Register.H &= 0xff
	if cpu.Register.H != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}
}

func jr_z_r(cpu *CPU, data []byte) {
	i := int32(data[1])
	if i > 127 {
		i =-((^i+1)&255)
	}
	cpu.Register.M = 2
	if cpu.Register.F & 0x80 == 0x80 {
		if i < 0 {
			cpu.Register.PC -= uint16(-i)
		} else {
			cpu.Register.PC += uint16(i)
		}

		//fmt.Printf("Jumped to 0x%x (F: 0x%x)\n", cpu.Register.PC, cpu.Register.F)

		cpu.Register.M++
	}
}

func ld_a_hli(cpu *CPU, data []byte) {
	//Z80._r.a=MMU.rb((Z80._r.h<<8)+Z80._r.l); Z80._r.l=(Z80._r.l+1)&255; if(!Z80._r.l) Z80._r.h=(Z80._r.h+1)&255; Z80._r.m=2; },
	cpu.Register.A = cpu.ReadByte((uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L))
	cpu.Register.L = (cpu.Register.L + 1) & 255
	if cpu.Register.L == 0 {
		cpu.Register.H = (cpu.Register.H + 1) & 255
	}
	cpu.Register.M = 2
}

func inc_l(cpu *CPU, data []byte) {
	cpu.Register.L++

	if cpu.Register.L != 0 {
		cpu.Register.F = 0x00
	} else {
		// Overflow
		cpu.Register.F = 0x80
	}
}

func ld_l_d(cpu *CPU, data []byte) {
	cpu.Register.L = data[1]
	cpu.Register.M = 2
}

func cpl(cpu *CPU, data []byte) {
	cpu.Register.A ^= 255

	if cpu.Register.A != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}
	cpu.Register.M = 1
}

func jr_nc_r(cpu *CPU, data []byte) {
	//var i=MMU.rb(Z80._r.pc); if(i>127) i=-((~i+1)&255); Z80._r.pc++; Z80._r.m=2; if((Z80._r.f&0x10)==0x00) { Z80._r.pc+=i; Z80._r.m++; } },
	i := int32(data[1])
	if i > 127 {
		i = -((^i+1)&255)
	}
	cpu.Register.M = 2
	if cpu.Register.F & 0x10 == 0x00 {
		if i < 0 {
			cpu.Register.PC -= uint16(-i)
		} else {
			cpu.Register.PC += uint16(i)
		}
		cpu.Register.M++
	}
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

func inc_r_hl(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L)
	i := cpu.ReadByte(addr) + 1
	i &= 255
	cpu.WriteByte(addr, i)

	if i != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}

	cpu.Register.M = 3
}

func dec_r_hl(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L)
	i := cpu.ReadByte(addr) - 1
	i &= 255

	cpu.WriteByte(addr, i)

	if i != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}

	cpu.Register.M = 3
}

func ld_hl_d(cpu *CPU, data []byte) {
	//MMU.wb((Z80._r.h<<8)+Z80._r.l, MMU.rb(Z80._r.pc)); Z80._r.pc++; Z80._r.m=3;
	cpu.WriteByte((uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L), data[1])
	cpu.Register.M = 3
}

func dec_a(cpu *CPU, data []byte) {
	cpu.Register.A--
	cpu.Register.A &= 0xff
	if cpu.Register.A != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}
}

func ld_a_n(cpu *CPU, data []byte) {
	cpu.Register.A = data[1]
	cpu.Register.M = 2
}

func ld_b_a(cpu *CPU, data []byte) {
	cpu.Register.B = cpu.Register.A
	cpu.Register.M = 1
}

func ld_c_a(cpu *CPU, data []byte) {
	cpu.Register.C = cpu.Register.A
	cpu.Register.M = 1
}

func ld_d_hl(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L)
	cpu.Register.D = cpu.ReadByte(addr)
	cpu.Register.M = 1
}

func ld_d_a(cpu *CPU, data []byte) {
	cpu.Register.D = cpu.Register.A
	cpu.Register.M = 1
}

func ld_e_hl(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L)
	cpu.Register.E = cpu.ReadByte(addr)
	cpu.Register.M = 1
}

func ld_e_a(cpu *CPU, data []byte) {
	cpu.Register.E = cpu.Register.A
	cpu.Register.M = 1
}

func ld_h_a(cpu *CPU, data []byte) {
	cpu.Register.H = cpu.Register.A
	cpu.Register.M = 1
}

func ld_hl_a(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H) << 8)+uint16(cpu.Register.L)
	cpu.WriteByte(addr, cpu.Register.A)
	cpu.Register.M = 2
}

func ld_a_b(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.Register.B
	cpu.Register.M = 1
}

func ld_a_c(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.Register.C
	cpu.Register.M = 1
}

func ld_a_e(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.Register.E
	cpu.Register.M = 1
}

func ld_a_h(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.Register.H
	cpu.Register.M = 1
}

func ld_a_l(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.Register.L
	cpu.Register.M = 1
}

func ld_a_hl(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.ReadByte((uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L))
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

func add_a_hl(cpu *CPU, data []byte) {
	a := cpu.Register.A
	m := cpu.ReadByte((uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L))

	cpu.Register.A += m
	if uint16(a) + uint16(m) > 255 {
		cpu.Register.F = 0x10
	} else {
		cpu.Register.F = 0x00
	}

	cpu.Register.A &= 255

	if cpu.Register.A == 0 {
		cpu.Register.F |= 0x80
	}

	if (cpu.Register.A ^ a ^ m) & 0x10 == 0x10 {
		cpu.Register.F |= 0x20
	}

	cpu.Register.M = 2
}

func add_a_a(cpu *CPU, data []byte) {
	a := cpu.Register.A
	cpu.Register.A += cpu.Register.A
	if uint16(a) + uint16(cpu.Register.A) > 255 {
		cpu.Register.F = 0x10
	} else {
		cpu.Register.F = 0x00
	}
	cpu.Register.A &= 0xff

	if cpu.Register.A == 0 {
		cpu.Register.F |= 0x80
	}

	if cpu.Register.A^cpu.Register.A^a & 0x10 == 0x10 {
		cpu.Register.F |= 0x20
	}

	cpu.Register.M = 1
}

func sub_b(cpu *CPU, data []byte) {
	a := int16(cpu.Register.A)
	cpu.Register.A -= cpu.Register.B

	if a - int16(cpu.Register.B) < 0 {
		cpu.Register.F = 0x50
	} else {
		cpu.Register.F = 0x40
	}

	cpu.Register.A = 0xff

	if cpu.Register.A == 0 {
		cpu.Register.F |= 0x80
	}

	if (cpu.Register.A ^ cpu.Register.B ^ uint8(a)) & 0x10 == 0x10 {
		cpu.Register.F |= 0x10
	}

	cpu.Register.F |= 0x20

	cpu.Register.M = 1
}

func and_c(cpu *CPU, data []byte) {
	cpu.Register.A &= cpu.Register.C
	cpu.Register.A &= 255

	if cpu.Register.A != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}

	cpu.Register.M = 1
}

func and_a(cpu *CPU, data []byte) {
	cpu.Register.A &= cpu.Register.A
	cpu.Register.A &= 255

	if cpu.Register.A != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}

	cpu.Register.M = 1
}

func xor_c(cpu *CPU, data []byte) {
	cpu.Register.A ^= cpu.Register.C
	cpu.Register.A &= 255
	if cpu.Register.A == 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
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

func or_b(cpu *CPU, data []byte) {
	cpu.Register.A |= cpu.Register.B
	cpu.Register.A &= 255

	if cpu.Register.A != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}

	cpu.Register.M = 1
}

func or_c(cpu *CPU, data []byte) {
	cpu.Register.A |= cpu.Register.C
	cpu.Register.A &= 255
	if cpu.Register.A == 0 {
		cpu.Register.F = 0x80
	} else {
		cpu.Register.F = 0x00
	}

	cpu.Register.M = 1
}

func cp_hl(cpu *CPU, data []byte) {
	//var i=Z80._r.a; var m=MMU.rb((Z80._r.h<<8)+Z80._r.l); i-=m; Z80._r.f=(i<0)?0x50:0x40; i&=255; if(!i) Z80._r.f|=0x80; if((Z80._r.a^i^m)&0x10) Z80._r.f|=0x20; Z80._r.m=2; },
	i := uint16(cpu.Register.A)

	addr := (uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L)
	m := cpu.ReadByte(addr)

	fmt.Printf("Compare 0x%x with 0x%x (located at 0x%x)\n", i, m, addr)

	i -= uint16(m)

	if i < 0 {
		cpu.Register.F = 0x50
	} else {
		cpu.Register.F = 0x40
	}

	i &= 0xff

	if i == 0 {
		cpu.Register.F |= 0x80
	}

	if (cpu.Register.A ^ byte(i) ^ m) & 0x10 == 0x10 {
		cpu.Register.F |= 0x20
	}

	cpu.Register.M = 2
}

func pop_bc(cpu *CPU, data []byte) {
	cpu.Register.C = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.B = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.M = 3
}

func jp_aa(cpu *CPU, data []byte) {
	addr := (uint16(data[2]) << 8) + uint16(data[1])
	cpu.Register.PC = addr
	cpu.Register.M = 3
}

func push_bc(cpu *CPU, data []byte) {
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.B)
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.C)
	cpu.Register.M = 3
}

func ret_z(cpu *CPU, data []byte) {
	cpu.Register.M = 1
	if cpu.Register.F & 0x80 == 0x80 {
		cpu.Register.PC = cpu.ReadWord(cpu.Register.SP)
		cpu.Register.SP += 2
		cpu.Register.M += 2
	}
}

func ret(cpu *CPU, data []byte) {
	cpu.Register.PC = cpu.ReadWord(cpu.Register.SP)
	cpu.Register.SP += 2
	cpu.Register.M = 3
}

func jp_z_aa(cpu *CPU, data []byte) {
	cpu.Register.M = 3
	if cpu.Register.F & 0x80 == 0x00 {
		cpu.Register.SP -= 2

		cpu.WriteWord(cpu.Register.SP, cpu.Register.PC)
		cpu.Register.PC = (uint16(data[2]) << 8) + uint16(data[1])
		cpu.Register.M += 2
	}
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

func pop_de(cpu *CPU, data []byte) {
	cpu.Register.E = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.D = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.M = 3
}

func push_de(cpu *CPU, data []byte) {
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.D)
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.E)
	cpu.Register.M = 3
}

func ldh_a_n(cpu *CPU, data []byte) {
	cpu.WriteByte(0xFF00 + uint16(data[1]), cpu.Register.A)
	cpu.Register.M = 3
}

func pop_hl(cpu *CPU, data []byte) {
	cpu.Register.L = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.H = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.M = 3
}

func ld_c_a_2(cpu *CPU, data []byte) {
	cpu.WriteByte(0xFF00+uint16(cpu.Register.C), cpu.Register.A)
	cpu.Register.M = 2
}

func push_hl(cpu *CPU, data []byte) {
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.H)
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.L)
	cpu.Register.M = 3
}

func and_d(cpu *CPU, data []byte) {
	cpu.Register.A &= data[1]
	cpu.Register.A &= 255

	if cpu.Register.A != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}

	cpu.Register.M = 2
}

func jp_hl(cpu *CPU, data []byte) {
	cpu.Register.PC = (uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L)
	cpu.Register.M = 1
}

func ld_aa_a(cpu *CPU, data []byte) {
	addr := (uint16(data[2]) << 8)+uint16(data[1])
	cpu.WriteByte(addr, cpu.Register.A)
	cpu.Register.M = 4
}

func rst_28h(cpu *CPU, data []byte) {
	rsv(cpu)
	cpu.Register.SP -= 2
	cpu.WriteWord(cpu.Register.SP, cpu.Register.PC)
	cpu.Register.PC = 0x28
	cpu.Register.M = 3
}

func ldh_a_a(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.ReadByte(0xFF00 + uint16(data[1]))
	cpu.Register.M = 3
}

func pop_af(cpu *CPU, data []byte) {
	cpu.Register.F = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.A = cpu.ReadByte(cpu.Register.SP)
	cpu.Register.SP++
	cpu.Register.M = 3
}

func di(cpu *CPU, data []byte) {
	cpu.Register.IME = 0
	cpu.Register.M = 1
}

func push_af(cpu *CPU, data []byte) {
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.A)
	cpu.Register.SP--
	cpu.WriteByte(cpu.Register.SP, cpu.Register.F)
	cpu.Register.M = 3
}

func ld_a_aa(cpu *CPU, data []byte) {
	cpu.Register.A = cpu.ReadByte((uint16(data[2]) << 8) + uint16(data[1]))
	cpu.Register.M = 4
}

func ei(cpu *CPU, data []byte) {
	cpu.Register.IME = 1
	cpu.Register.M = 1
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

func swap_a(cpu *CPU, data []byte) {
	tr := cpu.Register.A

	cpu.Register.A = ((tr & 0xF)<<4) | ((tr & 0xF0) >> 4)
	if cpu.Register.A != 0 {
		cpu.Register.F = 0
	} else {
		cpu.Register.F = 0x80
	}

	cpu.Register.M = 1
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

func res_0_a(cpu *CPU, data []byte) {
	cpu.Register.A &= 0xFE
	cpu.Register.M = 2
}

func set_7_hl(cpu *CPU, data []byte) {
	addr := (uint16(cpu.Register.H) << 8) + uint16(cpu.Register.L)
	i := cpu.ReadByte(addr)
	i |= 0x80
	cpu.WriteByte(addr, i)

	cpu.Register.M = 4
}

func rsv(cpu *CPU) {
	cpu.RSV.A = cpu.Register.A
	cpu.RSV.B = cpu.Register.B
	cpu.RSV.C = cpu.Register.C
	cpu.RSV.D = cpu.Register.D
	cpu.RSV.E = cpu.Register.E
	cpu.RSV.F = cpu.Register.F
	cpu.RSV.H = cpu.Register.H
	cpu.RSV.L = cpu.Register.L
}

func rss(cpu *CPU) {
	cpu.Register.A = cpu.RSV.A
	cpu.Register.B = cpu.RSV.B
	cpu.Register.C = cpu.RSV.C
	cpu.Register.D = cpu.RSV.D
	cpu.Register.E = cpu.RSV.E
	cpu.Register.F = cpu.RSV.F
	cpu.Register.H = cpu.RSV.H
	cpu.Register.L = cpu.RSV.L
}