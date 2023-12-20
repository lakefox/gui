package window

import (
	"gui/cstyle"
	"gui/fps"
	"image/color"

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
	Textures   []rl.Texture2D
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
	rl.SetTargetFPS(30)
	// Enable window resizing
	rl.SetWindowState(rl.FlagWindowResizable)
}

// CloseWindow closes the window
func (wm *WindowManager) CloseWindow() {
	rl.CloseWindow()
}

func (wm *WindowManager) LoadTextures(nodes []cstyle.Node) {
	wm.Textures = make([]rl.Texture2D, len(nodes))
	for i, node := range nodes {
		if node.Text.Image != nil {
			texture := rl.LoadTextureFromImage(rl.NewImageFromImage(node.Text.Image))
			wm.Textures[i] = texture
		}
	}
}

// Draw draws all nodes on the window
func (wm *WindowManager) Draw(nodes []cstyle.Node) {

	for i, node := range nodes {
		rl.DrawRectangle(int32(node.X), int32(node.Y), int32((node.Width)), int32(node.Height), node.Colors.Background)
		if node.Text.Image != nil {
			r, g, b, a := node.Text.Color.RGBA()
			rl.DrawTexture(wm.Textures[i], int32(node.X+node.Padding.Left), int32(node.Y+node.Padding.Top), color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
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
