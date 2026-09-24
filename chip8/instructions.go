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
	c.PC = address
}

func (c *CHIP8) instrSetVX(X uint8, val byte) {
	c.Registers[X] = val
}

func (c *CHIP8) instrAddVX(X uint8, val byte) {
	c.Registers[X] += val
}

func (c *CHIP8) instrSetI(address uint16) {
	c.I = address
}

func (c *CHIP8) instrDraw(X uint8, Y uint8, nibble byte) {
	initX := uint8(c.Registers[X] % uint8(DISPLAY_X))
	initY := uint8(c.Registers[Y] % uint8(DISPLAY_Y))
	posX := initX
	posY := initY

	// Set VF flag
	c.Registers[0xF] = 0x0

	for _, s := range c.Memory[c.I : c.I+uint16(nibble)] {
		for i := 7; i >= 0; i-- { // byte size
			pixel := (s >> i) &^ 0xFE

			shouldBeOn := false
			if pixel == 0x01 {
				shouldBeOn = true
			}

			if c.Display[posX][posY] && shouldBeOn {
				c.Display[posX][posY] = false
				c.Registers[0xF] = 0x01
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
	c.Stack.Push(c.PC)
	c.PC = addr
}

func (c *CHIP8) instrReturnSubroutine() {
	addr, err := c.Stack.Pop()
	if err != nil {
		panic(ErrInvalidStackManipulation)
	}

	c.PC = addr
}

func (c *CHIP8) instrSkipInstrEqualVX(X uint8, hex byte) {
	valX := c.Registers[X]

	if valX == hex {
		c.PC += 2
	}
}

func (c *CHIP8) instrSkipInstrNotEqualVX(X uint8, hex byte) {
	valX := c.Registers[X]

	if valX != hex {
		c.PC += 2
	}
}

func (c *CHIP8) instrSkipInstrEqualVXVY(X uint8, Y uint8) {
	valX := c.Registers[X]
	valY := c.Registers[Y]

	if valX == valY {
		c.PC += 2
	}
}

func (c *CHIP8) instrStoreVYinVX(X uint8, Y uint8) {
	valY := c.Registers[Y]
	c.Registers[X] = valY
}

func (c *CHIP8) instrVXorVY(X uint8, Y uint8) {
	c.Registers[X] |= c.Registers[Y]
}

func (c *CHIP8) instrVXandVY(X uint8, Y uint8) {
	c.Registers[X] &= c.Registers[Y]
}

func (c *CHIP8) instrVXxorVY(X uint8, Y uint8) {
	c.Registers[X] ^= c.Registers[Y]
}

func (c *CHIP8) instrAddVYtoVX(X uint8, Y uint8) {
	valX := c.Registers[X]
	valY := c.Registers[Y]

	// Overflow detection
	if valX+valY < valX {
		c.Registers[0xF] = 0x1
	} else {
		c.Registers[0xF] = 0x0
	}

	c.Registers[X] += valY
}

func (c *CHIP8) instrSubVYfromVX(X uint8, Y uint8) {
	valX := c.Registers[X]
	valY := c.Registers[Y]

	// Borrow detection
	if valX > valY {
		c.Registers[0xF] = 0x1
	} else {
		c.Registers[0xF] = 0x0
	}

	c.Registers[X] = valX - valY
}

func (c *CHIP8) instrSubVXfromVYtoVX(X uint8, Y uint8) {
	valX := c.Registers[X]
	valY := c.Registers[Y]

	// Borrow detection
	if valY > valX {
		c.Registers[0xF] = 0x1
	} else {
		c.Registers[0xF] = 0x0
	}

	c.Registers[X] = valY - valX
}

func (c *CHIP8) instrShiftVYright(X uint8, Y uint8) {
	valY := c.Registers[Y]
	c.Registers[X] = valY >> 1
	c.Registers[0xF] = valY &^ 0xFE
}

func (c *CHIP8) instrShiftVYleft(X uint8, Y uint8) {
	valY := c.Registers[Y]
	c.Registers[X] = valY << 1
	c.Registers[0xF] = valY &^ 0x7F
}

func (c *CHIP8) instrSkipInstrNotEqualVXVY(X uint8, Y uint8) {
	valX := c.Registers[X]
	valY := c.Registers[Y]

	if valX != valY {
		c.PC += 2
	}
}

func (c *CHIP8) instrJumpToAddrByV0(addr uint16) {
	c.PC = addr + uint16(c.Registers[0x0])
}

func (c *CHIP8) instrSetVXToRandAndMask(X uint8, hex byte) {
	num := uint8(rand.UintN(math.MaxUint8 + 1))
	c.Registers[X] = num & hex
}
