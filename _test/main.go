package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gio-eui/ivgconv"
)

func loadIcon(fn string) *widget.Icon {
	b, err := ivgconv.FromFile(fn, ivgconv.WithOutputSize(48))
	if err != nil {
		panic(err)
	}
	icon, err := widget.NewIcon(b)
	if err != nil {
		panic(err)
	}
	return icon
}

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.MinSize(unit.Dp(128), unit.Dp(128)))
		err := run(window, loadIcon(os.Args[1]))
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window, icon *widget.Icon) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	btn := widget.Clickable{}
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			ib := material.IconButton(theme, &btn, icon, "Hello, Gio")
			ib.Size = unit.Dp(128)
			ib.Color = maroon

			// Change the position of the label.
			// title.Alignment = text.Middle

			// Draw the label to the graphics context.
			ib.Layout(gtx)
			// icon.Layout(gtx, maroon)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
