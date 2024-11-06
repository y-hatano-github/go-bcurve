package main

import (
	"github.com/nsf/termbox-go"
)

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output256)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	ch := make(chan termbox.Event)
	go keyEvent(ch)

loop:
	for {

		select {
		case ev := <-ch:
			switch ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyEsc {
					break loop
				}
			}
			break
		default:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			drawCurve()
			termbox.Flush()
		}
	}

}

func keyEvent(ch chan termbox.Event) {
	for {
		ch <- termbox.PollEvent()
	}
}

func drawCurve() {
	ps := Points{}
	ps = append(ps, Point{X: 0, Y: 15})
	ps = append(ps, Point{X: 8, Y: 5})
	ps = append(ps, Point{X: 18, Y: 5})
	ps = append(ps, Point{X: 24, Y: 15})
	ps = append(ps, Point{X: 40, Y: 15})
	ps = append(ps, Point{X: 45, Y: 5})

	var cps *Points = &Points{}

	getCurvePoint(ps, 0.5, cps)

	lineP := polyline(*cps)

	for _, p := range lineP {
		termbox.SetCell(p.X, p.Y, ' ', termbox.ColorDefault, termbox.ColorGreen)
	}

	for i, p := range ps {
		if i == 0 || i == len(ps)-1 {
			termbox.SetCell(p.X, p.Y, ' ', termbox.ColorDefault, termbox.ColorLightBlue)
		} else {
			termbox.SetCell(p.X, p.Y, ' ', termbox.ColorDefault, termbox.ColorLightRed)
		}
	}
}

func getCurvePoint(ps Points, t float64, cps *Points) {
	*cps = append(*cps, ps[0])
	if len(ps) > 1 {
		newPs := Points{}
		for i := 0; i < len(ps)-1; i++ {
			x := (1-t)*float64(ps[i].X) + t*float64(ps[i+1].X)
			y := (1-t)*float64(ps[i].Y) + t*float64(ps[i+1].Y)
			//fmt.Printf("[%g, %g]", x, y)
			newPs = append(newPs, Point{X: int(x), Y: int(y)})
		}
		//fmt.Println()
		getCurvePoint(newPs, t, cps)
		*cps = append(*cps, ps[len(ps)-1])
	}
}

type Points []Point

type Point struct {
	X int
	Y int
}

func polyline(ps Points) []Point {
	pl := Points{}

	for i := 0; i < len(ps)-1; i++ {
		pl = append(pl, line(ps[i].X, ps[i].Y, ps[i+1].X, ps[i+1].Y)...)
	}

	return pl
}

func line(x1, y1, x2, y2 int) []Point {

	var dx, dy, sx, sy int

	if x2 > x1 {
		sx = 1
	} else {
		sx = -1
	}
	if x2 > x1 {
		dx = x2 - x1
	} else {
		dx = x1 - x2
	}
	if y2 > y1 {
		sy = 1
	} else {
		sy = -1
	}
	if y2 > y1 {
		dy = y2 - y1
	} else {
		dy = y1 - y2
	}

	ps := []Point{}
	ps = append(ps, Point{X: x1, Y: y1})

	x := x1
	y := y1

	if dx >= dy {
		e := -dx
		for i := 0; i <= dx; i++ {
			ps = append(ps, Point{X: x, Y: y})
			x += sx
			e += 2 * dy
			if e >= 0 {
				y += sy
				e -= 2 * dx
			}
		}

	} else {
		e := -dy
		for i := 0; i <= dy; i++ {
			ps = append(ps, Point{X: x, Y: y})
			y += sy
			e += 2 * dx
			if e >= 0 {
				x += sx
				e -= 2 * dy
			}
		}
	}

	return ps
}
