package main

// 👁️ https://github.com/ardnew/drivers/blob/rgb75/rgb75/rgb75.go

import (
	"image/color"
	"machine"
	"math/rand"
	"rgb75/driver"
	"time"
)

func demo() {

	// panel layout and color depth
	config := driver.Config{
		Width:      64,
		Height:     32,
		ColorDepth: 4,
		DoubleBuf:  true,
	}
	display := &screen{
		// actual rgb75 Device object
		Device: driver.New(
			machine.HUB75_OE, machine.HUB75_LAT, machine.HUB75_CLK,
			[6]machine.Pin{
				machine.HUB75_R1, machine.HUB75_G1, machine.HUB75_B1,
				machine.HUB75_R2, machine.HUB75_G2, machine.HUB75_B2,
			},
			[]machine.Pin{
				machine.HUB75_ADDR_A, machine.HUB75_ADDR_B, machine.HUB75_ADDR_C,
				machine.HUB75_ADDR_D, machine.HUB75_ADDR_E,
			}),
		// private data structure representing a multi-colored, continuous
		// stream of moving pixels
		trail: []*trail{
			newTrail(config.Width, 1.0),
			newTrail(config.Width, 2.0),
			newTrail(config.Width, 2.0),
			newTrail(config.Width, 0.5),
			newTrail(config.Width, 1.5),
			newTrail(config.Width, 1.5),
			newTrail(config.Width, 0.7),
			newTrail(config.Width, 1.3),
		},
	}

	if err := display.Configure(config); nil != err {
		for {
			println("error: " + err.Error())
			time.Sleep(time.Second)
		}
	}
	display.Resume()

	for {
		for _, tr := range display.trail {
			// update trail head
			if !display.contains(tr.inc()) {
				tr.wrap()
			}
			// draw each valid pixel of the trail
			for _, px := range tr.pix {
				if display.contains(px.point) {
					x, y := px.point.pos()
					display.SetPixel(x, y, px.color)
				}
			}
		}
		if err := display.Display(); nil != err {
			println("error: " + err.Error())
		}

		time.Sleep(10 * time.Millisecond)
	}
}

// screen associates an embedded rgb75.Device with a list of trails to animate.
type screen struct {
	*driver.Device
	trail []*trail
}

// contains returns true if and only if the given point exists within the
// receiver's screen dimensions.
func (s *screen) contains(p point) bool {
	width, height := s.Size()
	aboveMin := p.x >= 0 && p.y >= 0
	belowMax := int16(p.x+0.5) < width && int16(p.y+0.5) < height
	return aboveMin && belowMax
}

// pixel represents an individual RGB color with given point coordinates.
type pixel struct {
	point point
	color color.RGBA
}

// point represents a 2-dimensional point in space.
type point struct{ x, y float32 }

// noPoint represents a point that exists outside of any screen's coordinate
// space.
var noPoint = point{x: -256.0, y: -256.0}

// pos returns the x and y components of the receiver, rounded to the nearest
// int16 integer. Note that x and y may be negative.
func (p point) pos() (x, y int16) {
	round := func(f float32) int16 {
		// naive rounding (half-away)
		if f < 0 {
			f -= 0.5
		} else {
			f += 0.5
		}
		return int16(f)
	}
	return round(p.x), round(p.y)
}

// next returns a new point whose x and y components are equal to the receiver's
// x and y components incremented by random deltas.
// The x component delta is a random value in the interval (-1,1).
// The y component delta is a random value in the interval [0,1), multiplied by
// the given factor speed.
func (p point) next(speed float32) point {
	dx, dy := rand.Float32(), rand.Float32()*speed
	if rand.Intn(2) == 0 {
		return point{x: p.x - dx, y: p.y + dy}
	}
	return point{x: p.x + dx, y: p.y + dy}
}

// top returns a new point positioned at the top of a screen.
// The x component is a random integer in the interval [xMin,xMax).
// The y component is always 0.
func top(xMin, xMax int) point {
	if xMin > xMax {
		xMin, xMax = xMax, xMin
	}
	return point{
		x: float32(int32(xMin) + rand.Int31n(int32(xMax-xMin))),
		y: 0.0,
	}
}

// trail contains a queue of pixels, a coordinate for the head of the queue to
// move into, and a width dimension defining the horizontal range of head.
type trail struct {
	pix []pixel
	pos point
	dim int
	fac float32
}

func newTrail(xSpan int, ySpeed float32) *trail {
	return &trail{
		pix: []pixel{
			{point: noPoint, color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0x1, G: 0x0, B: 0x0, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0x3, G: 0x0, B: 0x0, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0x7, G: 0x0, B: 0x0, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0xF, G: 0x0, B: 0x0, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0xF, G: 0x0, B: 0x0, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0xF, G: 0x0, B: 0x1, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0x7, G: 0x0, B: 0x3, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0x3, G: 0x0, B: 0x7, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0x1, G: 0x0, B: 0xF, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0x0, G: 0xF, B: 0x0, A: 0xF}},
			{point: noPoint, color: color.RGBA{R: 0x0, G: 0xF, B: 0x0, A: 0xF}},
		},
		pos: top(0, xSpan),
		dim: xSpan,
		fac: ySpeed,
	}
}

// push enqueues each of the given points to the receiver trail in the order
// they were provided.
func (t *trail) push(ps ...point) {
	for _, p := range ps {
		for i, px := range t.pix[1:] {
			t.pix[i].point = px.point
		}
		t.pix[len(t.pix)-1].point = p
	}
}

// inc pushes the trail head onto the pixel queue and then increments head,
// returning its new coordinates.
func (t *trail) inc() point {
	// add the current position to the list
	t.push(t.pos)
	// update the current position
	t.pos = t.pos.next(t.fac)
	return t.pos
}

// wrap resets the trail head to the top of the screen.
func (t *trail) wrap() {
	t.pos = top(0, t.dim)
}
