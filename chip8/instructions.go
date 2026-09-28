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

// Clear the screen
func (c *CHIP8) instr00E0() {
	clear(c.display[:][:])
}

// Return from a subroutine
func (c *CHIP8) instr00EE() {
	addr, err := c.stack.Pop()
	if err != nil {
		panic(ErrInvalidStackManipulation)
	}

	c.regPC = addr
}

// Jump to address NNN
func (c *CHIP8) instr1NNN(address uint16) {
	c.regPC = address
}

// Execute subroutine starting at address NNN
func (c *CHIP8) instr2NNN(addr uint16) {
	err := c.stack.Push(c.regPC)
	if err != nil {
		panic(ErrInvalidStackManipulation)
	}

	c.regPC = addr
}

// Skip the following instruction if the value of
// register VX equals NN
func (c *CHIP8) instr3XNN(X uint8, hex byte) {
	valX := c.registers[X]

	if valX == hex {
		c.regPC += 2
	}
}

// Skip the following instruction if the value of
// register VX is not equal to NN
func (c *CHIP8) instr4XNN(X uint8, hex byte) {
	valX := c.registers[X]

	if valX != hex {
		c.regPC += 2
	}
}

// Skip the following instruction if the value of
// register VX is equal to the value of register VY
func (c *CHIP8) instr5XY0(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	if valX == valY {
		c.regPC += 2
	}
}

// Store number NN in register VX
func (c *CHIP8) instr6XNN(X uint8, val byte) {
	c.registers[X] = val
}

// Add the value NN to register VX
func (c *CHIP8) instr7XNN(X uint8, val byte) {
	c.registers[X] += val
}

// Store the value of register VY in register VX
func (c *CHIP8) instr8XY0(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY
}

// Set VX to VX OR VY
func (c *CHIP8) instr8XY1(X uint8, Y uint8) {
	c.registers[X] |= c.registers[Y]

	// Clear VF if quirk enabled
	if !c.quirk2 {
		c.registers[0xF] = 0x0
	}
}

// Set VX to VX AND VY
func (c *CHIP8) instr8XY2(X uint8, Y uint8) {
	c.registers[X] &= c.registers[Y]

	// Clear VF if quirk enabled
	if !c.quirk2 {
		c.registers[0xF] = 0x0
	}
}

// Set VX to VX XOR VY
func (c *CHIP8) instr8XY3(X uint8, Y uint8) {
	c.registers[X] ^= c.registers[Y]

	// Clear VF if quirk enabled
	if !c.quirk2 {
		c.registers[0xF] = 0x0
	}
}

// Add the value of register VY to register VX
// Set VF to 01 if a carry occurs
// Set VF to 00 if a carry does not occur
func (c *CHIP8) instr8XY4(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	c.registers[X] += valY

	// Overflow detection
	// VF is set at the end
	if valX+valY < valX {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}
}

// Subtract the value of register VY from register VX
// Set VF to 00 if a borrow occurs
// Set VF to 01 if a borrow does not occur
func (c *CHIP8) instr8XY5(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	c.registers[X] = valX - valY

	// Borrow detection
	// VF is set at the end
	if valX >= valY {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}
}

// Store the value of register VY shifted right one bit in register VX
// Set register VF to the least significant bit prior to the shift
// VY is unchanged
func (c *CHIP8) instr8XY6(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY >> 1

	// VF is set at the end
	c.registers[0xF] = valY &^ 0xFE
}

// Set register VX to the value of VY minus VX
// Set VF to 00 if a borrow occurs
// Set VF to 01 if a borrow does not occur
func (c *CHIP8) instr8XY7(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	c.registers[X] = valY - valX

	// Borrow detection
	// VF is set at the end
	if valY >= valX {
		c.registers[0xF] = 0x1
	} else {
		c.registers[0xF] = 0x0
	}
}

// Store the value of register VY shifted left one bit in register VX
// Set register VF to the most significant bit prior to the shift
// VY is unchanged
func (c *CHIP8) instr8XYE(X uint8, Y uint8) {
	valY := c.registers[Y]
	c.registers[X] = valY << 1

	// VF is set at the end
	c.registers[0xF] = (valY &^ 0x7F) >> 7
}

// Skip the following instruction if the value of
// register VX is not equal to the value of register VY
func (c *CHIP8) instr9XY0(X uint8, Y uint8) {
	valX := c.registers[X]
	valY := c.registers[Y]

	if valX != valY {
		c.regPC += 2
	}
}

// Store memory address NNN in register I
func (c *CHIP8) instrANNN(address uint16) {
	c.regI = address
}

// Jump to address NNN + V0
func (c *CHIP8) instrBNNN(addr uint16) {
	c.regPC = addr + uint16(c.registers[0x0])
}

// Set VX to a random number with a mask of NN
func (c *CHIP8) instrCXNN(X uint8, hex byte) {
	num := uint8(rand.UintN(math.MaxUint8 + 1))
	c.registers[X] = num & hex
}

// Draw a sprite at position VX, VY with N bytes of sprite data
// starting at the address stored in I
// Set VF to 01 if any set pixels are changed to unset, and 00 otherwise
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

			// Clipping (do not wrap)
			posX++
			if posX >= uint8(DISPLAY_X) {
				break
			}
		}

		// Clipping (do not wrap)
		posY++
		posX = initX
		if posY >= uint8(DISPLAY_Y) {
			break
		}
	}
}

// Skip the following instruction if the key corresponding
// to the hex value currently stored in register VX is pressed
func (c *CHIP8) instrEX9E(X uint8) {
	valX := c.registers[X]

	if slices.Contains(c.keysPressed, valX) {
		c.regPC += 2
	}
}

// Skip the following instruction if the key corresponding
// to the hex value currently stored in register VX is not pressed
func (c *CHIP8) instrEXA1(X uint8) {
	valX := c.registers[X]

	if !slices.Contains(c.keysPressed, valX) {
		c.regPC += 2
	}
}

// Store the current value of the delay timer in register VX
func (c *CHIP8) instrFX07(X uint8) {
	c.registers[X] = c.regDT
}

// Wait for a keypress and store the result in register VX
func (c *CHIP8) instrFX0A(X uint8) {
	key := c.keyFunc()
	if key == 0x0 {
		c.regPC -= 2
		return
	}

	c.registers[X] = key
}

// Set the delay timer to the value of register VX
func (c *CHIP8) instrFX15(X uint8) {
	valX := c.registers[X]
	c.regDT = valX
}

// Set the sound timer to the value of register VX
func (c *CHIP8) instrFX18(X uint8) {
	valX := c.registers[X]
	c.regST = valX
}

// Add the value stored in register VX to register I
func (c *CHIP8) instrFX1E(X uint8) {
	valX := c.registers[X]
	c.regI += uint16(valX)
}

// Set I to the memory address of the sprite data corresponding
// to the hexadecimal digit stored in register VX
func (c *CHIP8) instrFX29(X uint8) {
	valX := c.registers[X]
	char := uint8(valX &^ 0xF0)

	c.regI = uint16(FONT_OFFSET) + (uint16(char) * uint16(FONT_CHAR_SIZE))
}

// Store the binary-coded decimal equivalent of the value
// stored in register VX at addresses I, I + 1, and I + 2
func (c *CHIP8) instrFX33(X uint8) {
	valX := c.registers[X]

	for i := 3 - 1; i >= 0; i-- {
		c.memory[int(c.regI)+i] = valX % 10
		valX /= 10
	}
}

// Store the values of registers V0 to VX
// inclusive in memory starting at address I
// I is set to I + X + 1 after operation
func (c *CHIP8) instrFX55(X uint8) {
	for i := 0; uint(i) <= uint(X); i++ {
		val := c.registers[i]
		c.memory[int(c.regI)+i] = val
	}

	// If the quirk is enabled, increment I
	if c.quirk1 {
		c.regI += uint16(X) + 1
	}
}

// Fill registers V0 to VX inclusive with the
// values stored in memory starting at address I
// I is set to I + X + 1 after operation
func (c *CHIP8) instrFX65(X uint8) {
	for i := 0; uint(i) <= uint(X); i++ {
		val := c.memory[int(c.regI)+i]
		c.registers[i] = val
	}

	// If the quirk is enabled, increment I
	if c.quirk1 {
		c.regI += uint16(X) + 1
	}
}
