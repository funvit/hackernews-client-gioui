package main

import (
	log "github.com/sirupsen/logrus"
	"hnclient/internal/myapp"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/unit"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug level mode on")
	}

	gofont.Register()

	w := app.NewWindow(
		app.Title("Hacker News reader (minimal)"),
		app.Size(
			unit.Dp(400),
			unit.Dp(800)),
	)

	myapp.New(w)

	app.Main()
}
