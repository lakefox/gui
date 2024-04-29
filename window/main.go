package window

import (
	"gui/element"
	"gui/fps"
	"gui/utils"
	"hash/fnv"
	"image"
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
	Textures   map[int]Texture
}

type Texture struct {
	Hash  uint64
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

func (wm *WindowManager) LoadTextures(nodes []element.State) {
	if wm.Textures == nil {
		wm.Textures = map[int]Texture{}
	}
	for i, node := range nodes {
		if node.Texture != nil {
			hash := computeImageHash(node.Texture)
			if wm.Textures[i].Hash != hash {
				rl.UnloadTexture(wm.Textures[i].Image)
				texture := rl.LoadTextureFromImage(rl.NewImageFromImage(node.Texture))
				wm.Textures[i] = Texture{
					Hash:  hash,
					Image: texture,
				}
			}

		}

	}
}

// Draw draws all nodes on the window
func (wm *WindowManager) Draw(nodes []element.State) {

	for i, node := range nodes {
		bw, _ := utils.ConvertToPixels(node.Border.Width, node.EM, node.Width)
		rad, _ := utils.ConvertToPixels(node.Border.Radius, node.EM, node.Width)

		p := node.Padding

		rect := rl.NewRectangle(node.X+bw,
			node.Y+bw,
			node.Width-(bw+bw),
			(node.Height + (p.Top + p.Bottom)),
		)

		// node.Background.A = 100
		// node.Background.R = uint8((255 / len(nodes)) * i)
		// node.Background.G = uint8((255 / len(nodes)) * i)
		// node.Background.B = uint8((255 / len(nodes)) * i)

		rl.DrawRectangleRoundedLines(rect, rad/200, 1000, bw, node.Border.Color)
		rl.DrawRectangleRounded(rect, rad/200, 1000, node.Background)

		// fmt.Println(node.Text.Image == nil, node.Text.Text)
		// fmt.Printf("%v\n", node.Text)

		if node.Texture != nil {
			r, g, b, a := node.Color.RGBA()
			rl.DrawTexture(wm.Textures[i].Image, int32(node.X+p.Left+bw), int32(node.Y+p.Top+bw), ic.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
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

func computeImageHash(img *image.RGBA) uint64 {
	var hash uint64
	hasher := fnv.New64a()

	// Combine the pixel values to generate the hash
	for _, pixel := range img.Pix {
		hasher.Write([]byte{pixel})
	}
	hash = hasher.Sum64()

	return hash
}
