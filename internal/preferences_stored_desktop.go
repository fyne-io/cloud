// +build !ios,!android,!mobile,!nacl
// This file is copied from fyne.io/fyne/app/preferences_other.go to avoid an import loop

package internal

import "path/filepath"

// storagePath returns the location of the settings storage
func (p *preferences) storagePath() string {
	return filepath.Join(p.rootDir, p.app.UniqueID(), "preferences.json")
}

func (p *preferences) watch() {
	watchFile(p.storagePath(), func() {
		p.prefLock.RLock()
		shouldIgnoreChange := p.ignoreChange
		p.prefLock.RUnlock()
		if shouldIgnoreChange {
			return
		}

		p.load()
	})
}
