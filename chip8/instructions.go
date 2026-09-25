package chip8

import (
	"errors"
	"math"
	"math/rand/v2"
)

var (
	ErrInvalidStackManipulation = errors.New("invalid usage of subroutine stack")
)

func (c *CHIP8) instrClearScreen() {
	clear(c.Display[:][:])
}

func (c *CHIP8) instrJump(address uint16) {
	c.regPC = address
}

func (c *CHIP8) instrSetVX(X uint8, val byte) {
	c.registers[X] = val
}

func (c *CHIP8) instrAddVX(X uint8, val byte) {
	c.registers[X] += val
}

func (c *CHIP8) instrSetI(address uint16) {
	c.regI = address
}

func (c *CHIP8) instrDraw(X uint8, Y uint8, nibble byte) {
	initX := uint8(c.registers[X] % uint8(DISPLAY_X))
	initY := uint8(c.registers[Y] % uint8(DISPLAY_Y))
	posX := initX
	posY := initY

	// Set VF flag
	c.registers[0xF] = 0x0

	for _, s := range c.memory[c.regI : c.regI+uint16(nibble)] {
		for i := 7; i >= 0; i-- { // byte size
			pixel := (s >> i) &^ 0xFE

			shouldBeOn := false
			if pixel == 0x01 {
				shouldBeOn = true
			}

			if c.Display[posX][posY] && shouldBeOn {
				c.Display[posX][posY] = false
				c.registers[0xF] = 0x01
			} else if shouldBeOn {
				c.Display[posX][posY] = true
			}

			posX++
			if posX >= uint8(DISPLAY_X) {
				break
			}
		}

		posY++
		posX = initX
		if posY >= uint8(DISPLAY_Y) {
			break
		}
	}
}

func (c *CHIP8) instrCallSubroutine(addr uint16) {
	c.stack.Push(c.regPC)
	c.regPC = addr
}

func (c *CHIP8) instrReturnSubroutine() {
	addr, err := c.stack.Pop()
	if err != nil {
		panic(ErrInvalidStackManipulation)
	}

	c.regPC = addr
}

func (c *CHIP8) instrSkipInstrEqualVX(X uint8, hex byte) {
	valX := c.registers[X]

	if valX == hex {
		c.regPC += 2
	}
}

func (c *CHIP8) instrSkipInstrNotEqualVX(X uint8, hex byte) {
	valX := c.registers[X]

	if valX != hex {
		c.regPC += 2
	}
}

func (c *CHIP8) instrSkipInstrEqualVXVY(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	if valX == valY {
		c.regPC += 2
	}
}

func (c *CHIP8) instrStoreVYinVX(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY
}

func (c *CHIP8) instrVXorVY(X uint8, Y uint8) {
	c.registers[X] |= c.registers[Y]
}

func (c *CHIP8) instrVXandVY(X uint8, Y uint8) {
	c.registers[X] &= c.registers[Y]
}

func (c *CHIP8) instrVXxorVY(X uint8, Y uint8) {
	c.registers[X] ^= c.registers[Y]
}

func (c *CHIP8) instrAddVYtoVX(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	// Overflow detection
	if valX+valY < valX {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}

	c.registers[X] += valY
}

func (c *CHIP8) instrSubVYfromVX(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	// Borrow detection
	if valX > valY {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}

	c.registers[X] = valX - valY
}

func (c *CHIP8) instrSubVXfromVYtoVX(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	// Borrow detection
	if valY > valX {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}

	c.registers[X] = valY - valX
}

func (c *CHIP8) instrShiftVYright(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY >> 1
	c.registers[0xF] = valY &^ 0xFE
}

func (c *CHIP8) instrShiftVYleft(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY << 1
	c.registers[0xF] = valY &^ 0x7F
}

func (c *CHIP8) instrSkipInstrNotEqualVXVY(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	if valX != valY {
		c.regPC += 2
	}
}

func (c *CHIP8) instrJumpToAddrByV0(addr uint16) {
	c.regPC = addr + uint16(c.registers[0x0])
}

func (c *CHIP8) instrSetVXToRandAndMask(X uint8, hex byte) {
	num := uint8(rand.UintN(math.MaxUint8 + 1))
	c.registers[X] = num & hex
}
