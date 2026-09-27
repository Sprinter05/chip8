package main

import (
	"flag"
	"os"

	"github.com/Sprinter05/chip-8/chip8"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const FPS float64 = 60.0
const FRAMETIME_US float64 = 1000000.0 / FPS // microseconds
const WINDOW_WIDTH = 640
const WINDOW_HEIGHT = 320
const SAMPLE_RATE = 48000
const AUDIO_BUFFER_SIZE = 4096
const BUZZER_FREQ = 440

var fileROM string

func init() {
	flag.StringVar(&fileROM, "rom", "rom.ch8", "ROM file to open")
	flag.Parse()
}

func load(c *chip8.CHIP8) error {
	f, err := os.ReadFile(fileROM)
	if err != nil {
		return err
	}

	if err := c.LoadROM(f); err != nil {
		return err
	}

	return nil
}

func main() {
	// Create and setup the emulator
	emu := new(chip8.CHIP8)
	emu.Reset() // Resets all values
	emu.InputFunc = inputHandler()
	if err := load(emu); err != nil {
		panic(err)
	}

	// Initialise raylib
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowResizable)
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "CHIP8 Interpreter")
	defer rl.CloseWindow()

	// Set audio buffer
	rl.InitAudioDevice()
	rl.SetAudioStreamBufferSizeDefault(AUDIO_BUFFER_SIZE)
	stream := rl.LoadAudioStream(SAMPLE_RATE, 32, 1)
	rl.PlayAudioStream(stream)
	rl.SetAudioStreamCallback(stream, audioCallback(emu))
	defer rl.UnloadAudioStream(stream)
	defer rl.CloseAudioDevice()

	// Create texture for emulator display
	canvas := rl.LoadRenderTexture(int32(chip8.DISPLAY_X), int32(chip8.DISPLAY_Y))
	defer rl.UnloadRenderTexture(canvas)

	// Main loop
	rl.SetTargetFPS(int32(FPS))
	for !rl.WindowShouldClose() {
		// CPU
		emu.Step()

		// AUDIO
		if emu.GetSoundTimer() > 0 {
		}

		// DISPLAY
		drawOnTexture(canvas, emu)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Resize and render texture
		src := rl.NewRectangle(0, 0, float32(chip8.DISPLAY_X), -float32(chip8.DISPLAY_Y))
		dst := rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()))
		rl.DrawTexturePro(canvas.Texture, src, dst, rl.NewVector2(0, 0), 0, rl.White)

		rl.EndDrawing()
	}
}
