package nes

type Mapper40 struct {
	*Cartridge
	console *Console
	bank    int
	cycles  int
}

func NewMapper40(console *Console, cartridge *Cartridge) Mapper {
	return &Mapper40{cartridge, console, 0, 0}
}

func (this *Mapper40) Step() {
	if this.cycles < 0 {
		return
	}
	this.cycles++
	if this.cycles%(4096*3) == 0 {
		this.cycles = 0
		this.console.CPU.triggerIRQ()
	}
}

func (this *Mapper40) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return this.CHR[address]
	case address >= 0x6000 && address < 0x8000:
		return this.PRG[address-0x6000+0x2000*6]
	case address >= 0x8000 && address < 0xa000:
		return this.PRG[address-0x8000+0x2000*4]
	case address >= 0xa000 && address < 0xc000:
		return this.PRG[address-0xa000+0x2000*5]
	case address >= 0xc000 && address < 0xe000:
		return this.PRG[address-0xc000+0x2000*uint16(this.bank)]
	case address >= 0xe000:
		return this.PRG[address-0xe000+0x2000*7]
	default:
		log_Fatalf("unhandled mapper40 read at address: 0x%04X", address)
	}
	return 0
}

func (this *Mapper40) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		this.CHR[address] = value
	case address >= 0x8000 && address < 0xa000:
		this.cycles = -1
	case address >= 0xa000 && address < 0xc000:
		this.cycles = 0
	case address >= 0xe000:
		this.bank = int(value)
	default:
		log_Fatalf("unhandled mapper40 write at address: 0x%04X", address)
	}
}
