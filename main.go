package main

import (
	"flag"
	"os"

	"github.com/Sprinter05/chip-8/chip8"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const FPS float64 = 60.0
const FRAMETIME_US float64 = 1000000.0 / FPS // microseconds
const SIZE = 10                              // size for pixels

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
	emu := new(chip8.CHIP8)
	emu.Reset()
	if err := load(emu); err != nil {
		panic(err)
	}

	rl.InitWindow(640, 320, "CHIP8 Interpreter")
	defer rl.CloseWindow()

	rl.SetTargetFPS(int32(FPS))
	rl.SetConfigFlags(rl.FlagVsyncHint)

	// dur := int64(math.Round(FRAMETIME_US))
	// ticker := time.NewTicker(time.Duration(dur) * time.Microsecond)

	for !rl.WindowShouldClose() {
		// CPU
		emu.Step()

		// DISPLAY
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		for x := range emu.Display {
			for y, v := range emu.Display[x] {
				if v {
					rl.DrawRectangle(int32(x), int32(y), SIZE, SIZE, rl.White)
				}
			}
		}

		rl.EndDrawing()

		// WAIT
		// <-ticker.C
	}
}
