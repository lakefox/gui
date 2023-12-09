package fps

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// FPSCounter represents a simple FPS counter.
type FPSCounter struct {
	frameCount int
	fps        int
	frameTime  float32
}

// NewFPSCounter creates a new FPSCounter instance.
func NewFPSCounter() *FPSCounter {
	return &FPSCounter{}
}

// Update updates the FPS counter.
func (f *FPSCounter) Update() {
	f.frameCount++
	f.frameTime += rl.GetFrameTime()

	if f.frameTime >= 1.0 {
		f.fps = f.frameCount
		f.frameCount = 0
		f.frameTime = 0.0
	}
}

// Draw draws the FPS counter.
func (f *FPSCounter) Draw(x, y, fontSize int, color rl.Color) {
	text := fmt.Sprintf("FPS: %d", f.fps)
	textLen := rl.MeasureText(text, int32(fontSize))
	rl.DrawRectangle(int32(x)-int32(fontSize/2), int32(y)-int32(fontSize/2), textLen+int32(fontSize), int32(fontSize)+int32(fontSize), rl.Black)
	rl.DrawText(text, int32(x), int32(y), int32(fontSize), color)
}
