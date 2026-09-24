package chip8

import (
	"errors"

	"github.com/Sprinter05/chip-8/models"
)

var (
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
	if instr == 0x00E0 {
		c.instrClearScreen()
	} else if (instr &^ 0x0FFF) == 0x1000 {
		c.instrJump(address)
	} else if (instr &^ 0x0FFF) == 0x6000 {
		c.instrSetVX(varX, hexadecimal)
	} else if (instr &^ 0x0FFF) == 0x7000 {
		c.instrAddVX(varX, hexadecimal)
	} else if (instr &^ 0x0FFF) == 0xA000 {
		c.instrSetI(address)
	} else if (instr &^ 0x0FFF) == 0xD000 {
		c.instrDraw(varX, varY, nibble)
	} else {
		panic(ErrNotImplemented)
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
	// Clear registers
	c.I = 0x0
	for i := range c.Registers {
		c.Registers[i] = 0x0
	}

	// Clear Memory
	for i := range c.Memory {
		c.Memory[i] = 0x0
	}

	// Clear display
	for i := range c.Display {
		for j := range c.Display[i] {
			c.Display[i][j] = false
		}
	}

	// Clear stack
	c.Stack = models.NewStack[uint16](STACK_SIZE)

	// Clear pointers and timers
	c.PC = uint16(PC_START)
	c.SP = 0x0
	c.DT = 0x0
	c.ST = 0x0

	// Load font onto memory
	c.LoadFont()
}
