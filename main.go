package main

import (
	"fmt"
	"time"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New() // Create a new app
	myWindow := myApp.NewWindow("Stopwatch")
	myWindow.Resize(fyne.NewSize(700, 500)) // Adjust window size

	stopwatch := widget.NewLabel("00:00.00")
	stopwatch.TextStyle.Bold = true

	// Convert the RGBA values (0-255) to fit the 0-1 range
	red := uint8(2)
	green := uint8(132)
	blue := uint8(199)
	alpha := uint8(255) // 1 in the range of 0-1 becomes 255 in the 0-255 range

	textColor := color.RGBA{red, green, blue, alpha}

	infoText := canvas.NewText("Stopwatch", textColor)
	infoText.TextSize = 30 // Adjust text size

	var (
		startTime time.Time
		elapsed time.Duration
		running bool
	)

	var startStopBtn *widget.Button

	// Add the start/pause functionality
	startStopBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if !running {
			startTime = time.Now().Add(-elapsed)
			running = true
			startStopBtn.SetIcon(theme.MediaPauseIcon())
			// Run the stopwatch as a goroutine
			go runStopwatch(stopwatch, &startTime, &running, &elapsed)
		} else {
			elapsed = time.Since(startTime)
			running = false
			startStopBtn.SetIcon(theme.MediaPlayIcon())
		}
	})

	// Add the reset functionality
	resetBtn := widget.NewButtonWithIcon("", theme.MediaReplayIcon(), func() {
		stopwatch.SetText("00:00.00")
		elapsed = 0
		running = false
		startStopBtn.SetIcon(theme.MediaPlayIcon())
	})

	// Create a container for the buttons
	buttons := container.NewHBox(layout.NewSpacer(), resetBtn, startStopBtn, layout.NewSpacer())
	// Create a container for the stopwatch
	counterContainer := container.NewVBox(
		layout.NewSpacer(),
		container.New(layout.NewCenterLayout(), stopwatch),
		layout.NewSpacer(),
	)

	// Create a container for topmost text and place it above the counter
	infoContainer := container.NewVBox(
		layout.NewSpacer(),
		container.New(layout.NewCenterLayout(), infoText),
		layout.NewSpacer(),
	)

	// Add the containers to the window and add spacers appropriately
	myWindow.SetContent(container.NewVBox(infoContainer, layout.NewSpacer(), counterContainer, layout.NewSpacer(), buttons,))
	// run the window
	myWindow.ShowAndRun()
}

// runStopwatch runs the stopwatch and updates the label
func runStopwatch(stopwatch *widget.Label, startTime *time.Time, running *bool, elapsed *time.Duration) {
	for *running {
		*elapsed = time.Since(*startTime)
		centiseconds := (*elapsed).Milliseconds() / 10 // Centiseconds
		minutes := centiseconds / (60 * 100)
		centiseconds -= minutes * (60 * 100)
		seconds := centiseconds / 100
		centiseconds -= seconds * 100

		// Format centiseconds without trailing zeros
		centisecondsString := fmt.Sprintf("%02d", centiseconds)
		
		stopwatch.SetText(fmt.Sprintf("%02d:%02d.%s", minutes, seconds, centisecondsString))

		// Sleep for 10 milliseconds before running again
		time.Sleep(10 * time.Millisecond)
	}
}
