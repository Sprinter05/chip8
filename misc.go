package main

import (
	"fmt"
	"os"

	"github.com/Sprinter05/chip-8/chip8"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

/* DEFINITIONS */

const MIN_FPS = 30
const MAX_FPS = 960
const DEFAULT_FPS = 60
const MIN_BUZZER_FREQ = 400
const MAX_BUZZER_FREQ = 1500
const DEFAULT_BUZZER_FREQ = 440

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

/* MISCELLANEOUS FUNCTIONS */

func defaultState() (*chip8.CHIP8, error) {
	emu := new(chip8.CHIP8)
	emu.Reset() // Resets all values
	emu.SetInputFunction(inputHandler())

	f, err := os.ReadFile(fileROM)
	if err != nil {
		return nil, err
	}

	if err := emu.LoadROM(f); err != nil {
		return nil, err
	}

	return emu, nil
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
	display := emu.GetDisplay()
	for x := range display {
		for y, v := range display[x] {
			if v {
				rl.DrawPixel(int32(x), int32(y), rl.White)
			}
		}
	}
	rl.EndTextureMode()
}

func audioCallback(c *chip8.CHIP8, freq *float32) rl.AudioCallback {
	index := 0
	return func(data []float32, frames int) {
		if c.GetSoundTimer() <= 0 {
			return
		}

		wavelength := SAMPLE_RATE / int(*freq)

		for i := range frames {
			if index < wavelength/2 {
				data[i] = 1.0
			} else {
				data[i] = -1.0
			}

			index = (index + 1) % wavelength
		}

	}
}

func renderInterface(emu *chip8.CHIP8, fps *float32, buzzer *float32) {
	// FPS Control
	gui.Slider(
		rl.NewRectangle(GUI_PADDING_X, WINDOW_HEIGHT+GUI_PADDING_Y, 120, 30),
		"FPS", fmt.Sprintf("%2.2f", *fps),
		fps, MIN_FPS, MAX_FPS,
	)
	butR1 := gui.Button(
		rl.NewRectangle(GUI_PADDING_X+165, WINDOW_HEIGHT+GUI_PADDING_Y, 30, 30),
		"#77#",
	)
	if butR1 {
		*fps = DEFAULT_FPS
	}

	// Buzzer Control
	gui.Slider(
		rl.NewRectangle(GUI_PADDING_X, WINDOW_HEIGHT+GUI_PADDING_Y+40, 120, 30),
		"Buzzer", fmt.Sprintf("%2.0f Hz", *buzzer),
		buzzer, MIN_BUZZER_FREQ, MAX_BUZZER_FREQ,
	)
	butR2 := gui.Button(
		rl.NewRectangle(GUI_PADDING_X+165, WINDOW_HEIGHT+GUI_PADDING_Y+40, 30, 30),
		"#77#",
	)
	if butR2 {
		*buzzer = DEFAULT_BUZZER_FREQ
	}

	// Pause
	pauseText := "#132#"
	if emu.IsPaused() {
		pauseText = "#131#"
	}
	if gui.Button(
		rl.NewRectangle(GUI_PADDING_X+255, WINDOW_HEIGHT+GUI_PADDING_Y, 70, 70),
		pauseText,
	) {
		emu.TogglePause()
	}

}
