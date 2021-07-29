package cloud

import (
	"fyne.io/cloud/internal/settings"
	"fyne.io/fyne/v2"
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
	store := settings.Load()
	if store == nil {
		return ""
	}
	return store.CloudName
}
