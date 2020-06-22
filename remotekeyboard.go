// Graphical interface to the KDE Connect remote keyboard plugin
package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// Special key codes for KDE Connect remote keyboard plugin
var specialKeys = map[uint]int{
	gdk.KEY_BackSpace: 1,
	gdk.KEY_Tab: 2,
	gdk.KEY_Left: 4,
	gdk.KEY_Up: 5,
	gdk.KEY_Right: 6,
	gdk.KEY_Down: 7,
	gdk.KEY_Page_Up: 8,
	gdk.KEY_Page_Down: 9,
	gdk.KEY_Home: 10,
	gdk.KEY_End: 11,
	gdk.KEY_Return: 12,
	gdk.KEY_Delete: 13,
	gdk.KEY_Escape: 14,
	gdk.KEY_Sys_Req: 15,
	gdk.KEY_Scroll_Lock: 16,
	gdk.KEY_F1: 21,
	gdk.KEY_F2: 22,
	gdk.KEY_F3: 23,
	gdk.KEY_F4: 24,
	gdk.KEY_F5: 25,
	gdk.KEY_F6: 26,
	gdk.KEY_F7: 27,
	gdk.KEY_F8: 28,
	gdk.KEY_F9: 29,
	gdk.KEY_F10: 30,
	gdk.KEY_F11: 31,
	gdk.KEY_F12: 32,
}

// Entry point to the program
func main() {
	// Start dbus connection
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Get KDE Connect device ID
	basePath := getPath(conn)
	if len(basePath) == 0 {
		fmt.Fprintln(os.Stderr, "Could not get device D-Bus path")
		os.Exit(1)
	}

	// Get dbus object for remote keyboard plugin
	bus := conn.Object(
		"org.kde.kdeconnect",
		dbus.ObjectPath(basePath + "/remotekeyboard"),
	)

	// Create GUI
	createWindow(&bus)
}

// Gets a device ID from command line arguments if possible, or gets the first
// available device with the plugin enabled, and returns its D-Bus path
func getPath(conn *dbus.Conn) string {
	if len(os.Args) > 1 {
		// Get ID from command line arguments
		device := os.Args[1]

		pluginEnabled, path := checkPlugin(conn, device)
		if !pluginEnabled {
			// Remote keyboard not enabled on the device
			fmt.Fprintln(os.Stderr, "Remote keyboard plugin not enabled!")
			os.Exit(1)
		}

		return path
	} else {
		// Get device list from D-Bus and find the first one with the plugin
		var devices []string

		// Get devices
		obj := conn.Object("org.kde.kdeconnect", "/modules/kdeconnect")
		err := obj.Call("org.kde.kdeconnect.daemon.devices", 0).Store(&devices)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not get devices:", err)
			os.Exit(1)
		}

		// Check if any device has the plugin enabled
		for _, id := range devices {
			pluginEnabled, path := checkPlugin(conn, id)
			if pluginEnabled {
				return path
			}
		}

		// Remote keyboard not enabled on any device
		fmt.Fprintln(os.Stderr, "Remote keyboard plugin not enabled!")
		os.Exit(1)
	}

	return ""
}

// Checks if the remote keyboard plugin is enabled for a given device
func checkPlugin(conn *dbus.Conn, device string) (bool, string) {
	// Get dbus object for the device
	path := "/modules/kdeconnect/devices/" + device
	obj := conn.Object(
		"org.kde.kdeconnect",
		dbus.ObjectPath(path),
	)

	var enabled bool

	// Check device has the plugin enabled
	err := obj.Call(
		"org.kde.kdeconnect.device.isPluginEnabled",
		0,
		"kdeconnect_remotekeyboard",
	).Store(&enabled)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not check plugin:", err)
		os.Exit(1)
	}

	return enabled, path
}

// Creates and shows the GUI window
func createWindow(bus *dbus.BusObject) {
	gtk.Init(nil)

	// Create window
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create window:", err)
		os.Exit(1)
	}

	win.SetTitle("KDE Connect Remote Keyboard")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Listen to the key press events
	win.Connect("key-press-event", func(_ *gtk.Window, event *gdk.Event) {
		eventKey := gdk.EventKeyNewFromEvent(event)
		keyVal := eventKey.KeyVal()

		// Get input character
		var char string
		if gdk.KeyvalToUnicode(keyVal) != 0 {
			char = string(keyVal)
		}

		special := specialKeys[keyVal]

		// The event state uses bits for each modifier
		state := eventKey.State()
		shift := (state & 1) != 0
		ctrl := (state >> 2 & 1) != 0
		alt := (state >> 3 & 1) != 0

		// Send input to remote device
		call := (*bus).Call(
			"org.kde.kdeconnect.device.remotekeyboard.sendKeyPress",
			0,
			// Method arguments
			char, special, shift, ctrl, alt,
		)
		if call.Err != nil {
			fmt.Fprintln(os.Stderr, "Could not send keypress:", call.Err)
			os.Exit(1)
		}
	})

	// Set size and show window
	win.SetDefaultSize(160, 90)
	win.ShowAll()
	gtk.Main()
}
