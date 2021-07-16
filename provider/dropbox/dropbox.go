//go:generate fyne bundle -o bundled.go -package dropbox Dropbox.svg

package dropbox

import (
	"errors"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
)

type dropbox struct {
	store fyne.URI
}

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

func (d *dropbox) Setup(a fyne.App) error {
	if !fyne.CurrentDevice().IsMobile() {
		return nil
	}

	err := make(chan error)
	for len(fyne.CurrentApp().Driver().AllWindows()) == 0 {
		time.Sleep(time.Millisecond*100)
	}
	win := fyne.CurrentApp().Driver().AllWindows()[0]

	ask := dialog.NewConfirm("Locate dropbox",
		"On the next screen please load the Dropbox file\nfynesync/"+a.UniqueID()+"/preferences.json",
		func(ok bool) {
			if !ok {
				err <- errors.New("user cancelled setup")
				return
			}

			dialog.ShowFileOpen(func(read fyne.URIReadCloser, e2 error) {
				if e2 != nil {
					err <- e2
					return
				}

				if read == nil {
					err <- errors.New("no file was chosen")
					return
				}

				// TODO check content?...
				d.store = read.URI()
				err <- nil
			}, win)
		}, win)

	ask.SetConfirmText("OK")
	ask.SetDismissText("Cancel")
	ask.Show()
	return <- err
}

func NewProvider() fyne.CloudProvider {
	return &dropbox{}
}

