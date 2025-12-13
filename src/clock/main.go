package main

import (
	_ "embed"
	"time"

	"m4-apps/lib/ntp"
	"m4-apps/lib/rgb75"
	"m4-apps/lib/utils"
	"m4-apps/lib/wifi"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

// 🟧 Display date/time synchronized by NTP

func main() {

	var err error
	var font = &proggy.TinySZ8pt7b
	var lastSyncd time.Time
	var syncd bool

	// 👇 noop in production
	utils.WaitForSerial()

	// 👇 connect to Wifi
	w := wifi.NewWifi()
	err = w.Connect()
	if err != nil {
		println("🔥 Wifi connection failed", err.Error())
		panic(err)
	}

	// 👇 ...and disconnect when we're done
	defer w.Disconnect()

	// 👇 prepare matrix display
	d := rgb75.NewDevice()

	// 👇 need to continually refresh display
	for {

		// 👇 resync clock with NTP every hour
		if !syncd || time.Since(lastSyncd).Minutes() > 10 {
			err := ntp.SyncSystemTime()
			if err != nil {
				println("🔥 SyncSystemTime failed", err.Error())
			} else {
				println("🐞 system time synchronized", time.Now().String())
			}
			// 👇 setup
			syncd = true
			lastSyncd = time.Now()
		}

		// TODO 🔥 HACK: PST only as TZ and ignoring DST
		pst := time.FixedZone("PST", -8*60*60)
		t := time.Now().In(pst)

		// 👇 finally!
		tinyfont.WriteLineColors(d, font, 15, 12, t.Format("Jan 2"), utils.Colors)
		tinyfont.WriteLineColors(d, font, 2, 26, t.Format("3:04:05pm"), utils.Colors)
		d.Display()

		// 👇 take a beat to minimize any flicker
		time.Sleep(1 * time.Millisecond)
	}

}
