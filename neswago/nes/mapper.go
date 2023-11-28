package nes

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
	Step()
}

func NewMapper(console *Console) (Mapper, error) {
	cartridge := console.Cartridge
	switch cartridge.Mapper {
	case 0:
		return NewMapper2(cartridge), nil
	case 1:
		return NewMapper1(cartridge), nil
	case 2:
		return NewMapper2(cartridge), nil
	case 3:
		return NewMapper3(cartridge), nil
	case 4:
		return NewMapper4(console, cartridge), nil
	case 7:
		return NewMapper7(cartridge), nil
	case 40:
		return NewMapper40(console, cartridge), nil
	case 225:
		return NewMapper225(cartridge), nil
	}
	err := errors_New("unsupported mapper: " + byte2str(cartridge.Mapper))
	return nil, err
}
