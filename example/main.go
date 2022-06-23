package main

import (
	"fmt"

	"fyne.io/cloud"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("io.fyne.cloud.example")
	w := a.NewWindow("Cloud")

	current := widget.NewLabel("No configured provider")
	ch := make(chan fyne.Settings)
	a.Settings().AddChangeListener(ch)
	cloud.Enable(a)

	choose := widget.NewButton("Cloud settings", func() {
		cloud.ShowSettings(a, w)
	})
	testEntry := widget.NewEntryWithData(binding.BindPreferenceString("test", a.Preferences()))
	testEntry.Validator = nil

	go func() {
		for range ch {
			if a.CloudProvider() == nil {
				return
			}
			current.SetText(fmt.Sprint("Current cloud: ", a.CloudProvider().ProviderName()))
		}
	}()

	w.SetContent(container.NewVBox(layout.NewSpacer(), current,
		container.NewBorder(nil, nil, widget.NewLabel("Test store value"), nil, testEntry),
		layout.NewSpacer(), choose))
	w.ShowAndRun()
}
