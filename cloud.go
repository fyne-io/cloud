// Package cloud provides various cloud provider implementations and utilities to add the cloud services
// into your Fyne app. Developers can choose to load a specific provider, or they can present a configuration
// user interface allowing the end-user to choose the provider they wish to use.
//
// A simple usage where an app uses AWS for cloud provision may look like this:
//
//   package main
//
//   import (
//   	"fyne.io/cloud/provider/aws"
//   	"fyne.io/fyne/v2/app"
//   	"fyne.io/fyne/v2/widget"
//   )
//
//   func main() {
//   	a := app.New()
//   	a.SetCloudProvider(aws.NewProvider()) // if aws provider existed ;)
//   	w := a.NewWindow("Cloud")
//
//   	w.SetContent(widget.NewLabel("Add content here"))
//   	w.ShowAndRun()
//   }
//
// Alternatively to allow the user to choose a cloud provider for their storage etc use:
//
//   package main
//
//   import (
//   	"fyne.io/cloud"
//   	"fyne.io/fyne/v2/app"
//   	"fyne.io/fyne/v2/widget"
//   )
//
//   func main() {
//   	a := app.New()
//   	cloud.Enable(a)
//   	w := a.NewWindow("Cloud")
//
//   	w.SetContent(widget.NewButton("Choose cloud provider", func() {
//   		cloud.ShowSettings(a, w)
//   	}))
//   	w.ShowAndRun()
//   }
package cloud // import "fyne.io/cloud"

import "fyne.io/fyne/v2"

var providers []fyne.CloudProvider

func Enable(a fyne.App) {
	a.SetCloudProvider(lookupConfiguredProvider())
}

func Register(p fyne.CloudProvider) {
	providers = append(providers, p)
}

func lookupConfiguredProvider() fyne.CloudProvider {
	name := currentProviderName()

	for _, p := range providers {
		if p.ProviderName() == name {
			return p
		}
	}
	return nil
}
