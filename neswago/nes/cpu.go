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

type CPU struct {
	*CPUMemory        // memory
	Cycles     uint64 // number of cycles
	PC         uint16 // program counter
	SP         byte   // stack pointer
	A          byte   // accumulator
	X          byte   // x register
	Y          byte   // y register
	C          byte   // carry flag
	Z          byte   // zero flag
	I          byte   // interrupt disable flag
	D          byte   // decimal mode flag
	B          byte   // break command flag
	U          byte   // unused flag
	V          byte   // overflow flag
	N          byte   // negative flag
	interrupt  byte   // interrupt type to perform
	stall      int    // number of cycles to stall
	table      [256]func(*stepInfo)
}

func (this *CPU) Dump() {
	print("{ ", this.Cycles, " ", this.PC, " ", this.SP, " ", this.A, " ", this.X, " ", this.Y, " ", this.C, " ", this.Z, " ", this.I, " ", this.D, " ", this.B, " ", this.U, " ", this.V, " ", this.N, " ", this.interrupt, " ", this.stall, " }")
	for _, v := range this.CPUMemory.console.RAM {
		print(v)
		print(" ")
	}
	println("")
}

func NewCPU(console *Console) *CPU {
	this := &CPU{CPUMemory: NewCPUMemory(console)}
	this.createTable()
	this.Reset()
	return this
}

// createTable builds a function table for each instruction
func (this *CPU) createTable() {
	this.table = [256]func(*stepInfo){
		this.brk, this.ora, this.kil, this.slo, this.nop, this.ora, this.asl, this.slo,
		this.php, this.ora, this.asl, this.anc, this.nop, this.ora, this.asl, this.slo,
		this.bpl, this.ora, this.kil, this.slo, this.nop, this.ora, this.asl, this.slo,
		this.clc, this.ora, this.nop, this.slo, this.nop, this.ora, this.asl, this.slo,
		this.jsr, this.and, this.kil, this.rla, this.bit, this.and, this.rol, this.rla,
		this.plp, this.and, this.rol, this.anc, this.bit, this.and, this.rol, this.rla,
		this.bmi, this.and, this.kil, this.rla, this.nop, this.and, this.rol, this.rla,
		this.sec, this.and, this.nop, this.rla, this.nop, this.and, this.rol, this.rla,
		this.rti, this.eor, this.kil, this.sre, this.nop, this.eor, this.lsr, this.sre,
		this.pha, this.eor, this.lsr, this.alr, this.jmp, this.eor, this.lsr, this.sre,
		this.bvc, this.eor, this.kil, this.sre, this.nop, this.eor, this.lsr, this.sre,
		this.cli, this.eor, this.nop, this.sre, this.nop, this.eor, this.lsr, this.sre,
		this.rts, this.adc, this.kil, this.rra, this.nop, this.adc, this.ror, this.rra,
		this.pla, this.adc, this.ror, this.arr, this.jmp, this.adc, this.ror, this.rra,
		this.bvs, this.adc, this.kil, this.rra, this.nop, this.adc, this.ror, this.rra,
		this.sei, this.adc, this.nop, this.rra, this.nop, this.adc, this.ror, this.rra,
		this.nop, this.sta, this.nop, this.sax, this.sty, this.sta, this.stx, this.sax,
		this.dey, this.nop, this.txa, this.xaa, this.sty, this.sta, this.stx, this.sax,
		this.bcc, this.sta, this.kil, this.ahx, this.sty, this.sta, this.stx, this.sax,
		this.tya, this.sta, this.txs, this.tas, this.shy, this.sta, this.shx, this.ahx,
		this.ldy, this.lda, this.ldx, this.lax, this.ldy, this.lda, this.ldx, this.lax,
		this.tay, this.lda, this.tax, this.lax, this.ldy, this.lda, this.ldx, this.lax,
		this.bcs, this.lda, this.kil, this.lax, this.ldy, this.lda, this.ldx, this.lax,
		this.clv, this.lda, this.tsx, this.las, this.ldy, this.lda, this.ldx, this.lax,
		this.cpy, this.cmp, this.nop, this.dcp, this.cpy, this.cmp, this.dec, this.dcp,
		this.iny, this.cmp, this.dex, this.axs, this.cpy, this.cmp, this.dec, this.dcp,
		this.bne, this.cmp, this.kil, this.dcp, this.nop, this.cmp, this.dec, this.dcp,
		this.cld, this.cmp, this.nop, this.dcp, this.nop, this.cmp, this.dec, this.dcp,
		this.cpx, this.sbc, this.nop, this.isc, this.cpx, this.sbc, this.inc, this.isc,
		this.inx, this.sbc, this.nop, this.sbc, this.cpx, this.sbc, this.inc, this.isc,
		this.beq, this.sbc, this.kil, this.isc, this.nop, this.sbc, this.inc, this.isc,
		this.sed, this.sbc, this.nop, this.isc, this.nop, this.sbc, this.inc, this.isc,
	}
}

// Reset resets the CPU to its initial powerup state
func (this *CPU) Reset() {
	this.PC = this.Read16(0xFFFC)
	this.SP = 0xFD
	this.SetFlags(0x24)
}

// pagesDiffer returns true if the two addresses reference different pages
func pagesDiffer(a, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// addBranchCycles adds a cycle for taking a branch and adds another cycle
// if the branch jumps to a new page
func (this *CPU) addBranchCycles(info *stepInfo) {
	this.Cycles++
	if pagesDiffer(info.pc, info.address) {
		this.Cycles++
	}
}

func (this *CPU) compare(a, b byte) {
	this.setZN(a - b)
	if a >= b {
		this.C = 1
	} else {
		this.C = 0
	}
}

// Read16 reads two bytes using Read to return a double-word value
func (this *CPU) Read16(address uint16) uint16 {
	lo := uint16(this.Read(address))
	hi := uint16(this.Read(address + 1))
	return hi<<8 | lo
}

// read16bug emulates a 6502 bug that caused the low byte to wrap without
// incrementing the high byte
func (this *CPU) read16bug(address uint16) uint16 {
	a := address
	b := (a & 0xFF00) | uint16(byte(a)+1)
	lo := this.Read(a)
	hi := this.Read(b)
	return uint16(hi)<<8 | uint16(lo)
}

// push pushes a byte onto the stack
func (this *CPU) push(value byte) {
	this.Write(0x100|uint16(this.SP), value)
	this.SP--
}

// pull pops a byte from the stack
func (this *CPU) pull() byte {
	this.SP++
	return this.Read(0x100 | uint16(this.SP))
}

// push16 pushes two bytes onto the stack
func (this *CPU) push16(value uint16) {
	hi := byte(value >> 8)
	lo := byte(value & 0xFF)
	this.push(hi)
	this.push(lo)
}

// pull16 pops two bytes from the stack
func (this *CPU) pull16() uint16 {
	lo := uint16(this.pull())
	hi := uint16(this.pull())
	return hi<<8 | lo
}

// Flags returns the processor status flags
func (this *CPU) Flags() byte {
	var flags byte
	flags |= this.C << 0
	flags |= this.Z << 1
	flags |= this.I << 2
	flags |= this.D << 3
	flags |= this.B << 4
	flags |= this.U << 5
	flags |= this.V << 6
	flags |= this.N << 7
	return flags
}

// SetFlags sets the processor status flags
func (this *CPU) SetFlags(flags byte) {
	this.C = (flags >> 0) & 1
	this.Z = (flags >> 1) & 1
	this.I = (flags >> 2) & 1
	this.D = (flags >> 3) & 1
	this.B = (flags >> 4) & 1
	this.U = (flags >> 5) & 1
	this.V = (flags >> 6) & 1
	this.N = (flags >> 7) & 1
}

// setZ sets the zero flag if the argument is zero
func (this *CPU) setZ(value byte) {
	if value == 0 {
		this.Z = 1
	} else {
		this.Z = 0
	}
}

// setN sets the negative flag if the argument is negative (high bit is set)
func (this *CPU) setN(value byte) {
	if value&0x80 != 0 {
		this.N = 1
	} else {
		this.N = 0
	}
}

// setZN sets the zero flag and the negative flag
func (this *CPU) setZN(value byte) {
	this.setZ(value)
	this.setN(value)
}

// triggerNMI causes a non-maskable interrupt to occur on the next cycle
func (this *CPU) triggerNMI() {
	this.interrupt = interruptNMI
}

// triggerIRQ causes an IRQ interrupt to occur on the next cycle
func (this *CPU) triggerIRQ() {
	if this.I == 0 {
		this.interrupt = interruptIRQ
	}
}

// stepInfo contains information that the instruction functions use
type stepInfo struct {
	address uint16
	pc      uint16
	mode    byte
}

var Halt bool

// Step executes a single CPU instruction
func (this *CPU) Step() int {
	if this.stall > 0 {
		this.stall--
		return 1
	}

	cycles := this.Cycles

	switch this.interrupt {
	case interruptNMI:
		this.nmi()
	case interruptIRQ:
		this.irq()
	}
	this.interrupt = interruptNone

	opcode := this.Read(this.PC)
	mode := instructionModes[opcode]

	var address uint16
	var pageCrossed bool
	switch mode {
	case modeAbsolute:
		address = this.Read16(this.PC + 1)
	case modeAbsoluteX:
		address = this.Read16(this.PC+1) + uint16(this.X)
		pageCrossed = pagesDiffer(address-uint16(this.X), address)
	case modeAbsoluteY:
		address = this.Read16(this.PC+1) + uint16(this.Y)
		pageCrossed = pagesDiffer(address-uint16(this.Y), address)
	case modeAccumulator:
		address = 0
	case modeImmediate:
		address = this.PC + 1
	case modeImplied:
		address = 0
	case modeIndexedIndirect:
		address = this.read16bug(uint16(this.Read(this.PC+1) + this.X))
	case modeIndirect:
		address = this.read16bug(this.Read16(this.PC + 1))
	case modeIndirectIndexed:
		address = this.read16bug(uint16(this.Read(this.PC+1))) + uint16(this.Y)
		pageCrossed = pagesDiffer(address-uint16(this.Y), address)
	case modeRelative:
		offset := uint16(this.Read(this.PC + 1))
		if offset < 0x80 {
			address = this.PC + 2 + offset
		} else {
			address = this.PC + 2 + offset - 0x100
		}
	case modeZeroPage:
		address = uint16(this.Read(this.PC + 1))
	case modeZeroPageX:
		address = uint16(this.Read(this.PC+1)+this.X) & 0xff
	case modeZeroPageY:
		address = uint16(this.Read(this.PC+1)+this.Y) & 0xff
	}

	this.PC += uint16(instructionSizes[opcode])
	this.Cycles += uint64(instructionCycles[opcode])
	if pageCrossed {
		this.Cycles += uint64(instructionPageCycles[opcode])
	}
	info := &stepInfo{address, this.PC, mode}

	if Halt {
		println(info.address, " ", info.pc, " ", info.mode, " ", instructionNames[opcode])
	}

	this.table[opcode](info)

	return int(this.Cycles - cycles)
}

// NMI - Non-Maskable Interrupt
func (this *CPU) nmi() {
	this.push16(this.PC)
	this.php(nil)
	this.PC = this.Read16(0xFFFA)
	this.I = 1
	this.Cycles += 7
}

// IRQ - IRQ Interrupt
func (this *CPU) irq() {
	this.push16(this.PC)
	this.php(nil)
	this.PC = this.Read16(0xFFFE)
	this.I = 1
	this.Cycles += 7
}

// ADC - Add with Carry
func (this *CPU) adc(info *stepInfo) {
	a := this.A
	b := this.Read(info.address)
	c := this.C
	this.A = a + b + c
	this.setZN(this.A)
	if int(a)+int(b)+int(c) > 0xFF {
		this.C = 1
	} else {
		this.C = 0
	}
	if (a^b)&0x80 == 0 && (a^this.A)&0x80 != 0 {
		this.V = 1
	} else {
		this.V = 0
	}
}

// AND - Logical AND
func (this *CPU) and(info *stepInfo) {
	this.A = this.A & this.Read(info.address)
	this.setZN(this.A)
}

// ASL - Arithmetic Shift Left
func (this *CPU) asl(info *stepInfo) {
	if info.mode == modeAccumulator {
		this.C = (this.A >> 7) & 1
		this.A <<= 1
		this.setZN(this.A)
	} else {
		value := this.Read(info.address)
		this.C = (value >> 7) & 1
		value <<= 1
		this.Write(info.address, value)
		this.setZN(value)
	}
}

// BCC - Branch if Carry Clear
func (this *CPU) bcc(info *stepInfo) {
	if this.C == 0 {
		this.PC = info.address
		this.addBranchCycles(info)
	}
}

// BCS - Branch if Carry Set
func (this *CPU) bcs(info *stepInfo) {
	if this.C != 0 {
		this.PC = info.address
		this.addBranchCycles(info)
	}
}

// BEQ - Branch if Equal
func (this *CPU) beq(info *stepInfo) {
	if this.Z != 0 {
		this.PC = info.address
		this.addBranchCycles(info)
	}
}

// BIT - Bit Test
func (this *CPU) bit(info *stepInfo) {
	value := this.Read(info.address)
	this.V = (value >> 6) & 1
	this.setZ(value & this.A)
	this.setN(value)
}

// BMI - Branch if Minus
func (this *CPU) bmi(info *stepInfo) {
	if this.N != 0 {
		this.PC = info.address
		this.addBranchCycles(info)
	}
}

// BNE - Branch if Not Equal
func (this *CPU) bne(info *stepInfo) {
	if this.Z == 0 {
		this.PC = info.address
		this.addBranchCycles(info)
	}
}

// BPL - Branch if Positive
func (this *CPU) bpl(info *stepInfo) {
	if this.N == 0 {
		this.PC = info.address
		this.addBranchCycles(info)
	}
}

// BRK - Force Interrupt
func (this *CPU) brk(info *stepInfo) {
	this.push16(this.PC)
	this.php(info)
	this.sei(info)
	this.PC = this.Read16(0xFFFE)
}

// BVC - Branch if Overflow Clear
func (this *CPU) bvc(info *stepInfo) {
	if this.V == 0 {
		this.PC = info.address
		this.addBranchCycles(info)
	}
}

// BVS - Branch if Overflow Set
func (this *CPU) bvs(info *stepInfo) {
	if this.V != 0 {
		this.PC = info.address
		this.addBranchCycles(info)
	}
}

// CLC - Clear Carry Flag
func (this *CPU) clc(info *stepInfo) {
	this.C = 0
}

// CLD - Clear Decimal Mode
func (this *CPU) cld(info *stepInfo) {
	this.D = 0
}

// CLI - Clear Interrupt Disable
func (this *CPU) cli(info *stepInfo) {
	this.I = 0
}

// CLV - Clear Overflow Flag
func (this *CPU) clv(info *stepInfo) {
	this.V = 0
}

// CMP - Compare
func (this *CPU) cmp(info *stepInfo) {
	value := this.Read(info.address)
	this.compare(this.A, value)
}

// CPX - Compare X Register
func (this *CPU) cpx(info *stepInfo) {
	value := this.Read(info.address)
	this.compare(this.X, value)
}

// CPY - Compare Y Register
func (this *CPU) cpy(info *stepInfo) {
	value := this.Read(info.address)
	this.compare(this.Y, value)
}

// DEC - Decrement Memory
func (this *CPU) dec(info *stepInfo) {
	value := this.Read(info.address) - 1
	this.Write(info.address, value)
	this.setZN(value)
}

// DEX - Decrement X Register
func (this *CPU) dex(info *stepInfo) {
	this.X--
	this.setZN(this.X)
}

// DEY - Decrement Y Register
func (this *CPU) dey(info *stepInfo) {
	this.Y--
	this.setZN(this.Y)
}

// EOR - Exclusive OR
func (this *CPU) eor(info *stepInfo) {
	this.A = this.A ^ this.Read(info.address)
	this.setZN(this.A)
}

// INC - Increment Memory
func (this *CPU) inc(info *stepInfo) {
	value := this.Read(info.address) + 1
	this.Write(info.address, value)
	this.setZN(value)
}

// INX - Increment X Register
func (this *CPU) inx(info *stepInfo) {
	this.X++
	this.setZN(this.X)
}

// INY - Increment Y Register
func (this *CPU) iny(info *stepInfo) {
	this.Y++
	this.setZN(this.Y)
}

// JMP - Jump
func (this *CPU) jmp(info *stepInfo) {
	this.PC = info.address
}

// JSR - Jump to Subroutine
func (this *CPU) jsr(info *stepInfo) {
	this.push16(this.PC - 1)
	this.PC = info.address
}

// LDA - Load Accumulator
func (this *CPU) lda(info *stepInfo) {
	this.A = this.Read(info.address)
	this.setZN(this.A)
}

// LDX - Load X Register
func (this *CPU) ldx(info *stepInfo) {
	this.X = this.Read(info.address)
	this.setZN(this.X)
}

// LDY - Load Y Register
func (this *CPU) ldy(info *stepInfo) {
	this.Y = this.Read(info.address)
	this.setZN(this.Y)
}

// LSR - Logical Shift Right
func (this *CPU) lsr(info *stepInfo) {
	if info.mode == modeAccumulator {
		this.C = this.A & 1
		this.A >>= 1
		this.setZN(this.A)
	} else {
		value := this.Read(info.address)
		this.C = value & 1
		value >>= 1
		this.Write(info.address, value)
		this.setZN(value)
	}
}

// NOP - No Operation
func (this *CPU) nop(info *stepInfo) {
}

// ORA - Logical Inclusive OR
func (this *CPU) ora(info *stepInfo) {
	this.A = this.A | this.Read(info.address)
	this.setZN(this.A)
}

// PHA - Push Accumulator
func (this *CPU) pha(info *stepInfo) {
	this.push(this.A)
}

// PHP - Push Processor Status
func (this *CPU) php(info *stepInfo) {
	this.push(this.Flags() | 0x10)
}

// PLA - Pull Accumulator
func (this *CPU) pla(info *stepInfo) {
	this.A = this.pull()
	this.setZN(this.A)
}

// PLP - Pull Processor Status
func (this *CPU) plp(info *stepInfo) {
	this.SetFlags(this.pull()&0xEF | 0x20)
}

// ROL - Rotate Left
func (this *CPU) rol(info *stepInfo) {
	if info.mode == modeAccumulator {
		c := this.C
		this.C = (this.A >> 7) & 1
		this.A = (this.A << 1) | c
		this.setZN(this.A)
	} else {
		c := this.C
		value := this.Read(info.address)
		this.C = (value >> 7) & 1
		value = (value << 1) | c
		this.Write(info.address, value)
		this.setZN(value)
	}
}

// ROR - Rotate Right
func (this *CPU) ror(info *stepInfo) {
	if info.mode == modeAccumulator {
		c := this.C
		this.C = this.A & 1
		this.A = (this.A >> 1) | (c << 7)
		this.setZN(this.A)
	} else {
		c := this.C
		value := this.Read(info.address)
		this.C = value & 1
		value = (value >> 1) | (c << 7)
		this.Write(info.address, value)
		this.setZN(value)
	}
}

// RTI - Return from Interrupt
func (this *CPU) rti(info *stepInfo) {
	this.SetFlags(this.pull()&0xEF | 0x20)
	this.PC = this.pull16()
}

// RTS - Return from Subroutine
func (this *CPU) rts(info *stepInfo) {
	this.PC = this.pull16() + 1
}

// SBC - Subtract with Carry
func (this *CPU) sbc(info *stepInfo) {
	a := this.A
	b := this.Read(info.address)
	c := this.C
	this.A = a - b - (1 - c)
	this.setZN(this.A)
	if int(a)-int(b)-int(1-c) >= 0 {
		this.C = 1
	} else {
		this.C = 0
	}
	if (a^b)&0x80 != 0 && (a^this.A)&0x80 != 0 {
		this.V = 1
	} else {
		this.V = 0
	}
}

// SEC - Set Carry Flag
func (this *CPU) sec(info *stepInfo) {
	this.C = 1
}

// SED - Set Decimal Flag
func (this *CPU) sed(info *stepInfo) {
	this.D = 1
}

// SEI - Set Interrupt Disable
func (this *CPU) sei(info *stepInfo) {
	this.I = 1
}

// STA - Store Accumulator
func (this *CPU) sta(info *stepInfo) {
	if Halt {
		println("CPU.sta(), info.addr:", info.address, "this.A:", this.A)
	}
	this.Write(info.address, this.A)
}

// STX - Store X Register
func (this *CPU) stx(info *stepInfo) {
	this.Write(info.address, this.X)
}

// STY - Store Y Register
func (this *CPU) sty(info *stepInfo) {
	this.Write(info.address, this.Y)
}

// TAX - Transfer Accumulator to X
func (this *CPU) tax(info *stepInfo) {
	this.X = this.A
	this.setZN(this.X)
}

// TAY - Transfer Accumulator to Y
func (this *CPU) tay(info *stepInfo) {
	this.Y = this.A
	this.setZN(this.Y)
}

// TSX - Transfer Stack Pointer to X
func (this *CPU) tsx(info *stepInfo) {
	this.X = this.SP
	this.setZN(this.X)
}

// TXA - Transfer X to Accumulator
func (this *CPU) txa(info *stepInfo) {
	this.A = this.X
	this.setZN(this.A)
}

// TXS - Transfer X to Stack Pointer
func (this *CPU) txs(info *stepInfo) {
	this.SP = this.X
}

// TYA - Transfer Y to Accumulator
func (this *CPU) tya(info *stepInfo) {
	this.A = this.Y
	this.setZN(this.A)
}

// illegal opcodes below

func (this *CPU) ahx(info *stepInfo) {
}

func (this *CPU) alr(info *stepInfo) {
}

func (this *CPU) anc(info *stepInfo) {
}

func (this *CPU) arr(info *stepInfo) {
}

func (this *CPU) axs(info *stepInfo) {
}

func (this *CPU) dcp(info *stepInfo) {
}

func (this *CPU) isc(info *stepInfo) {
}

func (this *CPU) kil(info *stepInfo) {
}

func (this *CPU) las(info *stepInfo) {
}

func (this *CPU) lax(info *stepInfo) {
}

func (this *CPU) rla(info *stepInfo) {
}

func (this *CPU) rra(info *stepInfo) {
}

func (this *CPU) sax(info *stepInfo) {
}

func (this *CPU) shx(info *stepInfo) {
}

func (this *CPU) shy(info *stepInfo) {
}

func (this *CPU) slo(info *stepInfo) {
}

func (this *CPU) sre(info *stepInfo) {
}

func (this *CPU) tas(info *stepInfo) {
}

func (this *CPU) xaa(info *stepInfo) {
}
