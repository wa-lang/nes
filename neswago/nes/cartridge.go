package nes

var Cartridge_PRG []byte   // PRG-ROM banks
var Cartridge_CHR []byte   // CHR-ROM banks
var Cartridge_SRAM []byte  // Save RAM
var Cartridge_Mapper byte  // mapper type
var Cartridge_Mirror byte  // mirroring mode
var Cartridge_Battery byte // battery present

func InitCartridge(prg, chr []byte, mapper, mirror, battery byte) {
	Cartridge_PRG = prg
	Cartridge_CHR = chr
	Cartridge_SRAM = make([]byte, 0x2000)
	Cartridge_Mapper = mapper
	Cartridge_Mirror = mirror
	Cartridge_Battery = battery
}
