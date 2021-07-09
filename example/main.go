package main

import (
	"fmt"

	"fyne.io/cloud"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	cloud.Enable(a)
	w := a.NewWindow("Cloud")

	current := widget.NewLabel("No configured provider")
	choose := widget.NewButton("Choose cloud provider", func() {
		cloud.ShowSettings(a, w)
	})

	updateCloudName := func () {
		current.SetText(fmt.Sprintf("Using %s provider", a.CloudProvider().ProviderName()))
		choose.SetText("Change provider")
	}
	if a.CloudProvider() != nil {
		updateCloudName()
	}

	ch := make(chan fyne.Settings)
	a.Settings().AddChangeListener(ch)
	go func() {
		for range ch {
			updateCloudName()
		}
	}()

	w.SetContent(container.NewVBox(current, choose))
	w.ShowAndRun()
}
