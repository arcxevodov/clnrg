package main

// #include <unistd.h>
import "C"

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func updateRam(ram *widget.Label) {
	bTotal := C.sysconf(C._SC_PHYS_PAGES) * C.sysconf(C._SC_PAGE_SIZE)
	gbTotal := float64(bTotal) / 1024 / 1024 / 1024
	fmtTotal := fmt.Sprintf("Total RAM: %.1f GB", gbTotal)

	bFree := C.sysconf(C._SC_AVPHYS_PAGES) * C.sysconf(C._SC_PAGE_SIZE)
	gbFree := float64(bFree) / 1024 / 1024 / 1024
	fmtFree := fmt.Sprintf("Available RAM: %.1f GB", gbFree)

	ram.SetText(fmt.Sprintf("%s | %s", fmtTotal, fmtFree))
}

func main() {
	a := app.New()
	w := a.NewWindow("clnr")

	ram := widget.NewLabel("")
	updateRam(ram)
	go func() {
		for range time.Tick(time.Second) {
			updateRam(ram)
		}
	}()

	clean := widget.NewButton("Clean RAM", func() {
		out, err := exec.Command("clnr", "-r", "-s").Output()
		if err != nil {
			if err.Error() == "exec: \"clnr\": executable file not found in $PATH" {
				installPrompt(a).Show()
			}
		}
		cleanSucsessful(a, string(out)).Show()
	})

	content := container.New(layout.NewVBoxLayout(), ram, layout.NewSpacer(), clean)
	w.SetContent(container.New(layout.NewVBoxLayout(), content))
	w.Show()

	a.Run()
}

func installPrompt(a fyne.App) fyne.Window {
	w := a.NewWindow("clnr not found")
	prompt := widget.NewLabel("\"clnr\" executable file not found. Would you install this?")
	yButton := widget.NewButton("Yes", func() {
		cmd := "curl -fsSL https://github.com/arcxevodov/clnr/releases/download/v0.1.0/clnr > /usr/bin/clnr; chmod +x /usr/bin/clnr"
		err := exec.Command("bash", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
		}
		w.Close()
	})
	nButton := widget.NewButton("No", func() {
		w.Close()
	})
	content := container.New(layout.NewVBoxLayout(), prompt, layout.NewSpacer(), yButton, nButton)
	w.SetContent(container.New(layout.NewVBoxLayout(), content))
	return w
}

func cleanSucsessful(a fyne.App, out string) fyne.Window {
	w := a.NewWindow("Cleaning...")
	logging := widget.NewLabel(out)
	w.SetContent(logging)
	return w
}
