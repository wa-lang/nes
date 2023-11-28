package nes

import "log"

func byte2str(x byte) string {
	return string(x)
}

func log_Fatalf(format string, x uint16) {
	log.Fatalf(format, x)
}
