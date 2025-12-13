package main

import (
	"image/color"
	"m4-apps/lib/rgb75"
	"m4-apps/lib/utils"
	"math/rand"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
	"tinygo.org/x/tinyfont/notoemoji"
)

var emoji = notoemoji.NotoEmojiRegular12pt
var font = freesans.BoldOblique12pt7b
var msg = "Ho ho ho! ⭐⭐⭐ Merry Christmas! ❤ ❤ ❤"
var msgColors = make([]color.RGBA, len(msg))
var _, lineWidth = tinyfont.LineWidth(&font, msg)

func main() {

	firstTime := true
	lastColorTick := time.Now()
	lastShiftTick := time.Now()

	// 👇 initially off-screen
	x := -int16(lineWidth + 1)

	d := rgb75.NewDevice()

	for {

		// 👇 if totally shifted left off-screen, start over
		if x+int16(lineWidth) < 0 {
			x = 64
		}

		// 👇 change colors at random every N seconds
		if firstTime || time.Since(lastColorTick).Seconds() > 2 {
			for ix := 0; ix < len(msgColors); ix++ {
				color := utils.Colors[rand.Intn(len(utils.Colors))]
				msgColors[ix] = color
			}
			lastColorTick = time.Now()
		}

		// 👇 shift left 1 pixel every M milliseconds
		//    note that the M4 can't really keep up, in reality the
		//    shift takes a bit longer
		if firstTime || time.Since(lastShiftTick).Milliseconds() > 50 {
			x -= 1
			lastShiftTick = time.Now()
		}

		firstTime = false

		// 👇 draw the message starting at the latest X position
		//    we cheat a little looking for emojis and use the special
		//    font and a fixed color
		xx := x
		for ix, char := range msg {
			f := utils.Ternary(char > 256, &emoji, &font)
			c := utils.Ternary(char > 256, color.RGBA{255, 0, 0, 255}, msgColors[ix])
			tinyfont.DrawChar(d, f, xx, 24, char, c)
			_, cx := tinyfont.LineWidth(&font, string(char))
			xx += int16(cx)
		}

		// 👇 the M4 can barely keep up with this, so just display as fast
		//    as we can -- we could optimize by not drawing the entire message,
		//    only those letters that are currently on screen
		d.Display()
	}

}
