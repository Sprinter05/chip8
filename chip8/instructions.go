package chip8

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
