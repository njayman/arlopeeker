package main

import (
	"log"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ShowPeeker(imagePath string, duration time.Duration, speed float32) {
	rl.SetConfigFlags(rl.FlagWindowUndecorated | rl.FlagWindowTransparent | rl.FlagWindowAlwaysRun)
	rl.InitWindow(600, 600, "")
	rl.SetTargetFPS(60)

	texture := rl.LoadTexture(imagePath)

	maxSize := float32(600)
	scale := float32(1.0)

	if texture.Width > int32(maxSize) || texture.Height > int32(maxSize) {
		wRatio := maxSize / float32(texture.Width)
		hRatio := maxSize / float32(texture.Height)
		scale = float32(math.Min(float64(wRatio), float64(hRatio)))
	}

	scaledWidth := int32(float32(texture.Width) * scale)
	scaledHeight := int32(float32(texture.Height) * scale)

	if texture.ID == 0 {
		rl.CloseWindow()
		log.Printf("Failed to load image: %s\n", imagePath)
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
		rl.DrawTexturePro(
			texture,
			rl.NewRectangle(0, 0, float32(texture.Width), float32(texture.Height)),
			rl.NewRectangle(peekX, 100, float32(scaledWidth), float32(scaledHeight)),
			rl.NewVector2(0, 0),
			0,
			rl.White,
		)
		rl.EndDrawing()

		if elapsed > duration {
			break
		}
	}
}
