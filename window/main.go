package window

import (
	"gui/element"
	"gui/fps"
	"gui/utils"
	"hash/fnv"
	"image"
	ic "image/color"
	"slices"
	"sort"

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
	Fonts        map[string]rl.Font
	FPSCounterOn bool
	FPS          int32
	FPSCounter   fps.FPSCounter
	Textures     map[int]Texture
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

func (wm *WindowManager) SetFPS(fps int32) {
	wm.FPS = fps
	rl.SetTargetFPS(fps)
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
			// !TODO: Make a faster hash algo that minimises the time to detect if a image is different
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
	indexes := []float32{0}
	// !TODO: Only Draw whats in fov

	for a := 0; a < len(indexes); a++ {
		for i, node := range nodes {

			if node.Z == indexes[a] {
				rad := utils.ConvertToPixels(node.Border.Radius, node.EM, node.Width)
				rad = rad * 3
				p := node.Padding

				rect := rl.NewRectangle(node.X,
					node.Y,
					node.Width+node.Border.Left.Width+node.Border.Right.Width,
					node.Height+node.Border.Top.Width+node.Border.Bottom.Width,
				)

				rl.DrawRectangleRounded(rect, rad/200, 1000, node.Background)

				// Draw the border based on the style for each side
				img := image.NewRGBA(image.Rect(int(node.X),
					int(node.Y), int(node.X+
						node.Width+node.Border.Left.Width+node.Border.Right.Width),
					int(node.Y+node.Height+node.Border.Top.Width+node.Border.Bottom.Width)))
				DrawBorder(img, img.Bounds(), node.Border, 0)

				rl.DrawTexture(rl.LoadTextureFromImage(rl.NewImageFromImage(img)), int32(node.X), int32(node.Y), rl.White)

				if node.Texture != nil {
					r, g, b, a := node.Color.RGBA()
					rl.DrawTexture(wm.Textures[i].Image, int32(node.X+p.Left+node.Border.Left.Width), int32(node.Y+p.Top+node.Border.Top.Width), ic.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
				}
			} else {
				if !slices.Contains(indexes, node.Z) {
					indexes = append(indexes, node.Z)
					sort.Slice(indexes, func(i, j int) bool {
						return indexes[i] < indexes[j]
					})
				}
			}
		}
	}

	if wm.FPSCounterOn {
		wm.FPSCounter.Update()
		wm.FPSCounter.Draw(10, 10, 10, rl.DarkGray)
	}

}

// WindowShouldClose returns true if the window should close
func (wm *WindowManager) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}

func computeImageHash(img *image.RGBA) uint64 {
	hasher := fnv.New64a()
	hasher.Write(img.Pix)
	return hasher.Sum64()
}
