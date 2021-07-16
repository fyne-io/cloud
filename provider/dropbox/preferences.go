package dropbox

import (
	"os"
	"path/filepath"

	"fyne.io/cloud/internal"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func (d *dropbox) CloudPreferences(a fyne.App) fyne.Preferences {
	if d.store != nil {
		return internal.NewPreferences(a, d.store)
	}

	filePath := filepath.Join(dropboxLocation(), a.UniqueID(), "preferences.json")
	return internal.NewPreferences(a, storage.NewFileURI(filePath))
}

func dropboxLocation() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Dropbox", "fynesync")
}
