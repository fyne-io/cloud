package settings

import (
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func Load() *app.SettingsSchema {
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

func SetCurrentProviderName(s string) {
	store := Load()
	store.CloudName = s

	saveSettings(store)
}

func SetProviderConfig(s string) {
	store := Load()
	store.CloudConfig = s

	saveSettings(store)
}
