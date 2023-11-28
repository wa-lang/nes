package nes

type Mapper3 struct {
	*Cartridge
	chrBank  int
	prgBank1 int
	prgBank2 int
}

func NewMapper3(cartridge *Cartridge) Mapper {
	prgBanks := len(cartridge.PRG) / 0x4000
	return &Mapper3{cartridge, 0, 0, prgBanks - 1}
}

func (this *Mapper3) Step() {
}

func (this *Mapper3) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		index := this.chrBank*0x2000 + int(address)
		return this.CHR[index]
	case address >= 0xC000:
		index := this.prgBank2*0x4000 + int(address-0xC000)
		return this.PRG[index]
	case address >= 0x8000:
		index := this.prgBank1*0x4000 + int(address-0x8000)
		return this.PRG[index]
	case address >= 0x6000:
		index := int(address) - 0x6000
		return this.SRAM[index]
	default:
		log_Fatalf("unhandled mapper3 read at address: 0x%04X", address)
	}
	return 0
}

func (this *Mapper3) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		index := this.chrBank*0x2000 + int(address)
		this.CHR[index] = value
	case address >= 0x8000:
		this.chrBank = int(value & 3)
	case address >= 0x6000:
		index := int(address) - 0x6000
		this.SRAM[index] = value
	default:
		log_Fatalf("unhandled mapper3 write at address: 0x%04X", address)
	}
}
