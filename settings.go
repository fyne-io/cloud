package cloud

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowSettings(a fyne.App, w fyne.Window) {
	var selected fyne.CloudProvider
	ui := widget.NewList(func() int {
		return len(providers)
	},
	func() fyne.CanvasObject {
		return container.NewHBox(
			widget.NewIcon(nil),
			widget.NewLabel(""))
	},
	func(id widget.ListItemID, o fyne.CanvasObject) {
		p := providers[id]
		o.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(p.ProviderIcon())
		o.(*fyne.Container).Objects[1].(*widget.Label).SetText(p.ProviderName())
	})
	ui.OnSelected = func(id widget.ListItemID) {
		selected = providers[id]
	}

	dialog.ShowCustomConfirm("Choose Cloud provider", "Enable", "Cancel",
		container.NewVBox(widget.NewLabel("Choose one of the providers below to sync\nyour preferences to the cloud."), ui),
		func(ok bool) {
		if !ok || selected == nil {
			return
		}

		chooseProvider(a, selected)
	}, w)
}

func chooseProvider(a fyne.App, p fyne.CloudProvider) {
	setCurrentProviderName(p.ProviderName())
	a.SetCloudProvider(p)
}
