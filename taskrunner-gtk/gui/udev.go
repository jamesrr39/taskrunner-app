package gui

/*
import (
	udev "github.com/jochenvg/go-udev"
	"github.com/mattn/go-gtk/gtk"
)

type UdevDevicesScene struct {
}

func (u *UdevDevicesScene) IsCurrentlyRendered() bool {
	return false //stub
}

func (u *UdevDevicesScene) Content() gtk.IWidget {
	// Create Udev and Enumerate
	u := udev.Udev{}
	e := u.NewEnumerate()

	// Add some FilterAddMatchSubsystemDevtype
	e.AddMatchSubsystem("block")
	e.AddMatchIsInitialized()
	devices, _ := e.Devices()
	for i := range devices {
		device := devices[i]
		fmt.Println(device.Syspath())
	}

	return nil
}
*/
