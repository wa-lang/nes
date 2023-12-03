package nes

const CPUFrequency = 1789773

// interrupt types
const (
	_ = iota
	interruptNone
	interruptNMI
	interruptIRQ
)

// addressing modes
const (
	_ = iota
	modeAbsolute
	modeAbsoluteX
	modeAbsoluteY
	modeAccumulator
	modeImmediate
	modeImplied
	modeIndexedIndirect
	modeIndirect
	modeIndirectIndexed
	modeRelative
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
)

// instructionModes indicates the addressing mode for each instruction
var instructionModes = [256]byte{
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	1, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 8, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
}

// instructionSizes indicates the size of each instruction in bytes
var instructionSizes = [256]byte{
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	3, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 0, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 0, 3, 0, 0,
	2, 2, 2, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
}

// instructionCycles indicates the number of cycles used by each instruction,
// not including conditional cycles
var instructionCycles = [256]byte{
	7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 6, 2, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

// instructionPageCycles indicates the number of cycles used by each
// instruction when a page is crossed
var instructionPageCycles = [256]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
}

// instructionNames indicates the name of each instruction
var instructionNames = [256]string{
	"BRK", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"PHP", "ORA", "ASL", "ANC", "NOP", "ORA", "ASL", "SLO",
	"BPL", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"CLC", "ORA", "NOP", "SLO", "NOP", "ORA", "ASL", "SLO",
	"JSR", "AND", "KIL", "RLA", "BIT", "AND", "ROL", "RLA",
	"PLP", "AND", "ROL", "ANC", "BIT", "AND", "ROL", "RLA",
	"BMI", "AND", "KIL", "RLA", "NOP", "AND", "ROL", "RLA",
	"SEC", "AND", "NOP", "RLA", "NOP", "AND", "ROL", "RLA",
	"RTI", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"PHA", "EOR", "LSR", "ALR", "JMP", "EOR", "LSR", "SRE",
	"BVC", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"CLI", "EOR", "NOP", "SRE", "NOP", "EOR", "LSR", "SRE",
	"RTS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"PLA", "ADC", "ROR", "ARR", "JMP", "ADC", "ROR", "RRA",
	"BVS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"SEI", "ADC", "NOP", "RRA", "NOP", "ADC", "ROR", "RRA",
	"NOP", "STA", "NOP", "SAX", "STY", "STA", "STX", "SAX",
	"DEY", "NOP", "TXA", "XAA", "STY", "STA", "STX", "SAX",
	"BCC", "STA", "KIL", "AHX", "STY", "STA", "STX", "SAX",
	"TYA", "STA", "TXS", "TAS", "SHY", "STA", "SHX", "AHX",
	"LDY", "LDA", "LDX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"TAY", "LDA", "TAX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"BCS", "LDA", "KIL", "LAX", "LDY", "LDA", "LDX", "LAX",
	"CLV", "LDA", "TSX", "LAS", "LDY", "LDA", "LDX", "LAX",
	"CPY", "CMP", "NOP", "DCP", "CPY", "CMP", "DEC", "DCP",
	"INY", "CMP", "DEX", "AXS", "CPY", "CMP", "DEC", "DCP",
	"BNE", "CMP", "KIL", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CLD", "CMP", "NOP", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CPX", "SBC", "NOP", "ISC", "CPX", "SBC", "INC", "ISC",
	"INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "ISC",
	"BEQ", "SBC", "KIL", "ISC", "NOP", "SBC", "INC", "ISC",
	"SED", "SBC", "NOP", "ISC", "NOP", "SBC", "INC", "ISC",
}

var CPU_Cycles uint64  // number of cycles
var CPU_PC uint16      // program counter
var CPU_SP byte        // stack pointer
var CPU_A byte         // accumulator
var CPU_X byte         // x register
var CPU_Y byte         // y register
var CPU_C byte         // carry flag
var CPU_Z byte         // zero flag
var CPU_I byte         // interrupt disable flag
var CPU_D byte         // decimal mode flag
var CPU_B byte         // break command flag
var CPU_U byte         // unused flag
var CPU_V byte         // overflow flag
var CPU_N byte         // negative flag
var CPU_interrupt byte // interrupt type to perform
var CPU_stall int      // number of cycles to stall
var CPU_table [256]func(*stepInfo)

//func CPU_Dump() {
//	f := file_log
//
//	fmt.Fprint(f, "CPU: { ", CPU_Cycles, " ", CPU_PC, " ", CPU_SP)
//	fmt.Fprint(f, "\n ", CPU_A, " ", CPU_X, " ", CPU_Y, " ", CPU_C, " ", CPU_Z, " ", CPU_I, " ", CPU_D, " ", CPU_B, " ", CPU_U, " ", CPU_V, " ", CPU_N)
//	fmt.Fprint(f, "\n ", CPU_interrupt, " ", CPU_stall, " }")
//	fmt.Fprint(f, "\n RAM:")
//	for _, v := range nes_RAM {
//		fmt.Fprint(f, v)
//		fmt.Fprint(f, " ")
//	}
//	fmt.Fprintln(f, "")
//}

func CPU_Init() {
	CPU_createTable()
	CPU_Reset()
}

// createTable builds a function table for each instruction
func CPU_createTable() {
	CPU_table = [256]func(*stepInfo){
		CPU_brk, CPU_ora, CPU_kil, CPU_slo, CPU_nop, CPU_ora, CPU_asl, CPU_slo,
		CPU_php, CPU_ora, CPU_asl, CPU_anc, CPU_nop, CPU_ora, CPU_asl, CPU_slo,
		CPU_bpl, CPU_ora, CPU_kil, CPU_slo, CPU_nop, CPU_ora, CPU_asl, CPU_slo,
		CPU_clc, CPU_ora, CPU_nop, CPU_slo, CPU_nop, CPU_ora, CPU_asl, CPU_slo,
		CPU_jsr, CPU_and, CPU_kil, CPU_rla, CPU_bit, CPU_and, CPU_rol, CPU_rla,
		CPU_plp, CPU_and, CPU_rol, CPU_anc, CPU_bit, CPU_and, CPU_rol, CPU_rla,
		CPU_bmi, CPU_and, CPU_kil, CPU_rla, CPU_nop, CPU_and, CPU_rol, CPU_rla,
		CPU_sec, CPU_and, CPU_nop, CPU_rla, CPU_nop, CPU_and, CPU_rol, CPU_rla,
		CPU_rti, CPU_eor, CPU_kil, CPU_sre, CPU_nop, CPU_eor, CPU_lsr, CPU_sre,
		CPU_pha, CPU_eor, CPU_lsr, CPU_alr, CPU_jmp, CPU_eor, CPU_lsr, CPU_sre,
		CPU_bvc, CPU_eor, CPU_kil, CPU_sre, CPU_nop, CPU_eor, CPU_lsr, CPU_sre,
		CPU_cli, CPU_eor, CPU_nop, CPU_sre, CPU_nop, CPU_eor, CPU_lsr, CPU_sre,
		CPU_rts, CPU_adc, CPU_kil, CPU_rra, CPU_nop, CPU_adc, CPU_ror, CPU_rra,
		CPU_pla, CPU_adc, CPU_ror, CPU_arr, CPU_jmp, CPU_adc, CPU_ror, CPU_rra,
		CPU_bvs, CPU_adc, CPU_kil, CPU_rra, CPU_nop, CPU_adc, CPU_ror, CPU_rra,
		CPU_sei, CPU_adc, CPU_nop, CPU_rra, CPU_nop, CPU_adc, CPU_ror, CPU_rra,
		CPU_nop, CPU_sta, CPU_nop, CPU_sax, CPU_sty, CPU_sta, CPU_stx, CPU_sax,
		CPU_dey, CPU_nop, CPU_txa, CPU_xaa, CPU_sty, CPU_sta, CPU_stx, CPU_sax,
		CPU_bcc, CPU_sta, CPU_kil, CPU_ahx, CPU_sty, CPU_sta, CPU_stx, CPU_sax,
		CPU_tya, CPU_sta, CPU_txs, CPU_tas, CPU_shy, CPU_sta, CPU_shx, CPU_ahx,
		CPU_ldy, CPU_lda, CPU_ldx, CPU_lax, CPU_ldy, CPU_lda, CPU_ldx, CPU_lax,
		CPU_tay, CPU_lda, CPU_tax, CPU_lax, CPU_ldy, CPU_lda, CPU_ldx, CPU_lax,
		CPU_bcs, CPU_lda, CPU_kil, CPU_lax, CPU_ldy, CPU_lda, CPU_ldx, CPU_lax,
		CPU_clv, CPU_lda, CPU_tsx, CPU_las, CPU_ldy, CPU_lda, CPU_ldx, CPU_lax,
		CPU_cpy, CPU_cmp, CPU_nop, CPU_dcp, CPU_cpy, CPU_cmp, CPU_dec, CPU_dcp,
		CPU_iny, CPU_cmp, CPU_dex, CPU_axs, CPU_cpy, CPU_cmp, CPU_dec, CPU_dcp,
		CPU_bne, CPU_cmp, CPU_kil, CPU_dcp, CPU_nop, CPU_cmp, CPU_dec, CPU_dcp,
		CPU_cld, CPU_cmp, CPU_nop, CPU_dcp, CPU_nop, CPU_cmp, CPU_dec, CPU_dcp,
		CPU_cpx, CPU_sbc, CPU_nop, CPU_isc, CPU_cpx, CPU_sbc, CPU_inc, CPU_isc,
		CPU_inx, CPU_sbc, CPU_nop, CPU_sbc, CPU_cpx, CPU_sbc, CPU_inc, CPU_isc,
		CPU_beq, CPU_sbc, CPU_kil, CPU_isc, CPU_nop, CPU_sbc, CPU_inc, CPU_isc,
		CPU_sed, CPU_sbc, CPU_nop, CPU_isc, CPU_nop, CPU_sbc, CPU_inc, CPU_isc,
	}
}

// Reset resets the CPU to its initial powerup state
func CPU_Reset() {
	CPU_PC = CPU_Read16(0xFFFC)
	CPU_SP = 0xFD
	CPU_SetFlags(0x24)
}

// pagesDiffer returns true if the two addresses reference different pages
func pagesDiffer(a, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// addBranchCycles adds a cycle for taking a branch and adds another cycle
// if the branch jumps to a new page
func CPU_addBranchCycles(info *stepInfo) {
	CPU_Cycles++
	if pagesDiffer(info.pc, info.address) {
		CPU_Cycles++
	}
}

func CPU_compare(a, b byte) {
	CPU_setZN(a - b)
	if a >= b {
		CPU_C = 1
	} else {
		CPU_C = 0
	}
}

// Read16 reads two bytes using Read to return a double-word value
func CPU_Read16(address uint16) uint16 {
	lo := uint16(CPUMemory_Read(address))
	hi := uint16(CPUMemory_Read(address + 1))
	return hi<<8 | lo
}

// read16bug emulates a 6502 bug that caused the low byte to wrap without
// incrementing the high byte
func CPU_read16bug(address uint16) uint16 {
	a := address
	b := (a & 0xFF00) | uint16(byte(a)+1)
	lo := CPUMemory_Read(a)
	hi := CPUMemory_Read(b)
	return uint16(hi)<<8 | uint16(lo)
}

// push pushes a byte onto the stack
func CPU_push(value byte) {
	CPUMemory_Write(0x100|uint16(CPU_SP), value)
	CPU_SP--
}

// pull pops a byte from the stack
func CPU_pull() byte {
	CPU_SP++
	return CPUMemory_Read(0x100 | uint16(CPU_SP))
}

// push16 pushes two bytes onto the stack
func CPU_push16(value uint16) {
	hi := byte(value >> 8)
	lo := byte(value & 0xFF)
	CPU_push(hi)
	CPU_push(lo)
}

// pull16 pops two bytes from the stack
func CPU_pull16() uint16 {
	lo := uint16(CPU_pull())
	hi := uint16(CPU_pull())
	return hi<<8 | lo
}

// Flags returns the processor status flags
func CPU_Flags() byte {
	var flags byte
	flags |= CPU_C << 0
	flags |= CPU_Z << 1
	flags |= CPU_I << 2
	flags |= CPU_D << 3
	flags |= CPU_B << 4
	flags |= CPU_U << 5
	flags |= CPU_V << 6
	flags |= CPU_N << 7
	return flags
}

// SetFlags sets the processor status flags
func CPU_SetFlags(flags byte) {
	CPU_C = (flags >> 0) & 1
	CPU_Z = (flags >> 1) & 1
	CPU_I = (flags >> 2) & 1
	CPU_D = (flags >> 3) & 1
	CPU_B = (flags >> 4) & 1
	CPU_U = (flags >> 5) & 1
	CPU_V = (flags >> 6) & 1
	CPU_N = (flags >> 7) & 1
}

// setZ sets the zero flag if the argument is zero
func CPU_setZ(value byte) {
	if value == 0 {
		CPU_Z = 1
	} else {
		CPU_Z = 0
	}
}

// setN sets the negative flag if the argument is negative (high bit is set)
func CPU_setN(value byte) {
	if value&0x80 != 0 {
		CPU_N = 1
	} else {
		CPU_N = 0
	}
}

// setZN sets the zero flag and the negative flag
func CPU_setZN(value byte) {
	CPU_setZ(value)
	CPU_setN(value)
}

// triggerNMI causes a non-maskable interrupt to occur on the next cycle
func CPU_triggerNMI() {
	CPU_interrupt = interruptNMI
}

// triggerIRQ causes an IRQ interrupt to occur on the next cycle
func CPU_triggerIRQ() {
	if CPU_I == 0 {
		CPU_interrupt = interruptIRQ
	}
}

// stepInfo contains information that the instruction functions use
type stepInfo struct {
	address uint16
	pc      uint16
	mode    byte
}

// Step executes a single CPU instruction
func CPU_Step() int {
	if CPU_stall > 0 {
		CPU_stall--
		return 1
	}

	cycles := CPU_Cycles

	switch CPU_interrupt {
	case interruptNMI:
		CPU_nmi()
	case interruptIRQ:
		CPU_irq()
	}
	CPU_interrupt = interruptNone

	opcode := CPUMemory_Read(CPU_PC)
	mode := instructionModes[opcode]

	var address uint16
	var pageCrossed bool
	switch mode {
	case modeAbsolute:
		address = CPU_Read16(CPU_PC + 1)
	case modeAbsoluteX:
		address = CPU_Read16(CPU_PC+1) + uint16(CPU_X)
		pageCrossed = pagesDiffer(address-uint16(CPU_X), address)
	case modeAbsoluteY:
		address = CPU_Read16(CPU_PC+1) + uint16(CPU_Y)
		pageCrossed = pagesDiffer(address-uint16(CPU_Y), address)
	case modeAccumulator:
		address = 0
	case modeImmediate:
		address = CPU_PC + 1
	case modeImplied:
		address = 0
	case modeIndexedIndirect:
		address = CPU_read16bug(uint16(CPUMemory_Read(CPU_PC+1) + CPU_X))
	case modeIndirect:
		address = CPU_read16bug(CPU_Read16(CPU_PC + 1))
	case modeIndirectIndexed:
		address = CPU_read16bug(uint16(CPUMemory_Read(CPU_PC+1))) + uint16(CPU_Y)
		pageCrossed = pagesDiffer(address-uint16(CPU_Y), address)
	case modeRelative:
		offset := uint16(CPUMemory_Read(CPU_PC + 1))
		if offset < 0x80 {
			address = CPU_PC + 2 + offset
		} else {
			address = CPU_PC + 2 + offset - 0x100
		}
	case modeZeroPage:
		address = uint16(CPUMemory_Read(CPU_PC + 1))
	case modeZeroPageX:
		address = uint16(CPUMemory_Read(CPU_PC+1)+CPU_X) & 0xff
	case modeZeroPageY:
		address = uint16(CPUMemory_Read(CPU_PC+1)+CPU_Y) & 0xff
	}

	CPU_PC += uint16(instructionSizes[opcode])
	CPU_Cycles += uint64(instructionCycles[opcode])
	if pageCrossed {
		CPU_Cycles += uint64(instructionPageCycles[opcode])
	}
	info := &stepInfo{address, CPU_PC, mode}

	CPU_table[opcode](info)

	return int(CPU_Cycles - cycles)
}

// NMI - Non-Maskable Interrupt
func CPU_nmi() {
	CPU_push16(CPU_PC)
	CPU_php(nil)
	CPU_PC = CPU_Read16(0xFFFA)
	CPU_I = 1
	CPU_Cycles += 7
}

// IRQ - IRQ Interrupt
func CPU_irq() {
	CPU_push16(CPU_PC)
	CPU_php(nil)
	CPU_PC = CPU_Read16(0xFFFE)
	CPU_I = 1
	CPU_Cycles += 7
}

// ADC - Add with Carry
func CPU_adc(info *stepInfo) {
	a := CPU_A
	b := CPUMemory_Read(info.address)
	c := CPU_C
	CPU_A = a + b + c
	CPU_setZN(CPU_A)
	if int(a)+int(b)+int(c) > 0xFF {
		CPU_C = 1
	} else {
		CPU_C = 0
	}
	if (a^b)&0x80 == 0 && (a^CPU_A)&0x80 != 0 {
		CPU_V = 1
	} else {
		CPU_V = 0
	}
}

// AND - Logical AND
func CPU_and(info *stepInfo) {
	CPU_A = CPU_A & CPUMemory_Read(info.address)
	CPU_setZN(CPU_A)
}

// ASL - Arithmetic Shift Left
func CPU_asl(info *stepInfo) {
	if info.mode == modeAccumulator {
		CPU_C = (CPU_A >> 7) & 1
		CPU_A <<= 1
		CPU_setZN(CPU_A)
	} else {
		value := CPUMemory_Read(info.address)
		CPU_C = (value >> 7) & 1
		value <<= 1
		CPUMemory_Write(info.address, value)
		CPU_setZN(value)
	}
}

// BCC - Branch if Carry Clear
func CPU_bcc(info *stepInfo) {
	if CPU_C == 0 {
		CPU_PC = info.address
		CPU_addBranchCycles(info)
	}
}

// BCS - Branch if Carry Set
func CPU_bcs(info *stepInfo) {
	if CPU_C != 0 {
		CPU_PC = info.address
		CPU_addBranchCycles(info)
	}
}

// BEQ - Branch if Equal
func CPU_beq(info *stepInfo) {
	if CPU_Z != 0 {
		CPU_PC = info.address
		CPU_addBranchCycles(info)
	}
}

// BIT - Bit Test
func CPU_bit(info *stepInfo) {
	value := CPUMemory_Read(info.address)
	CPU_V = (value >> 6) & 1
	CPU_setZ(value & CPU_A)
	CPU_setN(value)
}

// BMI - Branch if Minus
func CPU_bmi(info *stepInfo) {
	if CPU_N != 0 {
		CPU_PC = info.address
		CPU_addBranchCycles(info)
	}
}

// BNE - Branch if Not Equal
func CPU_bne(info *stepInfo) {
	if CPU_Z == 0 {
		CPU_PC = info.address
		CPU_addBranchCycles(info)
	}
}

// BPL - Branch if Positive
func CPU_bpl(info *stepInfo) {
	if CPU_N == 0 {
		CPU_PC = info.address
		CPU_addBranchCycles(info)
	}
}

// BRK - Force Interrupt
func CPU_brk(info *stepInfo) {
	CPU_push16(CPU_PC)
	CPU_php(info)
	CPU_sei(info)
	CPU_PC = CPU_Read16(0xFFFE)
}

// BVC - Branch if Overflow Clear
func CPU_bvc(info *stepInfo) {
	if CPU_V == 0 {
		CPU_PC = info.address
		CPU_addBranchCycles(info)
	}
}

// BVS - Branch if Overflow Set
func CPU_bvs(info *stepInfo) {
	if CPU_V != 0 {
		CPU_PC = info.address
		CPU_addBranchCycles(info)
	}
}

// CLC - Clear Carry Flag
func CPU_clc(info *stepInfo) {
	CPU_C = 0
}

// CLD - Clear Decimal Mode
func CPU_cld(info *stepInfo) {
	CPU_D = 0
}

// CLI - Clear Interrupt Disable
func CPU_cli(info *stepInfo) {
	CPU_I = 0
}

// CLV - Clear Overflow Flag
func CPU_clv(info *stepInfo) {
	CPU_V = 0
}

// CMP - Compare
func CPU_cmp(info *stepInfo) {
	value := CPUMemory_Read(info.address)
	CPU_compare(CPU_A, value)
}

// CPX - Compare X Register
func CPU_cpx(info *stepInfo) {
	value := CPUMemory_Read(info.address)
	CPU_compare(CPU_X, value)
}

// CPY - Compare Y Register
func CPU_cpy(info *stepInfo) {
	value := CPUMemory_Read(info.address)
	CPU_compare(CPU_Y, value)
}

// DEC - Decrement Memory
func CPU_dec(info *stepInfo) {
	value := CPUMemory_Read(info.address) - 1
	CPUMemory_Write(info.address, value)
	CPU_setZN(value)
}

// DEX - Decrement X Register
func CPU_dex(info *stepInfo) {
	CPU_X--
	CPU_setZN(CPU_X)
}

// DEY - Decrement Y Register
func CPU_dey(info *stepInfo) {
	CPU_Y--
	CPU_setZN(CPU_Y)
}

// EOR - Exclusive OR
func CPU_eor(info *stepInfo) {
	CPU_A = CPU_A ^ CPUMemory_Read(info.address)
	CPU_setZN(CPU_A)
}

// INC - Increment Memory
func CPU_inc(info *stepInfo) {
	value := CPUMemory_Read(info.address) + 1
	CPUMemory_Write(info.address, value)
	CPU_setZN(value)
}

// INX - Increment X Register
func CPU_inx(info *stepInfo) {
	CPU_X++
	CPU_setZN(CPU_X)
}

// INY - Increment Y Register
func CPU_iny(info *stepInfo) {
	CPU_Y++
	CPU_setZN(CPU_Y)
}

// JMP - Jump
func CPU_jmp(info *stepInfo) {
	CPU_PC = info.address
}

// JSR - Jump to Subroutine
func CPU_jsr(info *stepInfo) {
	CPU_push16(CPU_PC - 1)
	CPU_PC = info.address
}

// LDA - Load Accumulator
func CPU_lda(info *stepInfo) {
	CPU_A = CPUMemory_Read(info.address)
	CPU_setZN(CPU_A)
}

// LDX - Load X Register
func CPU_ldx(info *stepInfo) {
	CPU_X = CPUMemory_Read(info.address)
	CPU_setZN(CPU_X)
}

// LDY - Load Y Register
func CPU_ldy(info *stepInfo) {
	CPU_Y = CPUMemory_Read(info.address)
	CPU_setZN(CPU_Y)
}

// LSR - Logical Shift Right
func CPU_lsr(info *stepInfo) {
	if info.mode == modeAccumulator {
		CPU_C = CPU_A & 1
		CPU_A >>= 1
		CPU_setZN(CPU_A)
	} else {
		value := CPUMemory_Read(info.address)
		CPU_C = value & 1
		value >>= 1
		CPUMemory_Write(info.address, value)
		CPU_setZN(value)
	}
}

// NOP - No Operation
func CPU_nop(info *stepInfo) {
}

// ORA - Logical Inclusive OR
func CPU_ora(info *stepInfo) {
	CPU_A = CPU_A | CPUMemory_Read(info.address)
	CPU_setZN(CPU_A)
}

// PHA - Push Accumulator
func CPU_pha(info *stepInfo) {
	CPU_push(CPU_A)
}

// PHP - Push Processor Status
func CPU_php(info *stepInfo) {
	CPU_push(CPU_Flags() | 0x10)
}

// PLA - Pull Accumulator
func CPU_pla(info *stepInfo) {
	CPU_A = CPU_pull()
	CPU_setZN(CPU_A)
}

// PLP - Pull Processor Status
func CPU_plp(info *stepInfo) {
	CPU_SetFlags(CPU_pull()&0xEF | 0x20)
}

// ROL - Rotate Left
func CPU_rol(info *stepInfo) {
	if info.mode == modeAccumulator {
		c := CPU_C
		CPU_C = (CPU_A >> 7) & 1
		CPU_A = (CPU_A << 1) | c
		CPU_setZN(CPU_A)
	} else {
		c := CPU_C
		value := CPUMemory_Read(info.address)
		CPU_C = (value >> 7) & 1
		value = (value << 1) | c
		CPUMemory_Write(info.address, value)
		CPU_setZN(value)
	}
}

// ROR - Rotate Right
func CPU_ror(info *stepInfo) {
	if info.mode == modeAccumulator {
		c := CPU_C
		CPU_C = CPU_A & 1
		CPU_A = (CPU_A >> 1) | (c << 7)
		CPU_setZN(CPU_A)
	} else {
		c := CPU_C
		value := CPUMemory_Read(info.address)
		CPU_C = value & 1
		value = (value >> 1) | (c << 7)
		CPUMemory_Write(info.address, value)
		CPU_setZN(value)
	}
}

// RTI - Return from Interrupt
func CPU_rti(info *stepInfo) {
	CPU_SetFlags(CPU_pull()&0xEF | 0x20)
	CPU_PC = CPU_pull16()
}

// RTS - Return from Subroutine
func CPU_rts(info *stepInfo) {
	CPU_PC = CPU_pull16() + 1
}

// SBC - Subtract with Carry
func CPU_sbc(info *stepInfo) {
	a := CPU_A
	b := CPUMemory_Read(info.address)
	c := CPU_C
	CPU_A = a - b - (1 - c)
	CPU_setZN(CPU_A)
	if int(a)-int(b)-int(1-c) >= 0 {
		CPU_C = 1
	} else {
		CPU_C = 0
	}
	if (a^b)&0x80 != 0 && (a^CPU_A)&0x80 != 0 {
		CPU_V = 1
	} else {
		CPU_V = 0
	}
}

// SEC - Set Carry Flag
func CPU_sec(info *stepInfo) {
	CPU_C = 1
}

// SED - Set Decimal Flag
func CPU_sed(info *stepInfo) {
	CPU_D = 1
}

// SEI - Set Interrupt Disable
func CPU_sei(info *stepInfo) {
	CPU_I = 1
}

// STA - Store Accumulator
func CPU_sta(info *stepInfo) {
	CPUMemory_Write(info.address, CPU_A)
}

// STX - Store X Register
func CPU_stx(info *stepInfo) {
	CPUMemory_Write(info.address, CPU_X)
}

// STY - Store Y Register
func CPU_sty(info *stepInfo) {
	CPUMemory_Write(info.address, CPU_Y)
}

// TAX - Transfer Accumulator to X
func CPU_tax(info *stepInfo) {
	CPU_X = CPU_A
	CPU_setZN(CPU_X)
}

// TAY - Transfer Accumulator to Y
func CPU_tay(info *stepInfo) {
	CPU_Y = CPU_A
	CPU_setZN(CPU_Y)
}

// TSX - Transfer Stack Pointer to X
func CPU_tsx(info *stepInfo) {
	CPU_X = CPU_SP
	CPU_setZN(CPU_X)
}

// TXA - Transfer X to Accumulator
func CPU_txa(info *stepInfo) {
	CPU_A = CPU_X
	CPU_setZN(CPU_A)
}

// TXS - Transfer X to Stack Pointer
func CPU_txs(info *stepInfo) {
	CPU_SP = CPU_X
}

// TYA - Transfer Y to Accumulator
func CPU_tya(info *stepInfo) {
	CPU_A = CPU_Y
	CPU_setZN(CPU_A)
}

// illegal opcodes below

func CPU_ahx(info *stepInfo) {
}

func CPU_alr(info *stepInfo) {
}

func CPU_anc(info *stepInfo) {
}

func CPU_arr(info *stepInfo) {
}

func CPU_axs(info *stepInfo) {
}

func CPU_dcp(info *stepInfo) {
}

func CPU_isc(info *stepInfo) {
}

func CPU_kil(info *stepInfo) {
}

func CPU_las(info *stepInfo) {
}

func CPU_lax(info *stepInfo) {
}

func CPU_rla(info *stepInfo) {
}

func CPU_rra(info *stepInfo) {
}

func CPU_sax(info *stepInfo) {
}

func CPU_shx(info *stepInfo) {
}

func CPU_shy(info *stepInfo) {
}

func CPU_slo(info *stepInfo) {
}

func CPU_sre(info *stepInfo) {
}

func CPU_tas(info *stepInfo) {
}

func CPU_xaa(info *stepInfo) {
}
