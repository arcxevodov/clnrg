package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	NoRootError  = 1
	UnknownError = 2
)

var localeArg = flag.String("lang", "en", "Set locale")

func main() {
	flag.Parse()
	a := app.New()

	_, err := exec.Command("clnr", "-i").Output()
	if err != nil {
		if err.Error() == "exec: \"clnr\": executable file not found in $PATH" {
			noClnr(a).Show()
		} else if err.Error() == "exit status 1" {
			fmt.Println("\033[0;31m" + localString("NoRoot"))
			os.Exit(NoRootError)
		} else {
			fmt.Println(localString("UnknownError"), err)
			os.Exit(UnknownError)
		}
	} else {
		mainWidget(a).Show()
	}

	a.Run()
}

func initLocalizer() *i18n.Localizer {
	var bundle *i18n.Bundle
	var localizer *i18n.Localizer

	path, err := os.Executable()
	if err != nil {
		fmt.Println(localString("UnkownError"), err)
		os.Exit(UnknownError)
	}

	if *localeArg == "ru" {
		bundle = i18n.NewBundle(language.Russian)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		_, err := bundle.LoadMessageFile(path[:len(path)-5] + "/locales/ru.json")
		if err != nil {
			fmt.Println(localString("UnknownError"), err)
			os.Exit(UnknownError)
		}
		localizer = i18n.NewLocalizer(bundle, language.Russian.String())
	} else {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		_, err := bundle.LoadMessageFile(path[:len(path)-5] + "/locales/en.json")
		if err != nil {
			fmt.Println(localString("UnknownError"), err)
			os.Exit(UnknownError)
		}
		localizer = i18n.NewLocalizer(bundle, language.English.String())
	}
	return localizer
}

func localString(id string) string {
	localizer := initLocalizer()
	localzeConfig := i18n.LocalizeConfig{
		MessageID: id,
	}
	result, err := localizer.Localize(&localzeConfig)
	if err != nil {
		fmt.Println(localString("UnknownError"), err)
		os.Exit(UnknownError)
	}
	return result
}

func updateRam(ram *widget.Label) {
	var out []byte
	var err error
	if *localeArg == "ru" {
		out, err = exec.Command("clnr", "-i", "-lang=ru").Output()
	} else {
		out, err = exec.Command("clnr", "-i").Output()
	}
	if err != nil {
		fmt.Println(localString("UnknownError"), err)
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

	clean := widget.NewButton(localString("CleanRamButton"), func() {
		var info []byte
		var err error
		if *localeArg == "ru" {
			info, err = exec.Command("clnr", "-r", "-s", "-lang=ru").Output()
		} else {
			info, err = exec.Command("clnr", "-r", "-s").Output()
		}
		if err != nil {
			fmt.Println(localString("UnknownError"), err)
			os.Exit(UnknownError)
		}
		cleanSucsessful(a, string(info)).Show()
	})

	content := container.New(layout.NewVBoxLayout(), ram, layout.NewSpacer(), clean)
	w.SetContent(container.New(layout.NewVBoxLayout(), content))
	return w
}

func noClnr(a fyne.App) fyne.Window {
	w := a.NewWindow(localString("NoClnrTitle"))
	prompt := widget.NewLabel(localString("ClnrNotFound"))
	w.SetContent(prompt)
	return w
}

func cleanSucsessful(a fyne.App, out string) fyne.Window {
	w := a.NewWindow(localString("CleanSuccessfulTitle"))
	logging := widget.NewLabel(out)
	w.SetContent(logging)
	return w
}
