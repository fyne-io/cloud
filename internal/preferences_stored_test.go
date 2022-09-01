// This file is copied from fyne.io/fyne/app/preferences_test.go to avoid an import loop
package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func loadPreferences() *preferences {
	uri := storage.NewFileURI("dummy.json")
	p := NewPreferences(test.NewApp(), uri)
	p.load()

	return p
}

func TestPreferences_Save(t *testing.T) {
	p := loadPreferences()
	p.WriteValues(func(val map[string]interface{}) {
		val["keyString"] = "value"
		val["keyInt"] = 4
		val["keyFloat"] = 3.5
		val["keyBool"] = true
	})

	path := filepath.Join(os.TempDir(), "fynePrefs.json")
	defer os.Remove(path)
	p.saveToFile(storage.NewFileURI(path))

	expected, err := ioutil.ReadFile(filepath.Join("dummy.json"))
	if err != nil {
		assert.Fail(t, "Failed to load, %v", err)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		assert.Fail(t, "Failed to load, %v", err)
	}
	assert.JSONEq(t, string(expected), string(content))
}

func TestPreferences_Load(t *testing.T) {
	p := loadPreferences()
	path, _ := filepath.Abs(filepath.Join("testdata", "dummy.json"))
	p.loadFromFile(storage.NewFileURI(path))

	assert.Equal(t, "value", p.String("keyString"))
	assert.Equal(t, 4, p.Int("keyInt"))
	assert.Equal(t, 3.5, p.Float("keyFloat"))
	assert.Equal(t, true, p.Bool("keyBool"))
}
