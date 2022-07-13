package lalStore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/fynelabs/lal"
)

type lalPreferences struct {
	values          map[string]interface{}
	lock            sync.RWMutex
	changeListeners []func()
	wg              *sync.WaitGroup

	db *lal.DB
}

// enforce conformity
var _ (fyne.Preferences) = (*lalPreferences)(nil)
var _ (fyne.CloudProvider) = (*lalStore)(nil)

// AddChangeListener allows code to be notified when some preferences change. This will fire on any update.
// The passed 'listener' should not try to write values.
func (p *lalPreferences) AddChangeListener(listener func()) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.changeListeners = append(p.changeListeners, listener)
}

// ChangeListeners
func (p *lalPreferences) ChangeListeners() []func() {
	return p.changeListeners
}

// ReadValues provides read access to the underlying value map - for internal use only...
// You should not retain a reference to the map nor write to the values in the callback function
func (p *lalPreferences) ReadValues(fn func(map[string]interface{})) {
	p.lock.RLock()
	fn(p.values)
	p.lock.RUnlock()
}

// WriteValues provides write access to the underlying value map - for internal use only...
// You should not retain a reference to the map passed to the callback function
func (p *lalPreferences) WriteValues(fn func(map[string]interface{})) {
	p.lock.Lock()
	fn(p.values)
	p.lock.Unlock()

	p.fireChange()
}

func (p *lalPreferences) set(key string, value interface{}) {
	p.lock.Lock()

	// fmt.Println("setting")
	if p.db == nil {
		fmt.Println("nil db")
	}

	//Update(CreateIfNotExistBucket+Gob.Encode+Put) for each Settter and View(Bucket+Get+Gob.Decode)
	if err := p.db.Update(func(tx *lal.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("widgets"))
		if err != nil {
			fyne.LogError("Can't create bucket", err)
		}
		var keyBuffer bytes.Buffer
		keyEnc := gob.NewEncoder(&keyBuffer)
		err = keyEnc.Encode(key)
		if err != nil {
			fyne.LogError("Can't encode key to bytes buffer", err)
		}

		var valueBuffer bytes.Buffer
		valEnc := gob.NewEncoder(&valueBuffer)
		err = valEnc.Encode(value)
		if err != nil {
			fyne.LogError("Can't encode value to bytes buffer", err)
		}

		if err := b.Put(keyBuffer.Bytes(), valueBuffer.Bytes()); err != nil {
			fyne.LogError("Can't Put", err)
		}
		return nil
	}); err != nil {
		fyne.LogError("Can't update bucket", err)
	}

	p.values[key] = value
	p.lock.Unlock()

	p.fireChange()
}

func (p *lalPreferences) get(key string) (interface{}, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	v, err := p.values[key]
	return v, err
}

func (p *lalPreferences) remove(key string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.values, key)
}

func (p *lalPreferences) fireChange() {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for _, l := range p.changeListeners {
		p.wg.Add(1)
		go func(listener func()) {
			defer p.wg.Done()
			listener()
		}(l)
	}

	p.wg.Wait()
}

// Bool looks up a boolean value for the key
func (p *lalPreferences) Bool(key string) bool {
	return p.BoolWithFallback(key, false)
}

// BoolWithFallback looks up a boolean value and returns the given fallback if not found
func (p *lalPreferences) BoolWithFallback(key string, fallback bool) bool {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}

	val, ok := value.(bool)
	if !ok {
		return false
	}
	return val
}

// SetBool saves a boolean value for the given key
func (p *lalPreferences) SetBool(key string, value bool) {
	p.set(key, value)
}

// Float looks up a float64 value for the key
func (p *lalPreferences) Float(key string) float64 {
	return p.FloatWithFallback(key, 0.0)
}

// FloatWithFallback looks up a float64 value and returns the given fallback if not found
func (p *lalPreferences) FloatWithFallback(key string, fallback float64) float64 {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}

	val, ok := value.(float64)
	if !ok {
		return 0.0
	}
	return val
}

// SetFloat saves a float64 value for the given key
func (p *lalPreferences) SetFloat(key string, value float64) {
	p.set(key, value)
}

// Int looks up an integer value for the key
func (p *lalPreferences) Int(key string) int {
	return p.IntWithFallback(key, 0)
}

// IntWithFallback looks up an integer value and returns the given fallback if not found
func (p *lalPreferences) IntWithFallback(key string, fallback int) int {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}

	// integers can be de-serialised as floats, so support both
	if intVal, ok := value.(int); ok {
		return intVal
	}
	val, ok := value.(float64)
	if !ok {
		return 0
	}
	return int(val)
}

// SetInt saves an integer value for the given key
func (p *lalPreferences) SetInt(key string, value int) {
	p.set(key, value)
}

// String looks up a string value for the key
func (p *lalPreferences) String(key string) string {
	return p.StringWithFallback(key, "")
}

// StringWithFallback looks up a string value and returns the given fallback if not found
func (p *lalPreferences) StringWithFallback(key, fallback string) string {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}

	return value.(string)
}

// SetString saves a string value for the given key
func (p *lalPreferences) SetString(key string, value string) {
	p.set(key, value)
}

// RemoveValue deletes a value on the given key
func (p *lalPreferences) RemoveValue(key string) {
	p.remove(key)
}

// NewLalPreferences creates a new preferences implementation
func NewLalPreferences(db *lal.DB) *lalPreferences {
	return &lalPreferences{
		values: make(map[string]interface{}),
		wg:     &sync.WaitGroup{},
		db:     db,
	}
}
