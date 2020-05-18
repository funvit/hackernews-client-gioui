package myapp

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	log "github.com/sirupsen/logrus"
	"hnclient/internal/myapp/views"
)

type MyApp struct {
	w     *app.Window
	theme *material.Theme
	env   Env

	page1 *views.Page1
}

type Env struct {
	Insets layout.Inset
}

func New(w *app.Window) *MyApp {
	a := &MyApp{
		w:     w,
		theme: material.NewTheme(),
		env:   Env{},

		//storiesUpdated: make(chan struct{}, 1),
	}
	a.page1 = views.NewPage1(*a.theme, &a.env, "HackerNews: Top stories")

	go a.loop()

	return a
}

func (e *Env) GetInsets() layout.Inset {
	return e.Insets
}

func (a *MyApp) loop() error {
	gtx := new(layout.Context)
	for {
		select {
		case e := <-a.w.Events():
			switch e := e.(type) {
			//case system.ClipboardEvent:
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx.Reset(e.Queue, e.Config, e.Size)

				a.env.Insets = layout.Inset{
					Top:    e.Insets.Top,
					Left:   e.Insets.Left,
					Right:  e.Insets.Right,
					Bottom: unit.Add(gtx, unit.Dp(24), e.Insets.Bottom),
				}

				// handle click events here

				// layout widgets here
				a.page1.Layout(gtx)

				e.Frame(gtx.Ops)
			}

		case <-a.page1.MustRedraw():
			log.Debug("must redraw event from page1")
			a.w.Invalidate()
		}
	}
}
