package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
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
		configDir = filepath.Join(os.Getenv("Home"), ".config")
	}

	appDir := filepath.Join(configDir, "arlopeeker")
	os.MkdirAll(appDir, os.ModePerm)

	return filepath.Join(appDir, "config.json")
}

//go:embed assets/icon.png
var assetsFS embed.FS

var mainThreadCh = make(chan func())

func callOnMainThread(f func()) {
	done := make(chan struct{})
	mainThreadCh <- func() {
		f()
		close(done)
	}
	<-done
}

func loadConfig() {
	configPath = getConfigPath()
	configFile, err := os.Open(configPath)

	if err != nil {
		config = Config{
			Image:    "assets/icon.png",
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
			Image:    "assets/icon.png",
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
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	go func() {
		for {
			<-mSetting.ClickedCh
			callOnMainThread(func() {

				imagePath, err := dialog.File().
					Filter("Image Files", "png", "jpg", "jpeg").
					Title("Select Pet Image").
					Load()

				if err == nil && imagePath != "" {
					config.Image = imagePath
					saveConfig()
				} else if err != nil {
					log.Println("Image selection canceled or failed:", err)
				}
			})
		}
	}()
}

func onExit() {
	now := time.Now()
	if _, err := io.WriteString(os.Stdout, fmt.Sprintf("Arlo Peeker exited at %s\n", now.String())); err != nil {
		log.Fatal(err)
	}

}

func main() {
	loadConfig()

	go func() {
		for f := range mainThreadCh {
			f()
		}
	}()

	systray.Run(onReady, onExit)
}
