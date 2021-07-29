//go:generate fyne bundle -o bundled.go -package dropbox Dropbox.svg

package dropbox

import (
	"errors"
	"os"
	"strings"
	"time"

	"fyne.io/cloud/internal/settings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
)

type dropbox struct {
	store  fyne.URI
	config string
	prefs  fyne.Preferences
}

func (d *dropbox) SetConfig(str string) {
	d.config = str
}

func (d *dropbox) Configure(w fyne.Window) (data string, err error) {
	if fyne.CurrentDevice().IsMobile() {
		data, err = d.mobileConfig(fyne.CurrentApp())
	} else {
		data, err = d.desktopConfig(w)
	}

	if err == nil {
		d.config = data
	}
	return data, err
}

func (d *dropbox) Disconnect() {
	d.prefs.(interface{ Disconnect() }).Disconnect()
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
	if d.config != "" {
		d.store = storage.NewFileURI(d.config)
		// TODO validate
		return nil
	}
	data, err := d.mobileConfig(a)
	settings.SetProviderConfig(data)

	return err
}

func (d *dropbox) desktopConfig(w fyne.Window) (string, error) {
	err := make(chan error)
	data := ""

	ask := dialog.NewConfirm("Set Dropbox location",
		"On the next screen please choose the the fynesync\nfolder, or whichever you would prefer",
		func(ok bool) {
			if !ok {
				err <- errors.New("user cancelled setup")
				return
			}

			dialog.ShowFolderOpen(func(list fyne.ListableURI, e2 error) {
				if e2 != nil {
					err <- e2
					return
				}

				if list == nil {
					err <- errors.New("no folder was chosen")
					return
				}

				// TODO check content?...
				rootDir := list.Path()
				home, _ := os.UserHomeDir()
				if len(rootDir) > len(home) && strings.Index(rootDir, home) == 0 {
					rootDir = "~" + rootDir[len(home):]
				}

				data = rootDir
				err <- nil
			}, w)
		}, w)

	ask.SetConfirmText("OK")
	ask.SetDismissText("Cancel")
	ask.Show()
	return data, <-err
}

func (d *dropbox) mobileConfig(a fyne.App) (string, error) {
	err := make(chan error)
	data := ""
	for len(fyne.CurrentApp().Driver().AllWindows()) == 0 {
		time.Sleep(time.Millisecond * 100)
	}
	win := fyne.CurrentApp().Driver().AllWindows()[0]

	ask := dialog.NewConfirm("Locate Dropbox files",
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
				data = d.store.Path()
				err <- nil
			}, win)
		}, win)

	ask.SetConfirmText("OK")
	ask.SetDismissText("Cancel")
	ask.Show()
	return data, <-err
}

func NewProvider() fyne.CloudProvider {
	return &dropbox{}
}
