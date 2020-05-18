package widgets

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image/color"
)

type Background struct {
	Color  color.RGBA
	Radius unit.Value
	Inset  layout.Inset
}

func (b *Background) Layout(gtx *layout.Context, w layout.Widget) {
	var macro op.MacroOp
	macro.Record(gtx.Ops)
	b.Inset.Layout(gtx, w)
	macro.Stop()

	var stack op.StackOp
	stack.Push(gtx.Ops)
	size := gtx.Dimensions.Size
	width, height := float32(size.X), float32(size.Y)
	if r := float32(gtx.Px(b.Radius)); r > 0 {
		if r > width/2 {
			r = width / 2
		}
		if r > height/2 {
			r = height / 2
		}
		clip.Rect{
			Rect: f32.Rectangle{Max: f32.Point{
				X: width, Y: height,
			}}, NW: r, NE: r, SW: r, SE: r,
		}.Op(gtx.Ops).Add(gtx.Ops)
	}
	paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: width, Y: height}}}.Add(gtx.Ops)
	macro.Add()
	stack.Pop()
}
