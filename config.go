package cloud

import (
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

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
		return nil
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

func setCurrentProviderName(s string) {
	store := decodeSettings()
	store.CloudName = s

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
