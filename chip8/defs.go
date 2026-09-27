package chip8

import (
	"github.com/Sprinter05/chip-8/models"
)

/* DEFINITIONS */

const DISPLAY_X uint = 64
const DISPLAY_Y uint = 32
const MEM_SIZE uint = 0x1000
const STACK_SIZE uint = 16
const PC_START uint = 0x200
const FONT_OFFSET uint = 0x050
const FONT_CHAR_SIZE uint = 5

var CHIP8_FONT = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x60, 0x20, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0xA0, 0xA0, 0xF0, 0x20, 0x20, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x10, 0x10, 0x10, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xF0, 0x50, 0x70, 0x50, 0xF0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xF0, 0x50, 0x50, 0x50, 0xF0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type CHIP8 struct {
	// Machine values
	registers [16]byte
	regI      uint16                     // Memory address register
	memory    [MEM_SIZE]byte             // 4 Kilobytes
	stack     models.Stack[uint16]       // Maximum amount of subroutines
	regPC     uint16                     // Program Counter
	regSP     uint16                     // Stack Pointer
	regDT     uint8                      // Delay Timer
	regST     uint8                      // Sound Timer
	display   [DISPLAY_X][DISPLAY_Y]bool // Black or white pixels

	// State values
	keysPressed []byte        // Keys held down that frame
	inputFunc   func() []byte // Function that returns keys pressed that frame
}

/* SETUP FUNCTIONS */

func (c *CHIP8) GetSoundTimer() uint {
	return uint(c.regST)
}

func (c *CHIP8) GetDisplay() [DISPLAY_X][DISPLAY_Y]bool {
	return c.display
}

func (c *CHIP8) SetInputFunction(fun func() []byte) {
	c.inputFunc = fun
}
