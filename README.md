# Fyne Cloud

A repository that helps to manage cloud integrations in a Fyne app.
Since Fyne 2.3 developers can provide their own cloud storage to a Fyne app,
but if you would like to use standard providers this repository will manage that for you.

There is also the ability to have developers choose cloud provider and configure settings.
Just add the following lines to your app:

```go
import "fyne.io/cloud"

func main() {
   	a := app.NewWithID("io.fyne.cloud.example")
	cloud.Enable(a)

    ...
}
```

You would then present the cloud choice UI using a button or menu item that simply calls:
```go
    cloud.ShowSettings(a, w)
```
Which will render the settings UI as a dialog on the specified window.