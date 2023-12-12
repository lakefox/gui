package window

import (
	"gui/fps"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rect struct {
	Node  rl.Rectangle
	Color rl.Color // Added a Color field
	Text  Text
}

type Text struct {
	Color rl.Color
	Size  float32
	Value string
	Font  string
}

// WindowManager manages the window and rectangles
type WindowManager struct {
	nodes      []Rect
	Fonts      map[string]rl.Font
	FPS        bool
	FPSCounter fps.FPSCounter
}

// NewWindowManager creates a new WindowManager instance
func NewWindowManager() *WindowManager {
	fpsCounter := fps.NewFPSCounter()

	return &WindowManager{
		Fonts:      make(map[string]rl.Font),
		FPSCounter: *fpsCounter,
	}
}

// OpenWindow opens the window
func (wm *WindowManager) OpenWindow(title string, width, height int32) {
	rl.InitWindow(width, height, title)
	rl.SetTargetFPS(60)
	// Enable window resizing
	rl.SetWindowState(rl.FlagWindowResizable)
}

// CloseWindow closes the window
func (wm *WindowManager) CloseWindow() {
	rl.CloseWindow()
}

// AddRectangle adds a rectangle to the window
func (wm *WindowManager) AddNode(rect Rect) {
	wm.nodes = append(wm.nodes, rect)
}

// RemoveAllNods removes all nodes from the window
func (wm *WindowManager) RemoveAllNodes() {
	wm.nodes = nil
}

// Draw draws all nodes on the window
func (wm *WindowManager) Draw() {

	for _, node := range wm.nodes {
		rl.DrawRectangleRec(node.Node, node.Color)
		if node.Text.Value != "" {
			// Draw text inside the rectangle
			textHeight := node.Text.Size
			textX := node.Node.X
			textY := node.Node.Y + (node.Node.Height-float32(textHeight))/2

			font := wm.Fonts[node.Text.Font]

			if font.Texture.ID == 0 {
				font := rl.LoadFont(node.Text.Font)
				wm.Fonts[node.Text.Font] = font
			}

			rl.DrawTextEx(font, node.Text.Value, rl.NewVector2(textX, textY), node.Text.Size, 2, node.Text.Color)
		}
	}

	if wm.FPS {
		wm.FPSCounter.Update()
		wm.FPSCounter.Draw(10, 10, 10, rl.DarkGray)
	}
}

// WindowShouldClose returns true if the window should close
func (wm *WindowManager) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}
