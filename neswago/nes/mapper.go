package nes

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
	Step()
}

func NewMapper() (Mapper, error) {
	switch Cartridge_Mapper {
	case 0:
		return NewMapper2(), nil
	case 1:
		return NewMapper1(), nil
	case 2:
		return NewMapper2(), nil
	case 3:
		return NewMapper3(), nil
	case 4:
		return NewMapper4(), nil
	case 7:
		return NewMapper7(), nil
	case 40:
		return NewMapper40(), nil
	case 225:
		return NewMapper225(), nil
	}
	err := errors_New("unsupported mapper: " + byte2str(Cartridge_Mapper))
	return nil, err
}
