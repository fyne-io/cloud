package lalStore

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func (d *lalStore) CloudPreferences(a fyne.App) fyne.Preferences {
	var p fyne.Preferences
	if d.store != nil {
		p = NewPreferences(a, d.store, d.db)
	} else {
		filePath := filepath.Join(d.dbLocation(), a.UniqueID(), "preferences.json")
		p = NewPreferences(a, storage.NewFileURI(filePath), d.db)
	}

	d.prefs = p
	return p
}

func (d *lalStore) dbLocation() string {
	home, _ := os.UserHomeDir()
	dir := d.config
	if dir == "" {
		return filepath.Join(home, "Lal", "fynesync")
	}
	if dir[0] == '~' {
		dir = home + dir[1:]
	}
	return dir
}
