package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	NoRootError  = 1
	UnknownError = 2
)

func main() {
	a := app.New()

	_, err := exec.Command("clnr", "-i").Output()
	if err != nil {
		if err.Error() == "exec: \"clnr\": executable file not found in $PATH" {
			noClnr(a).Show()
		} else if err.Error() == "exit status 1" {
			fmt.Println("\033[0;31mPlease run this program as superuser")
			os.Exit(NoRootError)
		} else {
			fmt.Println("Unknown Error: ", err)
			os.Exit(UnknownError)
		}
	} else {
		mainWidget(a).Show()
	}

	a.Run()
}

func updateRam(ram *widget.Label) {
	out, err := exec.Command("clnr", "-i").Output()
	if err != nil {
		fmt.Println("Unknown Error: ", err)
		os.Exit(UnknownError)
	}
	ram.SetText(string(out))
}

func mainWidget(a fyne.App) fyne.Window {
	w := a.NewWindow("clnr")
	ram := widget.NewLabel("")

	updateRam(ram)
	go func() {
		for range time.Tick(time.Second) {
			updateRam(ram)
		}
	}()

	clean := widget.NewButton("Clean RAM", func() {
		info, err := exec.Command("clnr", "-r", "-s").Output()
		if err != nil {
			fmt.Println("Unknown Error: ", err)
			os.Exit(UnknownError)
		}
		cleanSucsessful(a, string(info)).Show()
	})

	content := container.New(layout.NewVBoxLayout(), ram, layout.NewSpacer(), clean)
	w.SetContent(container.New(layout.NewVBoxLayout(), content))
	return w
}

func noClnr(a fyne.App) fyne.Window {
	w := a.NewWindow("clnr not found")
	prompt := widget.NewLabel("\"clnr\" executable file not found. Please run install_clnr.sh")
	w.SetContent(prompt)
	return w
}

func cleanSucsessful(a fyne.App, out string) fyne.Window {
	w := a.NewWindow("Cleaning...")
	logging := widget.NewLabel(out)
	w.SetContent(logging)
	return w
}
