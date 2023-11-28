package nes

import (
	"errors"
	"log"
)

type color_RGBA struct {
	R, G, B, A uint8
}

func byte2str(x byte) string {
	return string(x)
}

func log_Fatalf(format string, x uint16) {
	log.Fatalf(format, x)
}

func errors_New(s string) error {
	return errors.New(s)
}
