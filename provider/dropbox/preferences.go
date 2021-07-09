package dropbox

import (
	"os"
	"path/filepath"

	"fyne.io/cloud/internal"
	"fyne.io/fyne/v2"
)

func (d *dropbox) CloudPreferences(a fyne.App) fyne.Preferences {
	return internal.NewPreferences(a, dropboxLocation())
}

func dropboxLocation() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Dropbox", ".fyne")
}
