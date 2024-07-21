package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Canvas represents a drawing canvas
type Canvas struct {
	Context     *image.RGBA
	StrokeStyle color.RGBA
	FillStyle   color.RGBA
	path        []image.Point
	Font        *truetype.Font
}

// NewCanvas creates a new canvas with the specified width and height
func NewCanvas(width, height int) *Canvas {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	white := color.RGBA{255, 255, 255, 0}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	return &Canvas{
		Context:     img,
		StrokeStyle: color.RGBA{0, 0, 0, 255},
		FillStyle:   color.RGBA{0, 0, 0, 255},
	}
}

// SetFont sets the font for text rendering
func (c *Canvas) SetFont(fontPath string) error {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	c.Font = f
	return nil
}

// MoveTo starts a new sub-path at the given (x, y) coordinates
func (c *Canvas) MoveTo(x, y int) {
	c.path = append(c.path, image.Point{x, y})
}

// LineTo adds a line to the current path
func (c *Canvas) LineTo(x, y int) {
	if len(c.path) > 0 {
		lastPoint := c.path[len(c.path)-1]
		c.path = append(c.path, image.Point{x, y})
		drawLine(c.Context, lastPoint, image.Point{x, y}, c.StrokeStyle)
	}
}

// BeginPath starts a new path
func (c *Canvas) BeginPath() {
	c.path = nil
}

// Stroke draws the current path
func (c *Canvas) Stroke() {
	for i := 1; i < len(c.path); i++ {
		drawLine(c.Context, c.path[i-1], c.path[i], c.StrokeStyle)
	}
	c.BeginPath()
}

// Fill fills the current path
func (c *Canvas) Fill() {
	if len(c.path) > 0 {
		drawPolygon(c.Context, c.path, c.FillStyle)
	}
	c.BeginPath()
}

// Arc adds an arc to the current path
func (c *Canvas) Arc(x, y, radius, startAngle, endAngle float64) {
	steps := 100
	for i := 0; i <= steps; i++ {
		angle := startAngle + (endAngle-startAngle)*float64(i)/float64(steps)
		px := x + radius*math.Cos(angle)
		py := y + radius*math.Sin(angle)
		if i == 0 {
			c.MoveTo(int(px), int(py))
		} else {
			c.LineTo(int(px), int(py))
		}
	}
}

// Rect adds a rectangle to the path
func (c *Canvas) Rect(x, y, width, height int) {
	c.MoveTo(x, y)
	c.LineTo(x+width, y)
	c.LineTo(x+width, y+height)
	c.LineTo(x, y+height)
	c.LineTo(x, y)
}

// FillRect fills a rectangle
func (c *Canvas) FillRect(x, y, width, height int) {
	draw.Draw(c.Context, image.Rect(x, y, x+width, y+height), &image.Uniform{c.FillStyle}, image.Point{}, draw.Src)
}

// StrokeRect strokes a rectangle
func (c *Canvas) StrokeRect(x, y, width, height int) {
	c.Rect(x, y, width, height)
	c.Stroke()
}

// ClearRect clears a rectangle
func (c *Canvas) ClearRect(x, y, width, height int) {
	draw.Draw(c.Context, image.Rect(x, y, x+width, y+height), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)
}

// FillText draws filled text
func (c *Canvas) FillText(text string, x, y int, size float64) error {
	if c.Font == nil {
		return nil
	}
	face := truetype.NewFace(c.Font, &truetype.Options{Size: size})
	drawer := &font.Drawer{
		Dst:  c.Context,
		Src:  &image.Uniform{c.FillStyle},
		Face: face,
		Dot:  fixed.Point26_6{fixed.I(x), fixed.I(y)},
	}
	drawer.DrawString(text)
	return nil
}

// StrokeText draws stroked text
func (c *Canvas) StrokeText(text string, x, y int, size float64) error {
	// This example does not implement stroked text.
	// Implementing stroked text would require a more complex approach.
	return c.FillText(text, x, y, size)
}

// drawLine draws a line between two points on the image
func drawLine(img *image.RGBA, p1, p2 image.Point, col color.RGBA) {
	dx := math.Abs(float64(p2.X - p1.X))
	dy := math.Abs(float64(p2.Y - p1.Y))
	sx := -1
	if p1.X < p2.X {
		sx = 1
	}
	sy := -1
	if p1.Y < p2.Y {
		sy = 1
	}
	err := dx - dy

	for {
		img.Set(p1.X, p1.Y, col)
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}
	}
}

// drawPolygon fills a polygon defined by the points with the specified color
func drawPolygon(img *image.RGBA, points []image.Point, col color.RGBA) {
	for i := 1; i < len(points); i++ {
		drawLine(img, points[i-1], points[i], col)
	}
	if len(points) > 2 {
		drawLine(img, points[len(points)-1], points[0], col)
	}
}
