package chip8

import (
	"errors"
	"fmt"
	"log"

	"github.com/Sprinter05/chip-8/models"
)

/* ERRORS */

var (
	ErrRunningOnEmu       error = errors.New("cannot execute, running on an emulator")
	ErrUnknownInstruction error = errors.New("instruction not known")
	ErrTooBig             error = errors.New("program is too big")
)

/* PRIVATE */

// An instruction is two bytes from memory
func (c *CHIP8) fetch() uint16 {
	byte1 := c.memory[c.regPC]
	byte2 := c.memory[c.regPC+1]

	c.regPC += 2

	return uint16(byte1)<<8 | uint16(byte2)
}

// Decode instruction and run its function
func (c *CHIP8) decodeAndRun(instr uint16) {
	// for NNN/NN/N instructions
	address := instr &^ 0xF000
	hexadecimal := byte(instr &^ 0xFF00)
	nibble := byte(instr &^ 0xFFF0)

	// for X and Y instructions
	varX := uint8((instr >> 8) &^ 0xFFF0)
	varY := uint8((instr >> 4) &^ 0xFFF0)

	// Check all instructions
	switch instr {
	case 0x00E0:
		c.instr00E0()
	case 0x00EE:
		c.instr00EE()
	default:
		switch instr &^ 0x0F00 {
		case 0xE09E:
			c.instrEX9E(varX)
		case 0xE0A1:
			c.instrEXA1(varX)
		case 0xF007:
			c.instrFX07(varX)
		case 0xF00A:
			c.instrFX0A(varX)
		case 0xF015:
			c.instrFX15(varX)
		case 0xF018:
			c.instrFX18(varX)
		case 0xF01E:
			c.instrFX1E(varX)
		case 0xF029:
			c.instrFX29(varX)
		case 0xF033:
			c.instrFX33(varX)
		case 0xF055:
			c.instrFX55(varX)
		case 0xF065:
			c.instrFX65(varX)
		default:
			switch instr &^ 0x0FF0 {
			case 0x5000:
				c.instr5XY0(varX, varY)
			case 0x8000:
				c.instr8XY0(varX, varY)
			case 0x8001:
				c.instr8XY1(varX, varY)
			case 0x8002:
				c.instr8XY2(varX, varY)
			case 0x8003:
				c.instr8XY3(varX, varY)
			case 0x8004:
				c.instr8XY4(varX, varY)
			case 0x8005:
				c.instr8XY5(varX, varY)
			case 0x8006:
				c.instr8XY6(varX, varY)
			case 0x8007:
				c.instr8XY7(varX, varY)
			case 0x800E:
				c.instr8XYE(varX, varY)
			case 0x9000:
				c.instr9XY0(varX, varY)
			default:
				switch instr &^ 0x0FFF {
				case 0x0000:
					// Execute machine language subroutine at address NNN
					// Not implemented on emulator
					log.Printf("%x erroneus instruction!\n", instr)
					fmt.Print(ErrRunningOnEmu)
				case 0x1000:
					c.instr1NNN(address)
				case 0x2000:
					c.instr2NNN(address)
				case 0x3000:
					c.instr3XNN(varX, hexadecimal)
				case 0x4000:
					c.instr4XNN(varX, hexadecimal)
				case 0x6000:
					c.instr6XNN(varX, hexadecimal)
				case 0x7000:
					c.instr7XNN(varX, hexadecimal)
				case 0xA000:
					c.instrANNN(address)
				case 0xB000:
					c.instrBNNN(address)
				case 0xC000:
					c.instrCXNN(varX, hexadecimal)
				case 0xD000:
					c.instrDXYN(varX, varY, nibble)
				default:
					// Every instruction is implemented so this should not run
					log.Printf("%x unknown instruction!\n", instr)
					panic(ErrUnknownInstruction)
				}
			}
		}
	}
}

/* PUBLIC */

// Decrement both the delay and sound timer
// Should only run once per frame
func (c *CHIP8) DecrementTimers() {
	if c.regDT > 0 {
		c.regDT--
	}

	if c.regST > 0 {
		c.regST--
	}
}

// Toggle emulator pause state
func (c *CHIP8) TogglePause() {
	c.paused = !c.paused
}

// Load ROM into memory
func (c *CHIP8) LoadROM(program []byte) error {
	// Cant be bigger than memory size and base address
	if len(program) > (int(MEM_SIZE) - int(PC_START)) {
		return ErrTooBig
	}

	copy(c.memory[PC_START:], program[:])
	c.loaded = true

	return nil
}

// Run a single instruction
// Returns true if the last instruction ran was drawing
func (c *CHIP8) Step() bool {
	if c.paused || !c.loaded {
		return false
	}

	// Input handling
	keys := c.inputFunc()
	c.keysPressed = keys

	// Fetch, decode and execute
	instr := c.fetch()
	c.decodeAndRun(instr)

	// Return true if we have drawn onto the screen
	// This allows the DISPLAY_WAIT quirk to work
	if instr&^0x0FFF == 0xD000 && !c.quirk3 {
		return true
	}

	return false
}

// Reset emulator state
func (c *CHIP8) Reset() {
	// Clear data
	c.regI = 0x0
	clear(c.registers[:])
	clear(c.memory[:])
	clear(c.display[:][:])

	// Clear stack
	c.stack = models.NewStack[uint16](STACK_SIZE)

	// Clear pointers and timers
	c.regPC = uint16(PC_START)
	c.regSP = 0x0
	c.regDT = 0x0
	c.regST = 0x0

	// Clear state
	c.paused = false
	c.loaded = false
	clear(c.keysPressed)

	// Load font onto memory
	copy(c.memory[FONT_OFFSET:], CHIP8_FONT)
}
