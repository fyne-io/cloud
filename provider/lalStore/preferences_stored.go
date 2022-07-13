// This file is copied from fyne.io/fyne/app/preferences.go to avoid an import loop
package lalStore

import (
	"encoding/json"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/fynelabs/lal"
)

type preferences struct {
	*lalPreferences

	prefLock                   sync.RWMutex
	disconnected, ignoreChange bool

	app  fyne.App
	file fyne.URI
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*preferences)(nil)

func (p *preferences) Disconnect() {
	p.prefLock.Lock()
	p.disconnected = true
	p.prefLock.Unlock()
}

func (p *preferences) resetIgnore() {
	go func() {
		time.Sleep(time.Millisecond * 100) // writes are not always atomic. 10ms worked, 100 is safer.
		p.prefLock.Lock()
		p.ignoreChange = false
		p.prefLock.Unlock()
	}()
}

func (p *preferences) save() error {
	return p.saveToFile(p.file)
}

func (p *preferences) saveToFile(file fyne.URI) error {
	p.prefLock.Lock()
	p.ignoreChange = true
	p.prefLock.Unlock()
	defer p.resetIgnore()
	dir, _ := storage.Parent(file)
	if dir != nil { // may not be supported / needed if the file exists
		_ = storage.CreateListable(dir)
	}

	write, err := storage.Writer(file)
	if err != nil {
		return err
	}
	encode := json.NewEncoder(write)

	p.lalPreferences.ReadValues(func(values map[string]interface{}) {
		err = encode.Encode(&values)
	})

	write.Close()
	return err
}

func (p *preferences) load() {
	err := p.loadFromFile(p.file)
	if err != nil {
		fyne.LogError("Preferences load error:", err)
	}
}

func (p *preferences) loadFromFile(path fyne.URI) (err error) {
	read, err := storage.Reader(path)
	if err != nil {
		return err
	}
	defer func() {
		if r := read.Close(); r != nil && err == nil {
			err = r
		}
	}()
	decode := json.NewDecoder(read)

	p.lalPreferences.WriteValues(func(values map[string]interface{}) {
		err = decode.Decode(&values)
	})

	return err
}

func NewPreferences(a fyne.App, file fyne.URI, db *lal.DB) *preferences {
	p := &preferences{app: a, file: file}
	p.lalPreferences = NewLalPreferences(db)

	// don't load or watch if not setup
	if a.UniqueID() == "" {
		return p
	}
	p.load()

	p.AddChangeListener(func() {
		p.prefLock.RLock()
		shouldIgnoreChange := p.ignoreChange || p.disconnected
		p.prefLock.RUnlock()
		if shouldIgnoreChange { // callback after loading, no need to save
			return
		}

		err := p.save()
		if err != nil {
			fyne.LogError("Failed on saving preferences", err)
		}
	})
	p.watch()
	return p
}
