package main

import (
	"flag"

	"github.com/Sprinter05/chip-8/chip8"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

/* CONSTANTS */

const FPS = 60
const FRAMETIME_US float64 = 1000000.0 / FPS // microseconds
const WINDOW_WIDTH = 640
const WINDOW_HEIGHT = 320
const AUDIO_BUFFER_SIZE = 4096
const VOLUME_DIVISION_FACTOR = 500
const SAMPLE_RATE = 44100

/* FLAGS */

var (
	fileROM string
	quirk1  bool
	quirk2  bool
	quirk3  bool
)

func init() {
	flag.StringVar(&fileROM, "rom", "rom.ch8", "ROM file to open")
	flag.BoolVar(&quirk1, "quirk1", false, "Both FX55 and FX65 increment the I register")
	flag.BoolVar(&quirk2, "quirk2", false, "Clear VF on AND, OR and XOR instructions")
	flag.BoolVar(&quirk3, "quirk3", false, "Do not wait for the display before drawing")
	flag.Parse()
}

/* MAIN */

func main() {
	// Create and setup the emulator
	emu := new(chip8.CHIP8)
	emu.SetInputFunction(inputHandler())
	emu.SetKeyFunction(keyHandler())
	emu.SetQuirks(quirk1, quirk2, quirk3)
	emu.Reset() // Resets all values
	if err := loadProgram(emu); err != nil {
		panic(err)
	}

	// Controllable values
	ipf := float32(DEFAULT_IPF)
	vol := float32(DEFAULT_VOLUME)
	freq := float32(DEFAULT_BUZZER_FREQ)
	showGUI := true

	// Initialise raylib
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowResizable)
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT+GUI_HEIGHT, "CHIP8 Interpreter")
	rl.SetTargetFPS(FPS)
	gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, 20)
	defer rl.CloseWindow()

	// Set audio stream
	rl.InitAudioDevice()
	rl.SetAudioStreamBufferSizeDefault(AUDIO_BUFFER_SIZE)
	stream := rl.LoadAudioStream(SAMPLE_RATE, 32, 1)
	rl.PlayAudioStream(stream)
	rl.SetAudioStreamCallback(stream, audioCallback(emu, &freq))
	rl.SetAudioStreamVolume(stream, vol/VOLUME_DIVISION_FACTOR)
	defer rl.UnloadAudioStream(stream)
	defer rl.CloseAudioDevice()

	// Create texture for emulator display
	canvas := rl.LoadRenderTexture(int32(chip8.DISPLAY_X), int32(chip8.DISPLAY_Y))
	defer rl.UnloadRenderTexture(canvas)

	// Main loop
	for !rl.WindowShouldClose() {
		// STATE
		oldVol := vol

		// CPU
		for range int(ipf) {
			needDraw := emu.Step()
			if needDraw {
				break
			}
		}

		// TEXTURE
		drawOnTexture(canvas, emu)

		// GUI
		showGUI = checkInterfaceToggle(showGUI)
		if showGUI {
			renderInterface(emu, &ipf, &vol, &freq)
		}

		// STATE
		if vol != oldVol {
			rl.SetAudioStreamVolume(stream, vol/VOLUME_DIVISION_FACTOR)
		}

		// DRAW
		rl.BeginDrawing()
		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

		// RESIZE
		src := rl.NewRectangle(0, 0, float32(chip8.DISPLAY_X), -float32(chip8.DISPLAY_Y))
		guiHeight := GUI_HEIGHT
		if !showGUI {
			guiHeight = 0
		}
		dst := rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()-guiHeight))
		rl.DrawTexturePro(canvas.Texture, src, dst, rl.NewVector2(0, 0), 0, rl.White)

		rl.EndDrawing()
	}
}
