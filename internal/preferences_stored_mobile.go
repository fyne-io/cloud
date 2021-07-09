// +build mobile nacl
// This file is copied from fyne.io/fyne/app/preferences_mobile.go to avoid an import loop

package app

import "path/filepath"

// storagePath returns the location of the settings storage
func (p *preferences) storagePath() string {
	return filepath.Join(p.app.storageRoot(), "preferences.json")
}

// storageRoot returns the location of the app storage
func (a *fyneApp) storageRoot() string {
	return filepath.Join(rootConfigDir(), a.UniqueID())
}

func (p *preferences) watch() {
	// no-op as we are in mobile simulation mode
}