package raylib

import (
	adapter "gui/adapters"
	"gui/element"
	"slices"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Init() *adapter.Adapter {
	a := adapter.Adapter{}
	a.AddEventListener("cursor", func(e element.Event) {
		switch e.Data.(string) {
		case "":
			rl.SetMouseCursor(0)
		case "text":
			rl.SetMouseCursor(2)
		case "crosshair":
			rl.SetMouseCursor(3)
		case "pointer":
			rl.SetMouseCursor(4)
		case "ew-resize":
			rl.SetMouseCursor(5)
		case "ns-resize":
			rl.SetMouseCursor(6)
		case "nwse-resize":
			rl.SetMouseCursor(7)
		case "nesw-resize":
			rl.SetMouseCursor(8)
		case "grab":
			rl.SetMouseCursor(9)
		case "not-allowed":
			rl.SetMouseCursor(10)
		}
	})
	a.Options = adapter.Options{
		RenderText:     true,
		RenderElements: true,
		RenderBorders:  true,
	}
	wm := NewWindowManager(&a)
	a.Init = func(width, height int) {
		wm.OpenWindow(int32(width), int32(height))
		a.Library.UnloadCallback = func(key string) {
			t, exists := wm.Textures[key]
			if exists {
				rl.UnloadTexture(*t)
				delete(wm.Textures, key)
			}
		}
	}
	a.Load = wm.LoadTextures
	a.Render = func(state []element.State) {
		if rl.WindowShouldClose() {
			a.DispatchEvent(element.Event{Name: "close"})
		}
		wm.Draw(state)
	}
	return &a
}

// WindowManager manages the window and rectangles
type WindowManager struct {
	FPSCounterOn  bool
	FPS           int32
	Textures      map[string]*rl.Texture2D
	Width         int32
	Height        int32
	CurrentEvents map[int]bool
	MousePosition []int
	MouseState    bool
	ContextState  bool
	Adapter       *adapter.Adapter
}

// NewWindowManager creates a new WindowManager instance
func NewWindowManager(a *adapter.Adapter) *WindowManager {

	mp := rl.GetMousePosition()
	return &WindowManager{
		CurrentEvents: make(map[int]bool, 256),
		MousePosition: []int{int(mp.X), int(mp.Y)},
		Adapter:       a,
	}
}

// OpenWindow opens the window
func (wm *WindowManager) OpenWindow(width, height int32) {
	rl.InitWindow(width, height, "")
	rl.SetTargetFPS(30)
	wm.Width = width
	wm.Height = height
	// Enable window resizing
	rl.SetWindowState(rl.FlagWindowResizable)
}

func (wm *WindowManager) SetFPS(fps int) {
	wm.FPS = int32(fps)
	rl.SetTargetFPS(int32(fps))
}

func (wm *WindowManager) LoadTextures(nodes []element.State) {
	if wm.Textures == nil {
		wm.Textures = make(map[string]*rl.Texture2D)
	}

	for _, node := range nodes {
		if len(node.Textures) > 0 {
			for _, key := range node.Textures {
				rt, exists := wm.Textures[key]
				texture, inLibrary := wm.Adapter.Library.Get(key)
				matches := true
				if inLibrary && exists {
					tb := texture.Bounds()
					matches = (rt.Width == int32(tb.Dx()) && rt.Height == int32(tb.Dy()))
				}
				if (!exists && inLibrary) || !matches {
					textureLoaded := rl.LoadTextureFromImage(rl.NewImageFromImage(texture))
					wm.Textures[key] = &textureLoaded
				}
			}

		}
	}
}

// Draw draws all nodes on the window
func (wm *WindowManager) Draw(nodes []element.State) {
	indexes := []float32{0}
	rl.BeginDrawing()
	wm.GetEvents()
	for a := 0; a < len(indexes); a++ {
		for _, node := range nodes {
			if node.Hidden {
				continue
			}
			if node.Z == indexes[a] {
				// DrawRoundedRect(node.X,
				// 	node.Y,
				// 	node.Width+node.Border.Left.Width+node.Border.Right.Width,
				// 	node.Height+node.Border.Top.Width+node.Border.Bottom.Width,
				// 	node.Border.Radius.TopLeft, node.Border.Radius.TopRight, node.Border.Radius.BottomLeft, node.Border.Radius.BottomRight, node.Background)

				// Draw the border based on the style for each side

				if node.Textures != nil {
					for _, v := range node.Textures {
						texture, exists := wm.Textures[v]
						if exists {
							rl.DrawTexture(*texture, int32(node.X), int32(node.Y), rl.White)
						}
					}
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

	rl.EndDrawing()
}

func DrawRoundedRect(x, y, width, height float32, topLeftRadius, topRightRadius, bottomLeftRadius, bottomRightRadius float32, color rl.Color) {
	// Draw the main rectangle excluding corners
	rl.DrawRectangle(int32(x+topLeftRadius), int32(y), int32(width-topLeftRadius-topRightRadius), int32(height), color)
	rl.DrawRectangle(int32(x), int32(y+topLeftRadius), int32(topLeftRadius), int32(height-topLeftRadius-bottomLeftRadius), color)
	rl.DrawRectangle(int32(x+width-topRightRadius), int32(y+topRightRadius), int32(topRightRadius), int32(height-topRightRadius-bottomRightRadius), color)
	rl.DrawRectangle(int32(x+bottomLeftRadius), int32(y+height-bottomLeftRadius), int32(width-bottomLeftRadius-bottomRightRadius), int32(bottomLeftRadius), color)

	// Draw the corner circles
	rl.DrawCircleSector(rl.Vector2{X: x + topLeftRadius, Y: y + topLeftRadius}, topLeftRadius, 180, 270, 16, color)
	rl.DrawCircleSector(rl.Vector2{X: x + width - topRightRadius, Y: y + topRightRadius}, topRightRadius, 270, 360, 16, color)
	rl.DrawCircleSector(rl.Vector2{X: x + width - bottomRightRadius, Y: y + height - bottomRightRadius}, bottomRightRadius, 0, 90, 16, color)
	rl.DrawCircleSector(rl.Vector2{X: x + bottomLeftRadius, Y: y + height - bottomLeftRadius}, bottomLeftRadius, 90, 180, 16, color)

	// Draw rectangle parts to fill the gaps
	rl.DrawRectangle(int32(x+topLeftRadius), int32(y), int32(width-topLeftRadius-topRightRadius), int32(topLeftRadius), color)                                     // Top
	rl.DrawRectangle(int32(x), int32(y+topLeftRadius), int32(topLeftRadius), int32(height-topLeftRadius-bottomLeftRadius), color)                                  // Left
	rl.DrawRectangle(int32(x+width-topRightRadius), int32(y+topRightRadius), int32(topRightRadius), int32(height-topRightRadius-bottomRightRadius), color)         // Right
	rl.DrawRectangle(int32(x+bottomLeftRadius), int32(y+height-bottomLeftRadius), int32(width-bottomLeftRadius-bottomRightRadius), int32(bottomLeftRadius), color) // Bottom
}

func (wm *WindowManager) GetEvents() {
	cw := rl.GetScreenWidth()
	ch := rl.GetScreenHeight()
	if cw != int(wm.Width) || ch != int(wm.Height) {
		e := element.Event{
			Name: "windowresize",
			Data: map[string]int{"width": cw, "height": ch},
		}
		wm.Width = int32(cw)
		wm.Height = int32(ch)
		wm.Adapter.DispatchEvent(e)
	}

	// Other keys
	for i := 8; i <= 255; i++ {
		// for i := 32; i < 126; i++ {
		isDown := rl.IsKeyDown(int32(i))
		if wm.CurrentEvents[i] != isDown {
			if isDown {
				keydown := element.Event{
					Name: "keydown",
					Data: i,
				}

				wm.CurrentEvents[i] = true
				wm.Adapter.DispatchEvent(keydown)
			} else {
				keyup := element.Event{
					Name: "keyup",
					Data: i,
				}
				wm.CurrentEvents[i] = false
				wm.Adapter.DispatchEvent(keyup)
			}
		}
	}
	// mouse move, ctrl, shift etc

	mp := rl.GetMousePosition()
	if wm.MousePosition[0] != int(mp.X) || wm.MousePosition[1] != int(mp.Y) {
		wm.Adapter.DispatchEvent(element.Event{
			Name: "mousemove",
			Data: []int{int(mp.X), int(mp.Y)},
		})
		wm.MousePosition[0] = int(mp.X)
		wm.MousePosition[1] = int(mp.Y)
	}
	md := rl.IsMouseButtonDown(rl.MouseLeftButton)
	if md != wm.MouseState {
		if md {
			wm.Adapter.DispatchEvent(element.Event{
				Name: "mousedown",
			})
			wm.MouseState = true
		} else {
			wm.Adapter.DispatchEvent(element.Event{
				Name: "mouseup",
			})
			wm.MouseState = false
		}
	}

	cs := rl.IsMouseButtonPressed(rl.MouseRightButton)
	if cs != wm.ContextState {
		if cs {
			wm.Adapter.DispatchEvent(element.Event{
				Name: "contextmenudown",
			})
			wm.ContextState = true
		} else {
			wm.Adapter.DispatchEvent(element.Event{
				Name: "contextmenuup",
			})
			wm.ContextState = false
		}
	}

	wd := rl.GetMouseWheelMove()

	if wd != 0 {
		wm.Adapter.DispatchEvent(element.Event{
			Name: "scroll",
			Data: int(wd),
		})
	}
}
