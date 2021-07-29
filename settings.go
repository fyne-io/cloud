package cloud

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowSettings(a fyne.App, w fyne.Window) {
	prov := a.CloudProvider()
	if prov == nil {
		showChoice(a, w)
	}

	config := widget.NewButton("Configure", func() {
		go func() {
			str, err := prov.(Configurable).Configure(w)
			if err != nil {
				dialog.ShowError(err, w)
			} else {
				if dis, ok := a.CloudProvider().(Disconnectable); ok {
					dis.Disconnect()
				}

				setProviderConfig(str)
				prov.(Configurable).SetConfig(str)
				a.SetCloudProvider(prov) // we don't use setCloud here as we have a new config
			}
		}()
	})
	if _, ok := prov.(Configurable); !ok {
		config.Hide()
	}
	current := widget.NewLabel(fmt.Sprintf("Using %s provider", prov.ProviderName()))
	ch := make(chan fyne.Settings)
	a.Settings().AddChangeListener(ch)
	go func() {
		for range ch {
			current.SetText(fmt.Sprintf("Using %s provider", a.CloudProvider().ProviderName()))
		}
	}()

	var d dialog.Dialog
	d = dialog.NewCustomConfirm("Cloud configuration", "Change Provider", "Cancel",
		container.NewBorder(nil, nil, nil, config, current),
		func(change bool) {
			if !change {
				return
			}

			d.Hide()
			showChoice(a, w)
		}, w)
	d.Show()
}

func showChoice(a fyne.App, w fyne.Window) {
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
	setCloud(p, a)
}
