package scroll

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// ScrollBar represents a simple vertical or horizontal scrollbar.
type ScrollBar struct {
	position    float32
	trackSize   int32
	thumbSize   int32
	isVertical  bool
	isDragging  bool
	contentSize int32
	windowSize  int32
	dragOffset  float32
	mousePos    raylib.Vector2
	mouseOver   bool
}

// NewScrollBar creates a new ScrollBar instance.
func NewScrollBar(contentSize, windowSize int32, isVertical bool) *ScrollBar {
	sb := &ScrollBar{
		position:    0,
		trackSize:   16, // Adjust as needed
		thumbSize:   windowSize,
		isVertical:  isVertical,
		contentSize: contentSize,
		windowSize:  windowSize,
	}

	if !isVertical {
		sb.thumbSize = sb.windowSize
	}

	return sb
}

// Update updates the scrollbar.
func (sb *ScrollBar) Update(newWindowSize, newContentSize int32) float32 {
	sb.mousePos = raylib.GetMousePosition()

	thumbSize := float32(sb.windowSize) * (float32(sb.windowSize) / float32(sb.contentSize))
	println(thumbSize)
	if sb.isDragging {
		// Update position based on drag
		if sb.isVertical {
			sb.position = (sb.mousePos.Y - sb.dragOffset)
		} else {
			sb.position = (sb.mousePos.X - sb.dragOffset)
		}
	}
	if raylib.IsMouseButtonDown(raylib.MouseLeftButton) && !sb.isDragging && sb.insideScrollBar() {
		// Start dragging if mouse clicked on the thumb
		sb.isDragging = true
		if sb.isVertical {
			sb.dragOffset = sb.mousePos.Y - sb.position
			sb.mouseOver = sb.mousePos.Y > sb.position && sb.mousePos.Y < sb.position+thumbSize && sb.insideScrollBar()
		} else {
			sb.dragOffset = sb.mousePos.X - sb.position
			sb.mouseOver = sb.mousePos.X > sb.position && sb.mousePos.X < sb.position+thumbSize && sb.insideScrollBar()
		}
	}
	if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) && sb.isDragging {
		// Stop dragging
		sb.isDragging = false
		sb.mouseOver = false
	}

	// Update window size
	sb.windowSize = newWindowSize
	sb.contentSize = newContentSize

	// Clamp position within the track
	if sb.position < 0 {
		sb.position = 0
	} else if sb.position > float32(sb.windowSize)-thumbSize {
		sb.position = float32(sb.windowSize) - thumbSize
	}

	return -(sb.position * (float32(sb.contentSize / sb.windowSize)))
}

func (sb *ScrollBar) insideScrollBar() bool {
	if sb.isVertical {
		return sb.mousePos.X > float32(raylib.GetScreenWidth())-14
	} else {
		return sb.mousePos.Y > float32(raylib.GetScreenHeight())-14
	}
}

// Draw draws the scrollbar.
func (sb *ScrollBar) Draw() {
	var trackRect, thumbRect raylib.Rectangle

	thumbSize := float32(sb.windowSize) * (float32(sb.windowSize) / float32(sb.contentSize))

	if thumbSize >= float32(sb.windowSize) {
		return
	}

	if sb.isVertical {
		trackRect = raylib.NewRectangle(float32(raylib.GetScreenWidth())-16, 0, 16, float32(sb.windowSize))
		thumbRect = raylib.NewRectangle(float32(raylib.GetScreenWidth())-14, float32(sb.position), 12, float32(thumbSize))
	} else {
		trackRect = raylib.NewRectangle(0, float32(raylib.GetScreenHeight())-16, float32(sb.windowSize), 16)
		thumbRect = raylib.NewRectangle(float32(sb.position), float32(raylib.GetScreenHeight())-14, float32(thumbSize), 12)
	}

	// Draw track
	raylib.DrawRectangleRec(trackRect, raylib.LightGray)

	if sb.mouseOver {
		// Draw thumb
		raylib.DrawRectangleRec(thumbRect, raylib.DarkGray)
	} else {
		// Draw thumb
		raylib.DrawRectangleRec(thumbRect, raylib.Gray)
	}
}

// func main() {
// 	// Initialization
// 	const screenWidth = 800
// 	const screenHeight = 450

// 	raylib.InitWindow(screenWidth, screenHeight, "Scrollbar Example")
// 	raylib.SetTargetFPS(60)

// 	// Example content size
// 	contentWidth := int32(1000)
// 	contentHeight := int32(1000)

// 	// Create vertical scrollbar
// 	verticalScrollBar := NewScrollBar(contentHeight, screenHeight, true)

// 	// Create horizontal scrollbar
// 	horizontalScrollBar := NewScrollBar(contentWidth, screenWidth, false)

// 	for !raylib.WindowShouldClose() {
// 		// Update
// 		vO := verticalScrollBar.Update(screenHeight)
// 		hO := horizontalScrollBar.Update(screenWidth)

// 		// Draw
// 		raylib.BeginDrawing()

// 		raylib.ClearBackground(raylib.RayWhite)

// 		raylib.DrawRectangleLines(int32(0+hO), int32(0+vO), 1000, 1000, raylib.DarkGray)
// 		raylib.DrawText("Large Content", int32(10+hO), int32(10+vO), 20, raylib.DarkGray)

// 		// Draw vertical scrollbar
// 		verticalScrollBar.Draw()

// 		// Draw horizontal scrollbar
// 		horizontalScrollBar.Draw()

// 		raylib.EndDrawing()
// 	}

// 	// De-Initialization
// 	raylib.CloseWindow()
// }
