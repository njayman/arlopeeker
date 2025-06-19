package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/getlantern/systray"
)

type Config struct {
	Image    string  `json:"image"`
	Duration float64 `json:"duration"`
	Speed    float64 `json:"speed"`
}

var config Config
var configPath string

func getConfigPath() string {
	var configDir string
	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("AppData")
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	default:
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	appDir := filepath.Join(configDir, "arlopeeker")
	os.MkdirAll(appDir, os.ModePerm)

	return filepath.Join(appDir, "config.json")
}

//go:embed assets/icon.png
var assetsFS embed.FS

func loadConfig() {
	configPath = getConfigPath()
	configFile, err := os.Open(configPath)

	if err != nil {
		config = Config{
			Image:    "assets/photo.png",
			Duration: 3,
			Speed:    1.5,
		}

		return
	}

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)

	if err := decoder.Decode(&config); err != nil {
		log.Printf("Failed to parse config.json: %v. Using default config.\n", err)
		config = Config{
			Image:    "assets/photo.png",
			Duration: 3,
			Speed:    1.5,
		}
	}
}

func saveConfig() {
	file, err := os.Create(configPath)

	if err != nil {
		log.Println("Failed to save config:", err)

		return
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "")

	if err := encoder.Encode(config); err != nil {
		log.Println("Failed to encode config:", err)
	}
}

func onReady() {
	iconBytes, err := assetsFS.ReadFile("assets/icon.png")

	if err != nil {
		log.Fatalf("Failed to load systray icon: %v", err)
	}

	systray.SetIcon(iconBytes)
	systray.SetTitle("ArloPeeker")
	systray.SetTooltip("Arlo Peeker")

	mShow := systray.AddMenuItem("Show", "Configure settings")
	mSetting := systray.AddMenuItem("Settings", "Configure settings")
	mQuit := systray.AddMenuItem("Quit", "Quit arlo peeker")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				go ShowPeeker(config.Image, time.Duration(config.Duration*float64(time.Second)), float32(config.Speed))
			case <-mSetting.ClickedCh:
				uiActions <- func() {
					ShowSettingsWindow()
				}
				// fmt.Println("Settings")

			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	now := time.Now()
	if _, err := io.WriteString(os.Stdout, fmt.Sprintf("Arlo Peeker exited at %s\n", now.String())); err != nil {
		log.Fatal(err)
	}

}

var uiActions = make(chan func())

func main() {
	runtime.LockOSThread()

	flagSettings := flag.Bool("settings", false, "Open settings window")
	cliMode := flag.Bool("peek", false, "Trigger peek image once and exit")
	image := flag.String("image", "", "Path to image to peek (optional)")
	duration := flag.Float64("duration", 0, "Duration in seconds (optional)")
	speed := flag.Float64("speed", 0, "Speed factor (optional)")

	flag.Parse()

	loadConfig()

	if *cliMode {
		img := config.Image

		if *image != "" {
			img = *image
		}

		dur := time.Duration(config.Duration * float64(time.Second))

		if *duration > 0 {
			dur = time.Duration(*duration * float64(time.Second))
		}

		spd := float32(config.Speed)

		if *speed > 0 {
			spd = float32(*speed)
		}

		if *flagSettings {
			ShowSettingsWindow()
			return
		}

		ShowPeeker(img, dur, spd)

		os.Exit(0)
	}

	go func() {
		systray.Run(onReady, onExit)
	}()

	for action := range uiActions {
		action()
	}
}
