package views

import (
	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"hnclient/pkg/gfxHelpers"
)

type TopBar struct {
	theme material.Theme

	ShowBack bool
	ShowMenu bool

	backClick gesture.Click
	menuClick gesture.Click
}

type (
	BackEvent struct{}
	MenuEvent struct{}
)

func NewTopBar(th material.Theme) *TopBar {
	return &TopBar{
		theme: th,
	}
}

func (t *TopBar) Layout(gtx *layout.Context, insets layout.Inset, w layout.Widget) {
	stack := layout.Stack{Alignment: layout.SW}
	insets = layout.Inset{
		Top:    unit.Add(gtx, insets.Top, unit.Dp(16)),
		Bottom: unit.Dp(16),
		Left:   unit.Max(gtx, insets.Left, unit.Dp(16)),
		Right:  unit.Max(gtx, insets.Right, unit.Dp(16)),
	}
	stackContent := layout.Stacked(func() {
		insets.Layout(gtx, func() {
			flex := layout.Flex{Alignment: layout.Middle}
			//backChild := layout.Rigid(func() {
			//	if t.ShowBack {
			//		ico := (&icon{src: icons.NavigationArrowBack, size: unit.Dp(24)}).image(gtx, gfxHelpers.RGB(0xffffff))
			//		ico.Add(gtx.Ops)
			//		paint.PaintOp{Rect: f32.Rectangle{Max: gfxHelpers.ToPointF(ico.Size())}}.Add(gtx.Ops)
			//		gtx.Dimensions.Size = ico.Size()
			//		gtx.Dimensions.Size.X += gtx.Px(unit.Dp(4))
			//		pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			//		t.backClick.Add(gtx.Ops)
			//	}
			//})
			content := layout.Flexed(1, w)
			//menuChild := layout.Rigid(func() {
			//	if t.ShowMenu {
			//		ico := (&icon{src: icons.NavigationMenu, size: unit.Dp(24)}).image(gtx, gfxHelpers.RGB(0xffffff))
			//		ico.Add(gtx.Ops)
			//		paint.PaintOp{Rect: f32.Rectangle{Max: gfxHelpers.ToPointF(ico.Size())}}.Add(gtx.Ops)
			//		gtx.Dimensions.Size = ico.Size()
			//		gtx.Dimensions.Size.X += gtx.Px(unit.Dp(4))
			//		pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			//		t.menuClick.Add(gtx.Ops)
			//	}
			//})

			//flex.Layout(gtx, backChild, content, menuChild)
			flex.Layout(gtx, content)
		})
	})
	bg := layout.Expanded(func() {
		gfxHelpers.Fill{Color: t.theme.Color.Primary}.Layout(gtx)
	})
	stack.Layout(gtx, bg, stackContent)
}

func (t *TopBar) Event(gtx *layout.Context) interface{} {
	for _, e := range t.backClick.Events(gtx) {
		if e.Type == gesture.TypeClick {
			return BackEvent{}
		}
	}
	for _, e := range t.menuClick.Events(gtx) {
		if e.Type == gesture.TypeClick {
			return MenuEvent{}
		}
	}
	return nil
}
