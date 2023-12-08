package painter

import rl "github.com/gen2brain/raylib-go/raylib"

type Rect struct {
	Node  rl.Rectangle
	Color rl.Color // Added a Color field
}

// WindowManager manages the window and rectangles
type WindowManager struct {
	rectangles []Rect
}

// NewWindowManager creates a new WindowManager instance
func NewWindowManager() *WindowManager {
	return &WindowManager{}
}

// OpenWindow opens the window
func (wm *WindowManager) OpenWindow(title string, width, height int32) {
	rl.InitWindow(width, height, title)
	rl.SetTargetFPS(60)
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
	}

	rl.EndDrawing()
}

// WindowShouldClose returns true if the window should close
func (wm *WindowManager) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}
