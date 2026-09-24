package chip8

import (
	"errors"
	"log"

	"github.com/Sprinter05/chip-8/models"
)

var (
	ErrRunningOnEmu   error = errors.New("cannot execute, running on an emulator")
	ErrNotImplemented error = errors.New("instruction not implemented")
	ErrTooBig         error = errors.New("program is too big")
)

/* PRIVATE */

func (c *CHIP8) fetch() uint16 {
	byte1 := c.Memory[c.PC]
	byte2 := c.Memory[c.PC+1]

	c.PC += 2

	return uint16(byte1)<<8 | uint16(byte2)
}

func (c *CHIP8) handleTimers() {
	if c.DT > 0 {
		c.DT--
	}

	if c.ST > 0 {
		c.ST--
	}
}

func (c *CHIP8) decodeAndRun(instr uint16) {
	// for NNN/NN/N instructions
	address := instr &^ 0xF000
	hexadecimal := byte(instr &^ 0xFF00)
	nibble := byte(instr &^ 0xFFF0)

	// for X and Y instructions
	varX := uint8((instr >> 8) &^ 0xFFF0)
	varY := uint8((instr >> 4) &^ 0xFFF0)

	// Instruction formatitng is weird
	// so we just do if-else
	switch instr {
	case 0x00E0:
		c.instrClearScreen()
	case 0x00EE:
		c.instrReturnSubroutine()
	default:
		switch instr &^ 0xFF0 {
		case 0x5000:
			c.instrSkipInstrEqualVXVY(varX, varY)
		case 0x8000:
			c.instrStoreVYinVX(varX, varY)
		case 0x8001:
			c.instrVXorVY(varX, varY)
		case 0x8002:
			c.instrVXandVY(varX, varY)
		case 0x8003:
			c.instrVXxorVY(varX, varY)
		case 0x8004:
			c.instrAddVYtoVX(varX, varY)
		case 0x8005:
			c.instrSubVYfromVX(varX, varY)
		case 0x8006:
			c.instrShiftVYright(varX, varY)
		case 0x8007:
			c.instrSubVXfromVYtoVX(varX, varY)
		case 0x8008:
			c.instrShiftVYleft(varX, varY)
		case 0x9000:
			c.instrSkipInstrNotEqualVXVY(varX, varY)
		default:
			switch instr &^ 0x0FFF {
			case 0x0000:
				panic(ErrRunningOnEmu)
			case 0x1000:
				c.instrJump(address)
			case 0x2000:
				c.instrCallSubroutine(address)
			case 0x3000:
				c.instrSkipInstrEqualVX(varX, hexadecimal)
			case 0x4000:
				c.instrSkipInstrNotEqualVX(varX, hexadecimal)
			case 0x6000:
				c.instrSetVX(varX, hexadecimal)
			case 0x7000:
				c.instrAddVX(varX, hexadecimal)
			case 0xA000:
				c.instrSetI(address)
			case 0xB000:
				c.instrJumpToAddrByV0(address)
			case 0xC000:
				c.instrSetVXToRandAndMask(varX, hexadecimal)
			case 0xD000:
				c.instrDraw(varX, varY, nibble)
			default:
				log.Printf("[INSTR] Not implemented: %x\n", instr)
				panic(ErrNotImplemented)
			}
		}
	}
}

/* PUBLIC */

func (c *CHIP8) LoadROM(program []byte) error {
	// Cant be bigger than memory size and base address
	if len(program) > (int(MEM_SIZE) - int(PC_START)) {
		return ErrTooBig
	}

	copy(c.Memory[PC_START:], program[:])

	return nil
}

func (c *CHIP8) Step() {
	// TODO: handle inputs
	c.handleTimers()

	instr := c.fetch()
	c.decodeAndRun(instr)
}

func (c *CHIP8) Reset() {
	// Clear everything
	c.I = 0x0
	clear(c.Registers[:])
	clear(c.Memory[:])
	clear(c.Display[:][:])

	// Clear stack
	c.Stack = models.NewStack[uint16](STACK_SIZE)

	// Clear pointers and timers
	c.PC = uint16(PC_START)
	c.SP = 0x0
	c.DT = 0x0
	c.ST = 0x0

	// Load font onto memory
	copy(c.Memory[FONT_OFFSET:], CHIP8_FONT)
}
