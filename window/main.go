package window

import (
	"gui/element"
	"gui/fps"
	"gui/utils"
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
func (wm *WindowManager) OpenWindow(width, height int32) {
	rl.InitWindow(width, height, "")
	rl.SetTargetFPS(30)
	// Enable window resizing
	rl.SetWindowState(rl.FlagWindowResizable)
}

// CloseWindow closes the window
func (wm *WindowManager) CloseWindow() {
	rl.CloseWindow()
}

func (wm *WindowManager) LoadTextures(nodes []element.Node) {
	wm.Textures = make([]rl.Texture2D, len(nodes))
	for i, node := range nodes {
		if node.Properties.Text.Image != nil {
			texture := rl.LoadTextureFromImage(rl.NewImageFromImage(node.Properties.Text.Image))
			wm.Textures[i] = texture
		}
	}
}

// Draw draws all nodes on the window
func (wm *WindowManager) Draw(nodes []element.Node) {

	for i, node := range nodes {
		bw, _ := utils.ConvertToPixels(node.Properties.Border.Width, node.Properties.EM, node.Properties.Width)
		rad, _ := utils.ConvertToPixels(node.Properties.Border.Radius, node.Properties.EM, node.Properties.Width)

		rect := rl.NewRectangle(node.Properties.X+bw,
			node.Properties.Y+bw,
			node.Properties.Width-(bw+bw),
			(node.Properties.Height+(node.Properties.Padding.Top+node.Properties.Padding.Bottom))-(bw+bw),
		)

		rl.DrawRectangleRoundedLines(rect, rad/200, 1000, bw, node.Properties.Border.Color)
		rl.DrawRectangleRounded(rect, rad/200, 1000, node.Properties.Colors.Background)

		if node.Properties.Text.Image != nil {
			r, g, b, a := node.Properties.Text.Color.RGBA()
			rl.DrawTexture(wm.Textures[i], int32(node.Properties.X+node.Properties.Padding.Left+bw), int32(node.Properties.Y+node.Properties.Padding.Top), color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	if wm.FPS {
		wm.FPSCounter.Update()
		wm.FPSCounter.Draw(10, 10, 10, rl.DarkGray)
	}
	// touching := events.GetEvents(nodes)

	// if touching.Id != "" {
	// 	fmt.Println(touching.Id)
	// }

}

// WindowShouldClose returns true if the window should close
func (wm *WindowManager) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}
