package window

import (
	"gui/color"
	"gui/element"
	"gui/fps"
	"gui/utils"
	ic "image/color"

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
	Textures   map[string]TextTexture
}

type TextTexture struct {
	Text  string
	Image rl.Texture2D
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
	if wm.Textures == nil {
		wm.Textures = map[string]TextTexture{}
	}
	for _, node := range nodes {
		if node.Properties.Text.Image != nil {
			if wm.Textures[node.Properties.Id].Text != node.InnerText {
				rl.UnloadTexture(wm.Textures[node.Properties.Id].Image)
				texture := rl.LoadTextureFromImage(rl.NewImageFromImage(node.Properties.Text.Image))
				wm.Textures[node.Properties.Id] = TextTexture{
					Text:  node.InnerText,
					Image: texture,
				}
			}

		}

	}
}

// Draw draws all nodes on the window
func (wm *WindowManager) Draw(nodes []element.Node) {

	for _, node := range nodes {
		bw, _ := utils.ConvertToPixels(node.Properties.Border.Width, node.Properties.EM, node.Properties.Computed["width"])
		rad, _ := utils.ConvertToPixels(node.Properties.Border.Radius, node.Properties.EM, node.Properties.Computed["width"])

		p := utils.GetMP(node, "padding")

		rect := rl.NewRectangle(node.Properties.X+bw,
			node.Properties.Y+bw,
			node.Properties.Computed["width"]-(bw+bw),
			(node.Properties.Computed["height"]+(p.Top+p.Bottom))-(bw+bw),
		)

		rl.DrawRectangleRoundedLines(rect, rad/200, 1000, bw, node.Properties.Border.Color)
		rl.DrawRectangleRounded(rect, rad/200, 1000, color.Parse(node.Style, "background"))

		if node.Properties.Text.Image != nil {
			r, g, b, a := node.Properties.Text.Color.RGBA()
			rl.DrawTexture(wm.Textures[node.Properties.Id].Image, int32(node.Properties.X+p.Left+bw), int32(node.Properties.Y+p.Top), ic.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
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
