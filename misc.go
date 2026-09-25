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

func inputCallback() (byte, bool) {
	for code, expected := range KB_KEYS {
		if rl.IsKeyReleased(expected) {
			return byte(code), true
		}
	}

	return 0x0, false
}

func inputListCallback() ([]byte, bool) {
	list := make([]byte, 0, len(KB_KEYS))
	for code, expected := range KB_KEYS {
		if rl.IsKeyPressed(expected) {
			list = append(list, byte(code))
		}
	}

	return list, false
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
