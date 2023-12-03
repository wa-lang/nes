package nes

type Mapper7 struct {
	prgBank int
}

func NewMapper7() Mapper {
	return &Mapper7{0}
}

func (this *Mapper7) Step() {
}

func (this *Mapper7) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return Cartridge_CHR[address]
	case address >= 0x8000:
		index := this.prgBank*0x8000 + int(address-0x8000)
		return Cartridge_PRG[index]
	case address >= 0x6000:
		index := int(address) - 0x6000
		return Cartridge_SRAM[index]
	default:
		log_Fatalf("unhandled mapper7 read at address: 0x%04X", address)
	}
	return 0
}

func (this *Mapper7) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		Cartridge_CHR[address] = value
	case address >= 0x8000:
		this.prgBank = int(value & 7)
		switch value & 0x10 {
		case 0x00:
			Cartridge_Mirror = MirrorSingle0
		case 0x10:
			Cartridge_Mirror = MirrorSingle1
		}
	case address >= 0x6000:
		index := int(address) - 0x6000
		Cartridge_SRAM[index] = value
	default:
		log_Fatalf("unhandled mapper7 write at address: 0x%04X", address)
	}
}
