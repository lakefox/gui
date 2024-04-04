package document

import (
	"gui/cstyle"
	"gui/element"
	"gui/events"
	"gui/window"

	"gui/cstyle/plugins/block"
	"gui/cstyle/plugins/flex"
	"gui/cstyle/plugins/inline"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/net/html"
)

type Window struct {
	StyleSheets []string
	StyleTags   []string
	DOM         *html.Node
	Title       string
}

type Document struct {
	CSS cstyle.CSS
}

func (doc Document) Open(index string, script func(*element.Node)) {
	d := parse(index)

	wm := window.NewWindowManager()
	wm.FPS = true

	// Initialization
	var screenWidth int32 = 800
	var screenHeight int32 = 450

	// Open the window
	wm.OpenWindow(screenWidth, screenHeight)
	defer wm.CloseWindow()

	doc.CSS = cstyle.CSS{
		Width:  800,
		Height: 450,
	}
	doc.CSS.StyleSheet("./master.css")
	// css.AddPlugin(position.Init())
	doc.CSS.AddPlugin(inline.Init())
	doc.CSS.AddPlugin(block.Init())
	doc.CSS.AddPlugin(flex.Init())

	for _, v := range d.StyleSheets {
		doc.CSS.StyleSheet(v)
	}

	for _, v := range d.StyleTags {
		doc.CSS.StyleTag(v)
	}

	nodes := doc.CSS.CreateDocument(d.DOM)
	root := &nodes

	script(root)

	// fmt.Println(nodes.Style)

	evts := map[string]element.EventList{}

	eventStore := &evts

	// Main game loop
	for !wm.WindowShouldClose() {
		rl.BeginDrawing()

		// Check if the window size has changed
		newWidth := int32(rl.GetScreenWidth())
		newHeight := int32(rl.GetScreenHeight())

		if newWidth != screenWidth || newHeight != screenHeight {
			rl.ClearBackground(rl.RayWhite)
			// Window has been resized, handle the event
			screenWidth = newWidth
			screenHeight = newHeight

			doc.CSS.Width = float32(screenWidth)
			doc.CSS.Height = float32(screenHeight)

			nodes = doc.CSS.CreateDocument(d.DOM)
			root = &nodes
			script(root)
		}

		eventStore = events.GetEvents(root, eventStore)
		doc.CSS.ComputeNodeStyle(root)
		rd := doc.CSS.Render(*root)
		wm.LoadTextures(rd)
		wm.Draw(rd)

		events.RunEvents(eventStore)

		rl.EndDrawing()
	}
}
