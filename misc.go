package main

import (
	"github.com/Sprinter05/chip-8/chip8"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var KB_KEYS = []int32{
	rl.KeyOne, rl.KeyTwo, rl.KeyThree, rl.KeyFour,
	rl.KeyQ, rl.KeyW, rl.KeyE, rl.KeyR,
	rl.KeyA, rl.KeyS, rl.KeyD, rl.KeyF,
	rl.KeyZ, rl.KeyX, rl.KeyC, rl.KeyV,
}

func handleInput() func() []byte {
	return func() []byte {
		keys := make([]byte, 0, len(KB_KEYS))
		for code, expected := range KB_KEYS {
			if rl.IsKeyPressed(expected) {
				keys = append(keys, byte(code))
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
