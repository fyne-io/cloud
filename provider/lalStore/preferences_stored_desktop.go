//go:build !ios && !android && !mobile && !nacl
// +build !ios,!android,!mobile,!nacl

// This file is copied from fyne.io/fyne/app/preferences_other.go to avoid an import loop

package lalStore

func (p *preferences) watch() {
	watchFile(p.file.Path(), func() {
		p.prefLock.RLock()
		shouldIgnoreChange := p.ignoreChange
		p.prefLock.RUnlock()
		if shouldIgnoreChange {
			return
		}

		p.load()
	})
}
