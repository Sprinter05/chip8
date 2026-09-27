package main

import (
	"fmt"
	"os"

	"github.com/Sprinter05/chip-8/chip8"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	zenity "github.com/ncruces/zenity"
)

/* DEFINITIONS */

// VALUES
const MIN_FPS = 30
const MAX_FPS = 960
const DEFAULT_FPS = 60
const MIN_VOLUME = 0
const MAX_VOLUME = 100
const DEFAULT_VOLUME = 30
const MIN_BUZZER_FREQ = 400
const MAX_BUZZER_FREQ = 1500
const DEFAULT_BUZZER_FREQ = 440

// GUI
const GUI_HEIGHT = 100
const GUI_SLIDER_PADDING_X = 90
const GUI_BUTTON_PADDING_X = 90
const GUI_SLIDER_PADDING_Y = 10
const GUI_BUTTON_PADDING_Y = 20
const GUI_SLIDER_SEPARATION = 30
const GUI_BUTTON_SEPARATION = 70
const GUI_SLIDER_RESET_OFFSET = 200
const GUI_SLIDER_SIZE = 20
const GUI_SLIDER_LENGTH = 120
const GUI_BUTTON_SIZE = 60
const GUI_RESET_ICON_SIZE = 1
const GUI_ICON_SIZE = 3

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

func loadProgram(emu *chip8.CHIP8) error {
	f, err := os.ReadFile(fileROM)
	if err != nil {
		if os.IsNotExist(err) {
			emu.TogglePause()
			return nil
		}

		return err
	}

	if err := emu.LoadROM(f); err != nil {
		return err
	}

	return nil
}

func keyHandler() func() byte {
	return func() byte {
		code := rl.GetKeyPressed()
		hex, exists := KB_KEYS_MAP[code]
		if code == 0x0 || !exists {
			return 0x0
		}

		return hex
	}
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

func checkInterfaceToggle(show bool) bool {
	if rl.IsKeyPressed(rl.KeySpace) {
		if show {
			rl.SetWindowSize(rl.GetScreenWidth(), rl.GetScreenHeight()-GUI_HEIGHT)
		}

		if !show {
			rl.SetWindowSize(rl.GetScreenWidth(), rl.GetScreenHeight()+GUI_HEIGHT)
		}

		return !show
	}

	return show
}

func renderInterface(emu *chip8.CHIP8, fps *float32, volume *float32, buzzer *float32) {
	width := float32(rl.GetScreenWidth())
	height := float32(rl.GetScreenHeight() - GUI_HEIGHT)

	/* SLIDERS */

	// FPS Control
	gui.SetIconScale(GUI_RESET_ICON_SIZE)
	gui.Slider(
		rl.NewRectangle(
			GUI_SLIDER_PADDING_X, height+GUI_SLIDER_PADDING_Y,
			GUI_SLIDER_LENGTH, GUI_SLIDER_SIZE,
		), "FPS", fmt.Sprintf("%2.2f", *fps),
		fps, MIN_FPS, MAX_FPS,
	)
	butR1 := gui.Button(
		rl.NewRectangle(
			GUI_SLIDER_PADDING_X+GUI_SLIDER_RESET_OFFSET, height+GUI_SLIDER_PADDING_Y,
			GUI_SLIDER_SIZE, GUI_SLIDER_SIZE,
		), "#74#",
	)
	if butR1 {
		*fps = DEFAULT_FPS
	}

	// Volume Control
	gui.SetIconScale(GUI_RESET_ICON_SIZE)
	gui.Slider(
		rl.NewRectangle(
			GUI_SLIDER_PADDING_X, height+GUI_SLIDER_PADDING_Y+GUI_SLIDER_SEPARATION,
			GUI_SLIDER_LENGTH, GUI_SLIDER_SIZE,
		), "Volume", fmt.Sprintf("%2.0f %%", *volume),
		volume, MIN_VOLUME, MAX_VOLUME,
	)
	butR2 := gui.Button(
		rl.NewRectangle(
			GUI_SLIDER_PADDING_X+GUI_SLIDER_RESET_OFFSET, height+GUI_SLIDER_PADDING_Y+GUI_SLIDER_SEPARATION,
			GUI_SLIDER_SIZE, GUI_SLIDER_SIZE,
		), "#74#",
	)
	if butR2 {
		*volume = DEFAULT_VOLUME
	}

	// Buzzer Control
	gui.SetIconScale(GUI_RESET_ICON_SIZE)
	gui.Slider(
		rl.NewRectangle(
			GUI_SLIDER_PADDING_X, height+GUI_SLIDER_PADDING_Y+GUI_SLIDER_SEPARATION*2,
			GUI_SLIDER_LENGTH, GUI_SLIDER_SIZE,
		), "Buzzer", fmt.Sprintf("%2.0f Hz", *buzzer),
		buzzer, MIN_BUZZER_FREQ, MAX_BUZZER_FREQ,
	)
	butR3 := gui.Button(
		rl.NewRectangle(
			GUI_SLIDER_PADDING_X+GUI_SLIDER_RESET_OFFSET, height+GUI_SLIDER_PADDING_Y+GUI_SLIDER_SEPARATION*2,
			GUI_SLIDER_SIZE, GUI_SLIDER_SIZE,
		), "#74#",
	)
	if butR3 {
		*buzzer = DEFAULT_BUZZER_FREQ
	}

	/* BUTTONS */

	// Pause
	gui.SetIconScale(GUI_ICON_SIZE)
	pauseText := "#132#"
	if emu.IsPaused() {
		pauseText = "#131#"
	}
	if gui.Button(
		rl.NewRectangle(
			width-GUI_BUTTON_PADDING_X-GUI_BUTTON_SEPARATION*3, height+GUI_BUTTON_PADDING_Y,
			GUI_BUTTON_SIZE, GUI_BUTTON_SIZE,
		), pauseText,
	) {
		emu.TogglePause()
	}

	// Reset
	gui.SetIconScale(GUI_ICON_SIZE)
	if gui.Button(
		rl.NewRectangle(
			width-GUI_BUTTON_PADDING_X-GUI_BUTTON_SEPARATION*2, height+GUI_BUTTON_PADDING_Y,
			GUI_BUTTON_SIZE, GUI_BUTTON_SIZE,
		), "#76#",
	) {
		emu.Reset()
		if err := loadProgram(emu); err != nil {
			panic(err)
		}
	}

	// Load
	gui.SetIconScale(GUI_ICON_SIZE)
	if gui.Button(
		rl.NewRectangle(
			width-GUI_BUTTON_PADDING_X-GUI_BUTTON_SEPARATION, height+GUI_BUTTON_PADDING_Y,
			GUI_BUTTON_SIZE, GUI_BUTTON_SIZE,
		), "#5#",
	) {
		file, err := zenity.SelectFile()
		if err != nil {
			if err == zenity.ErrCanceled {
				return
			}

			panic(err)
		}

		fileROM = file
		emu.Reset()
		if err := loadProgram(emu); err != nil {
			panic(err)
		}
	}

	// Reset window
	gui.SetIconScale(GUI_ICON_SIZE)
	if gui.Button(
		rl.NewRectangle(
			width-GUI_BUTTON_PADDING_X, height+GUI_BUTTON_PADDING_Y,
			GUI_BUTTON_SIZE, GUI_BUTTON_SIZE,
		), "#105#",
	) {
		rl.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT+GUI_HEIGHT)
	}

}
