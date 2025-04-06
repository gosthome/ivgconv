package main

import (
	"context"
	"image/color"
	"log"
	"os"
	"time"

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
	println(fn)
	b, err := ivgconv.FromFile(fn)
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
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	btn := widget.Clickable{}
	icons := os.Args[1:]
	iconIndex := 0
	nextIndex := func() {
		iconIndex++
		if iconIndex >= len(icons) {
			iconIndex = 0
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		done := ctx.Done()
		t := time.NewTicker(20 * time.Millisecond)
		for {
			select {
			case <-done:
				return
			case <-t.C:
				btn.Click()
				window.Invalidate()
			}
		}
	}()
	icon := loadIcon(icons[iconIndex])
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			cancel()
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)
			for btn.Clicked(gtx) {
				nextIndex()
				icon = loadIcon(icons[iconIndex])
			}

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			ib := material.IconButton(theme, &btn, icon, "Hello, Gio")
			s := min(e.Size.X, e.Size.Y)
			ib.Size = unit.Dp(s * 3 / 4)
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
