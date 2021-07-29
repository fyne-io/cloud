package cloud

import (
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// Configurable interface describes the functions required for a cloud provider to store configuration.
type Configurable interface {
	// Configure requests that the cloud provider show some configuration options as a dialog on the current window.
	// It returns a serialised configuration or an error.
	Configure(fyne.Window) (string, error)
	// SetConfig is used to apply a previous configuration to this provider.
	SetConfig(string)
}

func currentProviderName() string {
	store := decodeSettings()
	if store == nil {
		return ""
	}
	return store.CloudName
}

func decodeSettings() *app.SettingsSchema {
	store := &app.SettingsSchema{}

	file, err := os.Open(store.StoragePath()) // #nosec
	if err != nil {
		if !os.IsNotExist(err) {
			fyne.LogError("Unable to read global settings", err)
		}
		return store
	}
	decode := json.NewDecoder(file)

	err = decode.Decode(&store)
	if err != nil {
		fyne.LogError("Failed to decode global settings", err)
		return nil
	}
	_ = file.Close()

	return store
}

func saveSettings(store *app.SettingsSchema) {
	file, err := os.Create(store.StoragePath()) // #nosec
	if err != nil {
		fyne.LogError("Failed to create global settings file", err)
		return
	}

	encode := json.NewEncoder(file)
	err = encode.Encode(&store)
	if err != nil {
		fyne.LogError("Failed to encode global settings", err)
	}
}

func setCurrentProviderName(s string) {
	store := decodeSettings()
	store.CloudName = s

	saveSettings(store)
}

func setProviderConfig(s string) {
	store := decodeSettings()
	store.CloudConfig = s

	saveSettings(store)
}
