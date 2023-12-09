package painter

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
	rectangles []Rect
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
func (wm *WindowManager) AddRectangle(rect Rect) {
	wm.rectangles = append(wm.rectangles, rect)
}

// RemoveAllRectangles removes all rectangles from the window
func (wm *WindowManager) RemoveAllRectangles() {
	wm.rectangles = nil
}

// DrawRectangles draws all rectangles on the window
func (wm *WindowManager) DrawRectangles() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	for _, pair := range wm.rectangles {
		rl.DrawRectangleRec(pair.Node, pair.Color)
		if pair.Text.Value != "" {
			// Draw text inside the rectangle
			textHeight := pair.Text.Size
			textX := pair.Node.X
			textY := pair.Node.Y + (pair.Node.Height-float32(textHeight))/2

			font := wm.Fonts[pair.Text.Font]

			if font.Texture.ID == 0 {
				font := rl.LoadFont(pair.Text.Font)
				wm.Fonts[pair.Text.Font] = font
			}

			rl.DrawTextEx(font, pair.Text.Value, rl.NewVector2(textX, textY), pair.Text.Size, 2, pair.Text.Color)
		}
	}

	if wm.FPS {
		wm.FPSCounter.Update()
		wm.FPSCounter.Draw(10, 10, 10, rl.DarkGray)
	}

	rl.EndDrawing()
}

// WindowShouldClose returns true if the window should close
func (wm *WindowManager) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}
