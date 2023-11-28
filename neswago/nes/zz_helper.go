package nes

import (
	"errors"
	"log"
	"unsafe"
)

func byte2str(x byte) string {
	return string(x)
}

func raw(hdr []iNESFileHeader) []byte {
	return (*(*[1 << 20]byte)(unsafe.Pointer(&hdr[0])))[:unsafe.Sizeof(hdr[0])]
}

func log_Fatalf(format string, x uint16) {
	log.Fatalf(format, x)
}

func errors_New(s string) error {
	return errors.New(s)
}
