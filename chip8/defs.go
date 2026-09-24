package chip8

import (
	"github.com/Sprinter05/chip-8/models"
)

const DISPLAY_X uint = 64
const DISPLAY_Y uint = 64
const MEM_SIZE uint = 0x1000
const STACK_SIZE uint = 16
const PC_START uint = 0x200

type CHIP8 struct {
	Registers [16]byte
	I         uint16                     // Memory address register
	Memory    [MEM_SIZE]byte             // 4 Kilobytes
	Display   [DISPLAY_X][DISPLAY_Y]bool // Black or white pixels
	Stack     models.Stack[uint16]       // Maximum amount of subroutines
	PC        uint16
	SP        uint16
	DT        uint8 // Delay timer
	ST        uint8 // Sound timer
}
