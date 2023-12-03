package nes

// CPU Memory Map
func CPUMemory_Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return nes_RAM[address%0x0800]
	case address < 0x4000:
		return PPU_readRegister(0x2000 + address%8)
	case address == 0x4014:
		return PPU_readRegister(address)
	case address == 0x4015:
		return 0
		//return this.console.APU.readRegister(address)
	case address == 0x4016:
		return Nes_Controller1.Read()
	case address == 0x4017:
		return Nes_Controller2.Read()
	case address < 0x6000:
		// TODO: I/O registers
	case address >= 0x6000:
		return Consol_Mapper.Read(address)
	default:
		log_Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0
}

func CPUMemory_Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		nes_RAM[address%0x0800] = value
	case address < 0x4000:
		PPU_writeRegister(0x2000+address%8, value)
	case address < 0x4014:
		//this.console.APU.writeRegister(address, value)
	case address == 0x4014:
		PPU_writeRegister(address, value)
	case address == 0x4015:
		//this.console.APU.writeRegister(address, value)
	case address == 0x4016:
		Nes_Controller1.Write(value)
		Nes_Controller2.Write(value)
	case address == 0x4017:
		//this.console.APU.writeRegister(address, value)
	case address < 0x6000:
		// TODO: I/O registers
	case address >= 0x6000:
		Consol_Mapper.Write(address, value)
	default:
		log_Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

// PPU Memory Map
func PPUMemory_Read(address uint16) byte {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		return Consol_Mapper.Read(address)
	case address < 0x3F00:
		mode := Cartridge_Mirror
		return PPU_nameTableData[MirrorAddress(mode, address)%2048]
	case address < 0x4000:
		return PPU_readPalette(address % 32)
	default:
		log_Fatalf("unhandled ppu memory read at address: 0x%04X", address)
	}
	return 0
}

func PPUMemory_Write(address uint16, value byte) {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		Consol_Mapper.Write(address, value)
	case address < 0x3F00:
		mode := Cartridge_Mirror
		PPU_nameTableData[MirrorAddress(mode, address)%2048] = value
	case address < 0x4000:
		PPU_writePalette(address%32, value)
	default:
		log_Fatalf("unhandled ppu memory write at address: 0x%04X", address)
	}
}

// Mirroring Modes

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorSingle0    = 2
	MirrorSingle1    = 3
	MirrorFour       = 4
)

var MirrorLookup = [...][4]uint16{
	{0, 0, 1, 1},
	{0, 1, 0, 1},
	{0, 0, 0, 0},
	{1, 1, 1, 1},
	{0, 1, 2, 3},
}

func MirrorAddress(mode byte, address uint16) uint16 {
	address = (address - 0x2000) % 0x1000
	table := address / 0x0400
	offset := address % 0x0400
	return 0x2000 + MirrorLookup[mode][table]*0x0400 + offset
}
