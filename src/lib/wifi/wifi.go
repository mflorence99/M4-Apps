package wifi

import (
	_ "embed"
	"encoding/json"

	"tinygo.org/x/drivers/netlink"
	"tinygo.org/x/drivers/netlink/probe"
)

// 🟧 Connect to the Wifi

//go:embed secrets.json
var secrets []byte

type Secrets struct {
	SSID     string `json:"ssid"`
	Password string `json:"pass"`
}

type Wifi struct {
	link netlink.Netlinker
}

// 🟦 Constructor

func NewWifi() *Wifi {
	w := new(Wifi)
	w.link, _ = probe.Probe()
	return w
}

// 🟦 Connect to Wifi

func (w *Wifi) Connect() error {
	// 👇 parse out user credentials
	var s Secrets
	if err := json.Unmarshal(secrets, &s); err != nil {
		return (err)
	}
	params := netlink.ConnectParams{
		Ssid:       s.SSID,
		Passphrase: s.Password,
	}
	// 👇 attempt connection
	if err := w.link.NetConnect(&params); err != nil {
		return (err)
	}
	return nil
}

// 🟦 Disconnect from Wifi

func (w *Wifi) Disconnect() {
	w.link.NetDisconnect()
}

// 🟦 Callback if Wifi goes down

func (w *Wifi) IfDownCall(cb func()) {
	w.link.NetNotify(func(e netlink.Event) {
		switch e {
		case netlink.EventNetUp:
		case netlink.EventNetDown:
			cb()
		}
	})
}
