package main

import (
	"fmt"
	"gui/color"
	"gui/cstyle"
	"gui/document"
	painter "gui/painter"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-shiori/dom"
)

func main() {
	css := cstyle.CSS{
		Width:  800,
		Height: 450,
	}
	css.StyleSheet("./master.css")

	d := document.Write("./src/index.html")

	for _, v := range d.StyleSheets {
		css.StyleSheet(v)
	}

	for _, v := range d.StyleTags {
		css.StyleTag(v)
	}

	p := css.Map(d.DOM)

	wm := painter.NewWindowManager()

	// Open the window
	wm.OpenWindow(d.Title, 800, 450)
	defer wm.CloseWindow()

	for _, v := range p.Render {
		styles := p.StyleMap[v.Id]

		x, _ := strconv.ParseFloat(styles["x"], 32)
		y, _ := strconv.ParseFloat(styles["y"], 32)
		width, _ := strconv.ParseFloat(styles["width"], 32)
		height, _ := strconv.ParseFloat(styles["height"], 32)

		fmt.Printf("%s %s %f %f %f %f\n", v.Id, dom.InnerText(v.Node), x, y, width, height)

		bgColor := color.Background(styles)

		node := painter.Rect{
			Node:  rl.NewRectangle(float32(x), float32(y), float32(width), float32(height)),
			Color: rl.NewColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A),
		}

		wm.AddRectangle(node)
	}

	// Main game loop
	for !wm.WindowShouldClose() {
		// Draw rectangles
		wm.DrawRectangles()
	}

}
