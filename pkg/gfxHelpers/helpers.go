package gfxHelpers

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

type Fill struct {
	Color color.RGBA
}

func (f Fill) Layout(gtx *layout.Context) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: f.Color}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d, Baseline: d.Y}
}

func RGB(c uint32) color.RGBA {
	return ARGB((0xff << 24) | c)
}

func ARGB(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func ToPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}
