package main

import (
	"os"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
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

	IME	byte
}

type MBC struct {
	rombank		byte
	rambank		byte
	ramon		uint16
	mode		uint16
}

type CPU struct {
	ram		[]byte
	rom		[]byte

	isCB	bool

	gpu		 	*GPU
	Register	Register
	RSV			Register

	Clock	uint16

	Ie		byte
	If		byte	// Interrupt flags
	inBios  bool

	romoffs uint16
	ramoffs uint16

	mbc1	MBC
}

func NewCPU() *CPU {
	cpu := new(CPU)

	cpu.Register.IME = 1
	cpu.romoffs = 0x4000
	cpu.mbc1 = MBC{}

	cpu.ram = make([]byte, 65535)
	cpu.gpu = NewGPU(cpu)

	return cpu
}

func (c *CPU) LoadROM(file string) {
	c.rom = make([]byte, 32768)

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	n1, err := f.Read(c.rom)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Read %d rom\n", n1)
}

func (c *CPU) WriteWord(addr uint16, data uint16) {
	fmt.Printf("Write 16 bits 0x%x to 0x%x\n", data, addr)

	c.ram[addr] = uint8(data & 0xff)
	c.ram[addr+1] = uint8(data>>8)
}

func (c *CPU) WriteByte(addr uint16, data byte) {
	//fmt.Printf("Write 0x%x to 0x%x\n", data, addr)

	switch addr & 0xf000 {
	case 0x0000, 0x1000:
		if data & 0xf == 0xa {
			c.mbc1.ramon = 1
		} else {
			c.mbc1.ramon = 0
		}

		break
	case 0x2000, 0x3000:
		c.mbc1.rombank &= 0x60
		data &= 0x1f
		if data == 0 {
			data = 1
		}
		c.mbc1.rombank |= data
		c.romoffs = uint16(c.mbc1.rombank) * 0x4000
		break
	case 0x4000, 0x5000:
		if c.mbc1.mode != 0 {
			c.mbc1.rambank = data&3
			c.ramoffs = uint16(c.mbc1.rambank) * 0x2000
		} else {
			c.mbc1.rombank &= 0x1f
			c.mbc1.rombank |= (data&3) << 5
			c.romoffs = uint16(c.mbc1.rombank) * 0x4000
		}
		break
	case 0x6000, 0x7000:
		c.mbc1.mode = uint16(data) & 1
		break
	case 0x8000, 0x9000:	// vram
		c.gpu.WriteVram(addr & 0x1fff, data)
		c.gpu.UpdateTile(addr & 0x1fff, data)
		break
	case 0xf000:
		switch addr & 0x0f00 {
		case 0x000, 0x100, 0x200, 0x300, 0x400, 0x500, 0x600, 0x700, 0x800, 0x900, 0xa00, 0xb00, 0xc00, 0xd00:
			c.ram[addr] = data
			return
		case 0xe00:
			if (addr&0xFF)<0xA0 {
				c.gpu.WriteOam(addr & 0xFF, data)
			}
			c.gpu.UpdateOam(addr,data)
			return
		case 0xf00:
			if addr == 0xffff {
				c.Ie = data
				return
			} else if addr > 0xFF7F {
				fmt.Printf("Write 0x%x to 0x%x\n", data, addr)
				if addr == 0xff85 {
					sdl.Delay(1000)
				}
				c.ram[addr] = data
				return
			} else {
				switch addr & 0xf0 {
				case 0x00:
					switch addr & 0xf {
					case 0:
						return	// JOYP
					case 4: case 5: case 6: case 7:
						return	// Timer
					case 15:
						c.If = data
						return
					default:
						return
					}
				case 0x10, 0x20, 0x30:
					return
				case 0x40, 0x50, 0x60, 0x70:
					c.gpu.WriteByte(addr, data)
					return
				}
			}
		}
	}

	c.ram[addr] = data
}

func (c *CPU) ReadWord(addr uint16) uint16 {
	return uint16(c.ReadByte(addr)) + (uint16(c.ReadByte(addr+1)) << 8)
}

func (c *CPU) ReadByte(addr uint16) byte {
	switch addr & 0xf000 {
	case 0x0000:
		if c.inBios {
			if addr < 0x0100 {
				return c.ram[addr]
			} else if c.Register.PC == 0x100 {
				c.inBios = false
				fmt.Println("Leave bios/bootloader")
			}

			//fmt.Printf("Read rom at 0x%x\n", addr)
			return c.rom[addr]
		} else {
			//fmt.Printf("Read rom at 0x%x\n", addr)
			return c.rom[addr]
		}
	case 0x1000, 0x2000, 0x3000:
		//fmt.Printf("Read rom at 0x%x\n", addr)
		return c.rom[addr]
	case 0x4000, 0x5000, 0x6000, 0x7000:
		return c.rom[c.romoffs+(addr & 0x3fff)]
	case 0xf000:
		switch addr & 0x0f00 {
		case 0x000, 0x100, 0x200, 0x300, 0x400, 0x500, 0x600, 0x700, 0x800, 0x900, 0xa00, 0xb00, 0xc00, 0xd00:
			return c.ram[addr]
		case 0xe00:
			if addr & 0xff < 0xa0 {
				return c.gpu.ReadOam(addr & 0xff)
			} else {
				return 0
			}
		case 0xf00:
			if addr == 0xffff {
				return c.Ie
			} else if addr > 0xFF7F {
				//fmt.Printf("Read 0x%x from 0x%x\n", c.ram[addr], addr)
				return c.ram[addr]
			} else {
				switch addr & 0xf0 {
				case 0x00:
					switch addr & 0xf {
					case 0:
						return 0	// JOYP
					case 4: case 5: case 6: case 7:
						return 0	// Timer
					case 15:
						return c.If
					default:
						return 0
					}
				case 0x10, 0x20, 0x30:
					return 0
				case 0x40, 0x50, 0x60, 0x70:
					return c.gpu.ReadByte(addr)
				}
			}
		}
	}

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

	c.inBios = true
}

func (c *CPU) Run() {
	c.Register.PC = 0

	for {
		code := c.ReadByte(c.Register.PC)

		var opcode Opcode
		var ok bool

		if c.isCB {
			opcode, ok = OpcodesCB[code]
			c.isCB = false

			if !ok {
				fmt.Printf("Unknown cb-opcode 0x%x at 0x%x\n", code, c.Register.PC)
				return
			}
		} else {
			opcode, ok = Opcodes[code]

			if !ok {
				fmt.Printf("Unknown opcode 0x%x at 0x%x\n", code, c.Register.PC)
				return
			}
		}

		if c.Register.PC >= 0xe9 {
			//fmt.Printf("Opcode is %s at 0x%x\n", opcode.Mnemonic, c.Register.PC)
			//return
		}

		data := make([]byte, opcode.Length)
		end := c.Register.PC + uint16(opcode.Length)
		i := 0
		for c.Register.PC < end {
			data[i] = c.ReadByte(c.Register.PC)
			c.Register.PC++
			i++
		}

		if opcode.Callback != nil {
			opcode.Callback(c, data)

			c.Clock += uint16(c.Register.M)
		} else {
			fmt.Println("Not implemented!")
		}

		// GPU action
		c.gpu.CheckLine()

		if !c.gpu.IsRunning() {
			break
		}
	}
}

func (c *CPU) ActivateCB() {
	c.isCB = true
}
