package nes

const iNESFileMagic = 0x1a53454e

type iNESFileHeader struct {
	Magic    uint32  // iNES magic number
	NumPRG   byte    // number of PRG-ROM banks (16KB each)
	NumCHR   byte    // number of CHR-ROM banks (8KB each)
	Control1 byte    // control bits
	Control2 byte    // control bits
	NumRAM   byte    // PRG-RAM size (x 8KB)
	_        [7]byte // unused padding
}

func readNESFileHeader(buf []byte, hdr *iNESFileHeader) (int, error) {
	arr := make([]iNESFileHeader, 1, 1)
	b := raw(arr)
	n := len(b)

	if n > len(buf) {
		return 0, errors_New("EOF")
	}

	copy(b, buf)
	*hdr = arr[0]
	return n, nil
}

// LoadNESFile reads an iNES file (.nes) and returns a Cartridge on success.
// http://wiki.nesdev.com/w/index.php/INES
// http://nesdev.com/NESDoc.pdf (page 28)
func LoadNESFile(romBytes []byte) (*Cartridge, error) {
	file := romBytes

	// read file header
	header := iNESFileHeader{}
	n, err := readNESFileHeader(file, &header)
	if err != nil {
		return nil, err
	}
	file = file[n:]

	// verify header magic number
	if header.Magic != iNESFileMagic {
		return nil, errors_New("invalid .nes file")
	}

	// mapper type
	mapper1 := header.Control1 >> 4
	mapper2 := header.Control2 >> 4
	mapper := mapper1 | mapper2<<4

	// mirroring type
	mirror1 := header.Control1 & 1
	mirror2 := (header.Control1 >> 3) & 1
	mirror := mirror1 | mirror2<<1

	// battery-backed RAM
	battery := (header.Control1 >> 1) & 1

	// read trainer if present (unused)
	if header.Control1&4 == 4 {
		trainer := make([]byte, 512)
		if len(file) < len(trainer) {
			return nil, errors_New("EOF")
		}
		copy(trainer, file)
		file = file[len(trainer):]
	}

	// read prg-rom bank(s)
	prg := make([]byte, int(header.NumPRG)*16384)
	if len(file) < len(prg) {
		return nil, errors_New("EOF")
	}
	copy(prg, file)
	file = file[len(prg):]

	// read chr-rom bank(s)
	chr := make([]byte, int(header.NumCHR)*8192)
	if len(file) < len(chr) {
		return nil, errors_New("EOF")
	}
	copy(chr, file)
	file = file[len(chr):]

	// provide chr-rom/ram if not in file
	if header.NumCHR == 0 {
		chr = make([]byte, 8192)
	}

	// success
	return NewCartridge(prg, chr, mapper, mirror, battery), nil
}
