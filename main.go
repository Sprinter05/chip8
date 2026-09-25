package main

import (
	"flag"
	"math"
	"os"
	"time"

	"github.com/Sprinter05/chip-8/chip8"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const FPS float64 = 60.0
const FRAMETIME_US float64 = 1000000.0 / FPS // microseconds
const WINDOW_WIDTH = 640
const WINDOW_HEIGHT = 320

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
	emu.SetInputCallback(inputCallback, inputListCallback)
	if err := load(emu); err != nil {
		panic(err)
	}

	// Initialise raylib
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowResizable)
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "CHIP8 Interpreter")
	defer rl.CloseWindow()

	// Create texture for emulator display
	canvas := rl.LoadRenderTexture(int32(chip8.DISPLAY_X), int32(chip8.DISPLAY_Y))
	defer rl.UnloadRenderTexture(canvas)

	// Create ticker for waiting when stepping
	dur := int64(math.Round(FRAMETIME_US))
	ticker := time.NewTicker(time.Duration(dur) * time.Microsecond)
	defer ticker.Stop()

	// Main loop
	rl.SetTargetFPS(int32(FPS))
	for !rl.WindowShouldClose() {
		// CPU
		emu.Step()

		// DISPLAY
		drawOnTexture(canvas, emu)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Resize and render texture
		src := rl.NewRectangle(0, 0, float32(chip8.DISPLAY_X), -float32(chip8.DISPLAY_Y))
		dst := rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()))
		rl.DrawTexturePro(canvas.Texture, src, dst, rl.NewVector2(0, 0), 0, rl.White)

		rl.EndDrawing()

		// WAIT
		<-ticker.C
	}
}
