package lalStore

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/cloud/internal/settings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/fynelabs/lal"
)

type lalStore struct {
	store  fyne.URI
	config string
	prefs  fyne.Preferences
	db     *lal.DB
}

func (d *lalStore) SetConfig(str string) {
	d.config = str
}

func (d *lalStore) Configure(a fyne.App, w fyne.Window) (data string, err error) {

	fmt.Println("configuring kv store")

	if fyne.CurrentDevice().IsMobile() {
		data, err = d.mobileConfig(a)
	} else {
		data, err = d.desktopConfig(w)
	}

	if err == nil {
		d.config = data
	}
	return data, err
}

func (d *lalStore) ProviderDescription() string {
	return "KV Store description...\n"
}

func (d *lalStore) ProviderIcon() fyne.Resource {
	return nil
}

func (d *lalStore) ProviderName() string {
	return "KV Store"
}

func (d *lalStore) Setup(a fyne.App) error {
	dbURI, err := storage.Child(a.Storage().RootURI(), "preference.laldb")
	if err != nil {
		fyne.LogError("Can't setup prefs", err)
	}

	d.db, _ = lal.Open(dbURI.String(), 0666, lal.DefaultOptions)

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

	fmt.Println("setup")

	return err
}

func (d *lalStore) desktopConfig(w fyne.Window) (string, error) {
	err := make(chan error)
	data := ""

	ask := dialog.NewConfirm("Set lalStore location",
		"On the next screen please choose the KV DB",
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

func (d *lalStore) mobileConfig(a fyne.App) (string, error) {
	err := make(chan error)
	data := ""
	for len(a.Driver().AllWindows()) == 0 {
		time.Sleep(time.Millisecond * 100)
	}
	win := a.Driver().AllWindows()[0]

	ask := dialog.NewConfirm("Locate kvstore files",
		"On the next screen please load the kvstore file\nfynesync/"+a.UniqueID()+"/preferences.json",
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
	return &lalStore{}
}

// tempfile returns a temporary file path.
func tempfile() string {
	f, err := os.CreateTemp("", "bolt-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
