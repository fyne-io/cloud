//go:generate fyne bundle -o bundled.go -package dropbox Dropbox.svg

package dropbox

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type dropbox struct {}

func (d *dropbox) ProviderDescription() string {
	return "Preferences and files are stored into the user's Dropbox storage.\n" +
		"You need to have Dropbox installed and running for this provider to work."
}

func (d *dropbox) ProviderIcon() fyne.Resource {
	return theme.NewThemedResource(resourceDropboxSvg)
}

func (d *dropbox) ProviderName() string {
	return "Dropbox"
}

func NewProvider() fyne.CloudProvider {
	return &dropbox{}
}

