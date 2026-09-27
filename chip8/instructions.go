package chip8

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
)

/* ERRORS */

var (
	ErrInvalidStackManipulation = errors.New("invalid usage of subroutine stack")
)

/* INSTRUCTIONS */

func (c *CHIP8) instr00E0() {
	clear(c.display[:][:])
}

func (c *CHIP8) instr00EE() {
	addr, err := c.stack.Pop()
	if err != nil {
		panic(ErrInvalidStackManipulation)
	}

	c.regPC = addr
}

func (c *CHIP8) instr1NNN(address uint16) {
	c.regPC = address
}

func (c *CHIP8) instr2NNN(addr uint16) {
	err := c.stack.Push(c.regPC)
	if err != nil {
		panic(ErrInvalidStackManipulation)
	}

	c.regPC = addr
}

func (c *CHIP8) instr3XNN(X uint8, hex byte) {
	valX := c.registers[X]

	if valX == hex {
		c.regPC += 2
	}
}

func (c *CHIP8) instr4XNN(X uint8, hex byte) {
	valX := c.registers[X]

	if valX != hex {
		c.regPC += 2
	}
}

func (c *CHIP8) instr5XY0(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	if valX == valY {
		c.regPC += 2
	}
}

func (c *CHIP8) instr6XNN(X uint8, val byte) {
	c.registers[X] = val
}

func (c *CHIP8) instr7XNN(X uint8, val byte) {
	c.registers[X] += val
}

func (c *CHIP8) instr8XY0(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY
}

func (c *CHIP8) instr8XY1(X uint8, Y uint8) {
	c.registers[0xF] = 0x0

	c.registers[X] |= c.registers[Y]
}

func (c *CHIP8) instr8XY2(X uint8, Y uint8) {
	c.registers[0xF] = 0x0

	c.registers[X] &= c.registers[Y]
}

func (c *CHIP8) instr8XY3(X uint8, Y uint8) {
	c.registers[0xF] = 0x0

	c.registers[X] ^= c.registers[Y]
}

func (c *CHIP8) instr8XY4(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	c.registers[X] += valY

	// Overflow detection
	if valX+valY < valX {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}
}

func (c *CHIP8) instr8XY5(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	c.registers[X] = valX - valY

	// Borrow detection
	if valX >= valY {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}
}

func (c *CHIP8) instr8XY7(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	c.registers[X] = valY - valX

	// Borrow detection
	if valY >= valX {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}
}

func (c *CHIP8) instr8XY6(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY >> 1
	c.registers[0xF] = valY &^ 0xFE
}

func (c *CHIP8) instr8XYE(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY << 1
	c.registers[0xF] = (valY &^ 0x7F) >> 7
}

func (c *CHIP8) instr9XY0(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	if valX != valY {
		c.regPC += 2
	}
}

func (c *CHIP8) instrANNN(address uint16) {
	c.regI = address
}

func (c *CHIP8) instrBNNN(addr uint16) {
	c.regPC = addr + uint16(c.registers[0x0])
}

func (c *CHIP8) instrCXNN(X uint8, hex byte) {
	num := uint8(rand.UintN(math.MaxUint8 + 1))
	c.registers[X] = num & hex
}

func (c *CHIP8) instrDXYN(X uint8, Y uint8, nibble byte) {
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

			if c.display[posX][posY] && shouldBeOn {
				c.display[posX][posY] = false
				c.registers[0xF] = 0x01
			} else if shouldBeOn {
				c.display[posX][posY] = true
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

func (c *CHIP8) instrEX9E(X uint8) {
	valX := c.registers[X]

	if slices.Contains(c.keysPressed, valX) {
		c.regPC += 2
	}
}

func (c *CHIP8) instrEXA1(X uint8) {
	valX := c.registers[X]

	if !slices.Contains(c.keysPressed, valX) {
		c.regPC += 2
	}
}

func (c *CHIP8) instrFX07(X uint8) {
	c.registers[X] = c.regDT
}

func (c *CHIP8) instrFX0A(X uint8) {
	key := c.keyFunc()
	if key == 0x0 {
		c.regPC -= 2
		return
	}

	c.registers[X] = key
}

func (c *CHIP8) instrFX15(X uint8) {
	valX := c.registers[X]
	c.regDT = valX
}

func (c *CHIP8) instrFX18(X uint8) {
	valX := c.registers[X]
	c.regST = valX
}

func (c *CHIP8) instrFX1E(X uint8) {
	valX := c.registers[X]
	c.regI += uint16(valX)
}

func (c *CHIP8) instrFX29(X uint8) {
	valX := c.registers[X]
	char := uint8(valX &^ 0xF0)

	c.regI = uint16(FONT_OFFSET) + (uint16(char) * uint16(FONT_CHAR_SIZE))
}

func (c *CHIP8) instrFX33(X uint8) {
	valX := c.registers[X]

	for i := 3 - 1; i >= 0; i-- {
		c.memory[int(c.regI)+i] = valX % 10
		valX /= 10
	}
}

func (c *CHIP8) instrFX55(X uint8) {
	for i := 0; uint(i) <= uint(X); i++ {
		val := c.registers[i]
		c.memory[int(c.regI)+i] = val
	}

	if c.quirk1 {
		c.regI += uint16(X)
	}
}

func (c *CHIP8) instrFX65(X uint8) {
	for i := 0; uint(i) <= uint(X); i++ {
		val := c.memory[int(c.regI)+i]
		c.registers[i] = val
	}

	if c.quirk1 {
		c.regI += uint16(X)
	}
}
