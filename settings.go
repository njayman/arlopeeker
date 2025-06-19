package main

import (
	"fmt"
	"log"
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sqweek/dialog"
)

func ShowSettingsWindow() {
	loadConfig()
	rl.InitWindow(500, 300, "ArloPeeker Settings")
	defer rl.CloseWindow()

	textImagePath := config.Image
	textDuration := fmt.Sprintf("%.2f", config.Duration)
	textSpeed := fmt.Sprintf("%.2f", config.Speed)

	marginTop := float32(25)
	labelX := float32(20)
	fieldX := float32(140)
	fieldWidth := float32(240)
	fieldHeight := float32(30)
	spacingY := float32(50)

	imageY := marginTop
	durationY := imageY + spacingY
	speedY := durationY + spacingY

	editImagePath := true
	editDuration := false
	editSpeed := false

	buttonSavePressed := false

	for !rl.WindowShouldClose() && !buttonSavePressed {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		mouse := rl.GetMousePosition()

		rl.DrawText("Image Path:", int32(labelX), int32(imageY+8), 20, rl.DarkGray)
		rl.DrawText("Duration (sec):", int32(labelX), int32(durationY+8), 20, rl.DarkGray)
		rl.DrawText("Speed:", int32(labelX), int32(speedY+8), 20, rl.Black)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			editImagePath = rl.CheckCollisionPointRec(mouse, rl.NewRectangle(fieldX, imageY, fieldWidth, fieldHeight))
			editDuration = rl.CheckCollisionPointRec(mouse, rl.NewRectangle(fieldX, durationY, fieldWidth, fieldHeight))
			editSpeed = rl.CheckCollisionPointRec(mouse, rl.NewRectangle(fieldX, speedY, fieldWidth, fieldHeight))

			if editImagePath {
				editDuration = false
				editSpeed = false
			} else if editDuration {
				editImagePath = false
				editSpeed = false
			} else if editSpeed {
				editImagePath = false
				editDuration = false
			}
		}

		rg.TextBox(rl.NewRectangle(fieldX, imageY, fieldWidth, fieldHeight), &textImagePath, 512, editImagePath)
		rg.TextBox(rl.NewRectangle(fieldX, durationY, fieldWidth, fieldHeight), &textDuration, 64, editDuration)
		rg.TextBox(rl.NewRectangle(fieldX, speedY, fieldWidth, fieldHeight), &textSpeed, 64, editSpeed)

		if rg.Button(rl.NewRectangle(fieldX+fieldWidth+10, imageY, 90, fieldHeight), "Browse") {
			path, err := dialog.File().Title("Select Imahe").Filter("Images", "png", "jpg", "jpeg", "bmp").Load()

			if err == nil {
				textImagePath = path
			} else {
				log.Println("File dialog error:", err)
			}
		}

		if rg.Button(rl.NewRectangle(200, 200, 100, 40), "Save") {
			buttonSavePressed = true
		}

		rl.EndDrawing()
	}

	if buttonSavePressed {
		dur, err1 := strconv.ParseFloat(textDuration, 64)
		spd, err2 := strconv.ParseFloat(textSpeed, 64)

		if err1 == nil && err2 == nil {
			config.Image = textImagePath
			config.Duration = dur
			config.Speed = spd
			saveConfig()
		} else {
			log.Println("Invalid Input: duration or speed could not be parsed")
		}
	}
}
