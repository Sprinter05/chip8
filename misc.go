package main

import (
	"github.com/Sprinter05/chip-8/chip8"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var KB_KEYS_MAP = map[int32]byte{
	rl.KeyOne:   0x1,
	rl.KeyTwo:   0x2,
	rl.KeyThree: 0x3,
	rl.KeyFour:  0xC,
	rl.KeyQ:     0x4,
	rl.KeyW:     0x5,
	rl.KeyE:     0x6,
	rl.KeyR:     0xD,
	rl.KeyA:     0x7,
	rl.KeyS:     0x8,
	rl.KeyD:     0x9,
	rl.KeyF:     0xE,
	rl.KeyZ:     0xA,
	rl.KeyX:     0x0,
	rl.KeyC:     0xB,
	rl.KeyV:     0xF,
}

func inputHandler() func() []byte {
	return func() []byte {
		keys := make([]byte, 0, len(KB_KEYS_MAP))
		for expected, code := range KB_KEYS_MAP {
			if rl.IsKeyDown(expected) {
				keys = append(keys, code)
			}
		}

		return keys
	}
}

func drawOnTexture(canvas rl.RenderTexture2D, emu *chip8.CHIP8) {
	rl.BeginTextureMode(canvas)
	rl.ClearBackground(rl.Black)
	for x := range emu.Display {
		for y, v := range emu.Display[x] {
			if v {
				rl.DrawPixel(int32(x), int32(y), rl.White)
			}
		}
	}
	rl.EndTextureMode()
}
