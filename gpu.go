package main

import (
	"github.com/veandco/go-sdl2/sdl"
	//"fmt"
	"fmt"
)

type ObjData struct {
	x		int16
	y		int16
	tile	byte
	palette bool
	xflip	bool
	yflip	bool
	prio	bool
	num		byte
}

type GPU struct {
	cpu		*CPU

	running	bool
	window	*sdl.Window
	surface	*sdl.Surface

	reg         []byte
	oam         []byte
	vram		[]byte
	paletteBg   []byte
	paletteObj0 []byte
	paletteObj1 []byte

	tilemap [512][8][8]byte
	scanrow [161]byte
	objdata []ObjData

	modeClocks	uint16
	lineMode	byte
	curLine		byte
	curScan		uint16
	lcdon		bool
	bgtilebase	uint16
	bgmapbase	uint16
	objsize		bool
	objon		bool
	bgon		bool

	yscrl		byte
	xscrl		byte
	raster		byte
}

func NewGPU(cpu *CPU) *GPU {
	ret := new(GPU)

	ret.cpu = cpu
	ret.init()

	return ret
}

func (g *GPU) IsRunning() bool {
	return g.running
}

func (g *GPU) ReadOam(addr uint16) byte {
	return g.oam[addr]
}

func (g *GPU) ReadByte(addr uint16) byte {
	gaddr := addr - 0xFF40
	switch gaddr {
	case 0:
		var ret byte
		if g.lcdon {
			ret |= 0x80
		}
		if g.bgtilebase == 0x0000 {
			ret |= 0x10
		}
		if g.bgmapbase == 0x1c00 {
			ret |= 0x08
		}
		if g.objsize {
			ret |= 0x04
		}
		if g.objon {
			ret |= 0x02
		}
		if g.bgon {
			ret |= 0x01
		}
		return ret
	case 1:
		if g.curLine == g.raster {
			return 4 | g.lineMode
		} else {
			return g.lineMode
		}
	case 2:
		return g.yscrl
	case 3:
		return g.xscrl
	case 4:
		return g.curLine
	case 5:
		return g.raster
	default:
		return g.reg[gaddr]
	}
}

func (g *GPU) WriteVram(addr uint16, value byte) {
	fmt.Printf("Write VRAM 0x%x <- 0x%x\n", addr, value)

	g.vram[addr] = value
}

func (g *GPU) WriteOam(addr uint16, value byte) {
	g.oam[addr] = value
}

func (g *GPU) WriteByte(addr uint16, value byte) {
	// address range should be 0xFF40 - 0xFF7F
	gaddr := addr - 0xFF40
	g.reg[gaddr] = value

	switch gaddr {
	case 0:
		g.lcdon = value & 0x80 == 0x80
		if value & 0x10 == 0x10 {
			g.bgtilebase = 0x0000
		} else {
			g.bgtilebase = 0x0800
		}
		if value & 0x08 == 0x08 {
			g.bgmapbase = 0x1c00
		} else {
			g.bgmapbase = 0x1800
		}
		g.objsize = value & 0x04 == 0x04
		g.objon = value & 0x02 == 0x02
		g.bgon = value & 0x01 == 0x01
		break
	case 2:
		g.yscrl = value
		break
	case 3:
		g.xscrl = value
		break
	case 5:
		g.raster = value
		break
	case 6:
		var v byte
		for i := uint16(0); i < 160; i++ {
			v = g.cpu.ReadByte((uint16(value) << 8)+i)
			g.oam[i] = v
			g.UpdateOam(0xFE00+i, v)
		}
		break
	case 7:
		for i := uint16(0); i < 4; i++ {
			switch (value >> (i*2))&3 {
			case 0:
				g.paletteBg[i] = 255
				break
			case 1:
				g.paletteBg[i] = 192
				break
			case 2:
				g.paletteBg[i] = 96
				break
			case 3:
				g.paletteBg[i] = 0
				break
			}
		}
		break
	case 8:
		for i := uint16(0); i < 4; i++ {
			switch (value >> (i*2))&3 {
			case 0:
				g.paletteObj0[i] = 255
				break
			case 1:
				g.paletteObj0[i] = 192
				break
			case 2:
				g.paletteObj0[i] = 96
				break
			case 3:
				g.paletteObj0[i] = 0
				break
			}
		}
		break
	case 9:
		for i := uint16(0); i < 4; i++ {
			switch (value >> (i*2))&3 {
			case 0:
				g.paletteObj1[i] = 255
				break
			case 1:
				g.paletteObj1[i] = 192
				break
			case 2:
				g.paletteObj1[i] = 96
				break
			case 3:
				g.paletteObj1[i] = 0
				break
			}
		}
		break
	}
}

func (g *GPU) UpdateOam(addr uint16, data byte) {
	fmt.Println("Update oam")

	addr -= 0xfe00

	obj := addr >> 2

	if obj < 40 {
		switch addr & 3{
		case 0:
			g.objdata[obj].y = int16(data) - 16
			break
		case 1:
			g.objdata[obj].y = int16(data) - 8
			break
		case 2:
			if g.objsize {
				g.objdata[obj].tile = data & 0xfe
			} else {
				g.objdata[obj].tile = data
			}
			break
		case 3:
			if data & 0x10 == 0x10 {
				g.objdata[obj].palette = true
			} else {
				g.objdata[obj].palette = false
			}

			if data & 0x20 == 0x20 {
				g.objdata[obj].xflip = true
			} else {
				g.objdata[obj].xflip = false
			}

			if data & 0x40 == 0x40 {
				g.objdata[obj].yflip = true
			} else {
				g.objdata[obj].yflip = false
			}

			if data & 0x80 == 0x80 {
				g.objdata[obj].prio = true
			} else {
				g.objdata[obj].prio = false
			}

			break
		}
	}

	//g.objdatasorted = g.objdata
	//g.objdatasorted.Sort() todo sort

}

func (g *GPU) UpdateTile(addr uint16, data byte) {
	saddr := addr
	if addr & 1 == 1 {
		saddr--
		addr--
	}

	tile := (addr >> 4) & 511
	y := (addr >> 1) & 7
	var sx byte

	for i := byte(0); i < 8; i++ {
		sx = 1 << (7-i)

		var s byte
		if g.vram[saddr] & sx == sx {
			s |= 1
		}
		if g.vram[saddr+1] & sx == sx {
			s |= 2
		}

		g.tilemap[tile][y][i] = s
	}
}

func (g *GPU) CheckLine() {
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			g.running = false
		}
	}

	g.modeClocks += uint16(g.cpu.Register.M)

	switch g.lineMode {
	case 0:	// hblank
		if g.modeClocks >= 51 {
			if g.curLine == 143 {
				g.lineMode = 1
				// todo write image to window
				g.cpu.If |= 1
			} else {
				g.lineMode = 2
			}
			g.curLine++
			g.curScan += 640
			g.modeClocks = 0
		}
		break
	case 1: // vblank
		if g.modeClocks >= 114 {
			g.modeClocks = 0
			g.curLine++
			if g.curLine > 153 {
				g.curLine = 0
				g.curScan = 0
				g.lineMode = 2
			}
		}
		break
	case 2: // OAM-read
		if g.modeClocks >= 20 {
			g.modeClocks = 0
			g.lineMode = 3
		}
		break
	case 3: //VRAM-read
		if g.modeClocks >= 43 {
			g.modeClocks = 0
			g.lineMode = 0

			//fmt.Println("VRAM-Read")

			// read lcd-on flag (at 0xFF40)
			if g.lcdon {	// lcd is on
				if g.bgon {
					linebase := g.curScan
					mapbase := g.bgmapbase + ((((uint16(g.curLine)+uint16(g.yscrl))&255)>>3)<<5)
					y := (g.curLine+g.yscrl) & 7
					x := g.xscrl & 7
					t := (g.xscrl>>3) & 31

					//var pixel byte
					w := 160

					if g.bgtilebase != 0 {
						tile := uint16(g.vram[mapbase+uint16(t)])
						if tile < 128 {
							tile += 256
						}

						tilerow := g.tilemap[tile][y]

						for w > 0 {
							g.scanrow[160-x] = tilerow[x]
							color := g.paletteBg[tilerow[x]]

							//fmt.Printf("Write 0x%0x to %v\n", color, linebase+3)
							g.surface.Pixels()[linebase+0] = color
							g.surface.Pixels()[linebase+1] = color
							g.surface.Pixels()[linebase+2] = color
							x++
							if x == 8 {
								t = (t+1)&31
								x = 0
								tile = uint16(g.vram[mapbase+uint16(t)])
								if tile < 128 {
									tile += 256
								}
								tilerow = g.tilemap[tile][y]
							}
							linebase+=4

							w--
						}
					} else {
						tilerow := g.tilemap[g.vram[mapbase+uint16(t)]][y]

						for w > 0 {
							g.scanrow[160-x] = tilerow[x]
							//fmt.Printf("Write 0x%0x to %v\n", g.paletteBg[tilerow[x]], linebase+3)
							g.surface.Pixels()[linebase+0] = g.paletteBg[tilerow[x]]
							g.surface.Pixels()[linebase+1] = g.paletteBg[tilerow[x]]
							g.surface.Pixels()[linebase+2] = g.paletteBg[tilerow[x]]
							x++
							if x == 8 {
								t = (t+1)&31
								x = 0
								tilerow = g.tilemap[g.vram[mapbase+uint16(t)]][y]
							}
							linebase += 4

							w--
						}
					}
				}
				if g.objon {
					cnt := 0
					if g.objsize {
						for i := 0; i < 40; i++ {

						}
					} else {
						var tilerow [8]byte
						var obj ObjData
						var pal []byte
						//var x uint8
						linebase := g.curScan

						for i := 0; i < 40; i++ {
							obj = g.objdata[i]	// todo change to sorted
							if obj.y <= int16(g.curLine) && (obj.y+8) > int16(g.curLine) {
								if obj.yflip {
									tilerow = g.tilemap[obj.tile][7-(int16(g.curLine)-obj.y)]
								} else {
									tilerow = g.tilemap[obj.tile][int16(g.curLine)-obj.y]
								}

								if obj.palette {
									pal = g.paletteObj1
								} else {
									pal = g.paletteObj0
								}

								linebase = uint16((int16(g.curLine) * 160 + obj.x)*4)
								if obj.xflip {
									for x := int16(0); x < 8; x++ {
										if obj.x+x >= 0 && obj.x+x < 160 {
											if tilerow[7-x] != 0 && (obj.prio || g.scanrow[x] == 0) {
												//fmt.Printf("Write 0x%0x to %v\n", pal[tilerow[7-x]], linebase+3)
												g.surface.Pixels()[linebase+0] = pal[tilerow[7-x]]
												g.surface.Pixels()[linebase+1] = pal[tilerow[7-x]]
												g.surface.Pixels()[linebase+2] = pal[tilerow[7-x]]
											}
										}
										linebase += 4
									}
								} else {
									for x := int16(0); x < 8; x++ {
										if obj.x+x >= 0 && obj.x+x < 160 {
											if tilerow[x] != 0 && (obj.prio || g.scanrow[x] == 0) {
												//fmt.Printf("Write 0x%0x to %v\n", pal[tilerow[x]], linebase+3)
												g.surface.Pixels()[linebase+0] = pal[tilerow[x]]
												g.surface.Pixels()[linebase+1] = pal[tilerow[x]]
												g.surface.Pixels()[linebase+2] = pal[tilerow[x]]
											}
										}
										linebase += 4
									}
								}
								cnt++
								if cnt > 10 {
									break
								}
							}
						}
					}
				}
			}
		}
		break
	}

	g.window.UpdateSurface()
}

func (g *GPU) init() {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("GoGB", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 160, 144, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	g.reg = make([]byte, 64)
	g.oam = make([]byte, 160)
	g.vram = make([]byte, 8192)
	g.paletteBg = make([]byte, 4)
	g.paletteObj0 = make([]byte, 4)
	g.paletteObj1 = make([]byte, 4)
	g.objdata = make([]ObjData, 40)

	for i := byte(0); i < 40; i++ {
		g.objdata[i] = ObjData{x: -8, y: -16, tile: 0, palette: false, yflip: false, xflip: false, prio: false, num: i}
	}

	g.running = true
	g.window = window
	surface, err := g.window.GetSurface()
	if err != nil {
		panic(err)
	}

	g.surface = surface

	rect := sdl.Rect{0, 0, 160, 144}
	g.surface.FillRect(&rect, 0xffffffff)
}

