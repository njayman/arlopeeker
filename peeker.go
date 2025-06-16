package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ShowPeeker(imagePath string, duration time.Duration, speed float32) {
	rl.SetConfigFlags(rl.FlagWindowUndecorated | rl.FlagWindowTransparent | rl.FlagWindowAlwaysRun)
	rl.InitWindow(300, 300, "")
	rl.SetTargetFPS(60)

	// image := rl.LoadImage(im
	// defer rl.UnloadImage(image)
	//
	// if image.Width == 0 {
	// 	panic("Failed to load image")
	// }

	texture := rl.LoadTexture(imagePath)

	if texture.ID == 0 {
		rl.CloseWindow()
		return
	}

	defer func() {
		rl.UnloadTexture(texture)
		rl.CloseWindow()

	}()

	start := time.Now()
	peekX := float32(-float32(texture.Width))
	targetX := float32(0)
	step := speed * 5

	for !rl.WindowShouldClose() {
		elapsed := time.Since(start)

		if elapsed < duration/3 && peekX < targetX {
			peekX += step
		} else if elapsed > 2*duration/3 && peekX > -float32(texture.Width) {
			peekX -= step
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Blank)
		rl.DrawTexture(texture, int32(peekX), 100, rl.White)
		rl.EndDrawing()

		if elapsed > duration {
			break
		}
	}
}
