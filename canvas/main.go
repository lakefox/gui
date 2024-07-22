package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"sort"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Canvas represents a drawing canvas
type Canvas struct {
	Context                  *image.RGBA
	StrokeStyle              color.RGBA
	FillStyle                color.RGBA
	path                     []image.Point
	Font                     *truetype.Font
	LineWidth                float64
	direction                string
	fillStyle                color.RGBA
	filter                   string
	fontKerning              string
	fontStretch              string
	fontVariantCaps          string
	globalAlpha              float64
	globalCompositeOperation string
	imageSmoothingEnabled    bool
	imageSmoothingQuality    string
	letterSpacing            float64
	lineCap                  string
	lineDashOffset           float64
	lineJoin                 string
	miterLimit               float64
	shadowBlur               float64
	shadowColor              color.RGBA
	shadowOffsetX            float64
	shadowOffsetY            float64
	strokeStyle              color.RGBA
	textAlign                string
	textBaseline             string
	textRendering            string
	wordSpacing              float64
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
		LineWidth:   1.0,
		globalAlpha: 1.0,
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
		drawLine(c.Context, lastPoint, image.Point{x, y}, c.StrokeStyle, c.LineWidth)
	}
}

// BeginPath starts a new path
func (c *Canvas) BeginPath() {
	c.path = nil
}

// Stroke draws the current path
func (c *Canvas) Stroke() {
	for i := 1; i < len(c.path); i++ {
		drawLine(c.Context, c.path[i-1], c.path[i], c.StrokeStyle, c.LineWidth)
	}
	c.BeginPath()
}

// Fill fills the current path
func (c *Canvas) Fill() {
	if len(c.path) > 0 {
		fillPolygon(c.Context, c.path, c.FillStyle)
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
	return c.FillText(text, x, y, size)
}

// RoundedRect adds a rounded rectangle to the path
func (c *Canvas) RoundedRect(x, y, width, height, radius int) {
	c.BeginPath()
	c.MoveTo(x+radius, y)
	c.LineTo(x+width-radius, y)
	c.Arc(float64(x+width-radius), float64(y+radius), float64(radius), 1.5*math.Pi, 2*math.Pi)
	c.LineTo(x+width, y+height-radius)
	c.Arc(float64(x+width-radius), float64(y+height-radius), float64(radius), 0, 0.5*math.Pi)
	c.LineTo(x+radius, y+height)
	c.Arc(float64(x+radius), float64(y+height-radius), float64(radius), 0.5*math.Pi, math.Pi)
	c.LineTo(x, y+radius)
	c.Arc(float64(x+radius), float64(y+radius), float64(radius), math.Pi, 1.5*math.Pi)
	c.ClosePath()
}

// ClosePath closes the current path
func (c *Canvas) ClosePath() {
	if len(c.path) > 0 {
		c.path = append(c.path, c.path[0])
	}
}

// Setters for various properties
func (c *Canvas) SetDirection(dir string) {
	c.direction = dir
}

func (c *Canvas) SetFillStyle(col color.RGBA) {
	c.FillStyle = col
}

func (c *Canvas) SetFilter(filter string) {
	c.filter = filter
}

func (c *Canvas) SetFontKerning(kerning string) {
	c.fontKerning = kerning
}

func (c *Canvas) SetFontStretch(stretch string) {
	c.fontStretch = stretch
}

func (c *Canvas) SetFontVariantCaps(caps string) {
	c.fontVariantCaps = caps
}

func (c *Canvas) SetGlobalAlpha(alpha float64) {
	c.globalAlpha = alpha
}

func (c *Canvas) SetGlobalCompositeOperation(op string) {
	c.globalCompositeOperation = op
}

func (c *Canvas) SetImageSmoothingEnabled(enabled bool) {
	c.imageSmoothingEnabled = enabled
}

func (c *Canvas) SetImageSmoothingQuality(quality string) {
	c.imageSmoothingQuality = quality
}

func (c *Canvas) SetLetterSpacing(spacing float64) {
	c.letterSpacing = spacing
}

func (c *Canvas) SetLineCap(cap string) {
	c.lineCap = cap
}

func (c *Canvas) SetLineDashOffset(offset float64) {
	c.lineDashOffset = offset
}

func (c *Canvas) SetLineJoin(join string) {
	c.lineJoin = join
}

func (c *Canvas) SetLineWidth(width float64) {
	c.LineWidth = width
}

func (c *Canvas) SetMiterLimit(limit float64) {
	c.miterLimit = limit
}

func (c *Canvas) SetShadowBlur(blur float64) {
	c.shadowBlur = blur
}

func (c *Canvas) SetShadowColor(col color.RGBA) {
	c.shadowColor = col
}

func (c *Canvas) SetShadowOffsetX(offset float64) {
	c.shadowOffsetX = offset
}

func (c *Canvas) SetShadowOffsetY(offset float64) {
	c.shadowOffsetY = offset
}

func (c *Canvas) SetStrokeStyle(col color.RGBA) {
	c.StrokeStyle = col
}

func (c *Canvas) SetTextAlign(align string) {
	c.textAlign = align
}

func (c *Canvas) SetTextBaseline(baseline string) {
	c.textBaseline = baseline
}

func (c *Canvas) SetTextRendering(rendering string) {
	c.textRendering = rendering
}

func (c *Canvas) SetWordSpacing(spacing float64) {
	c.wordSpacing = spacing
}

// ArcTo adds an arc to the path with control points and radius
func (c *Canvas) ArcTo(x1, y1, x2, y2, radius float64) {
	// ArcTo implementation
}

// BezierCurveTo adds a cubic Bezier curve to the path
func (c *Canvas) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	// BezierCurveTo implementation
}

// Clip sets the clipping region to the current path
func (c *Canvas) Clip() {
	// Clip implementation
}

// CreateConicGradient creates a conic gradient
func (c *Canvas) CreateConicGradient(startAngle, x, y float64) {
	// CreateConicGradient implementation
}

// CreateImageData creates a new blank ImageData object
func (c *Canvas) CreateImageData(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

// CreateLinearGradient creates a linear gradient
func (c *Canvas) CreateLinearGradient(x0, y0, x1, y1 float64) {
	// CreateLinearGradient implementation
}

// CreatePattern creates a pattern with an image
func (c *Canvas) CreatePattern(img image.Image, repetition string) {
	// CreatePattern implementation
}

// CreateRadialGradient creates a radial gradient
func (c *Canvas) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) {
	// CreateRadialGradient implementation
}

// DrawFocusIfNeeded draws focus ring around element
func (c *Canvas) DrawFocusIfNeeded() {
	// DrawFocusIfNeeded implementation
}

// DrawImage draws an image onto the canvas
func (c *Canvas) DrawImage(img image.Image, dx, dy int) {
	draw.Draw(c.Context, img.Bounds().Add(image.Point{dx, dy}), img, image.Point{}, draw.Over)
}

// Ellipse adds an ellipse to the path
func (c *Canvas) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64) {
	// Ellipse implementation
}

// GetContextAttributes returns the context attributes
func (c *Canvas) GetContextAttributes() {
	// GetContextAttributes implementation
}

// GetImageData gets the image data for the specified rectangle
func (c *Canvas) GetImageData(sx, sy, sw, sh int) *image.RGBA {
	return c.Context.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*image.RGBA)
}

// GetLineDash returns the current line dash pattern
func (c *Canvas) GetLineDash() {
	// GetLineDash implementation
}

// GetTransform returns the current transformation matrix
func (c *Canvas) GetTransform() {
	// GetTransform implementation
}

// IsContextLost returns whether the context is lost
func (c *Canvas) IsContextLost() bool {
	return false
}

// IsPointInPath returns whether the given point is in the current path
func (c *Canvas) IsPointInPath(x, y float64) bool {
	// IsPointInPath implementation
	return false
}

// IsPointInStroke returns whether the given point is in the current stroke
func (c *Canvas) IsPointInStroke(x, y float64) bool {
	// IsPointInStroke implementation
	return false
}

// MeasureText measures the width of the given text
func (c *Canvas) MeasureText(text string) (float64, float64) {
	if c.Font == nil {
		return 0, 0
	}
	face := truetype.NewFace(c.Font, &truetype.Options{Size: 12})
	drawer := &font.Drawer{
		Dst:  c.Context,
		Src:  &image.Uniform{c.FillStyle},
		Face: face,
	}
	bounds, _ := font.BoundString(drawer.Face, text)
	width := float64((bounds.Max.X - bounds.Min.X).Ceil())
	height := float64((bounds.Max.Y - bounds.Min.Y).Ceil())
	return width, height
}

// PutImageData puts the image data onto the canvas
func (c *Canvas) PutImageData(img *image.RGBA, dx, dy int) {
	draw.Draw(c.Context, img.Bounds().Add(image.Point{dx, dy}), img, image.Point{}, draw.Over)
}

// QuadraticCurveTo adds a quadratic Bezier curve to the path
func (c *Canvas) QuadraticCurveTo(cpx, cpy, x, y float64) {
	// QuadraticCurveTo implementation
}

// Reset resets the canvas state
func (c *Canvas) Reset() {
	// Reset implementation
}

// ResetTransform resets the transformation matrix
func (c *Canvas) ResetTransform() {
	// ResetTransform implementation
}

// Restore restores the most recently saved canvas state
func (c *Canvas) Restore() {
	// Restore implementation
}

// Rotate rotates the canvas around the given angle
func (c *Canvas) Rotate(angle float64) {
	// Rotate implementation
}

// RoundRect adds a rounded rectangle to the path
func (c *Canvas) RoundRect(x, y, w, h, radius float64) {
	c.MoveTo(int(x+radius), int(y))
	c.LineTo(int(x+w-radius), int(y))
	c.Arc(x+w-radius, y+radius, radius, 1.5*math.Pi, 2*math.Pi)
	c.LineTo(int(x+w), int(y+h-radius))
	c.Arc(x+w-radius, y+h-radius, radius, 0, 0.5*math.Pi)
	c.LineTo(int(x+radius), int(y+h))
	c.Arc(x+radius, y+h-radius, radius, 0.5*math.Pi, math.Pi)
	c.LineTo(int(x), int(y+radius))
	c.Arc(x+radius, y+radius, radius, math.Pi, 1.5*math.Pi)
	c.ClosePath()
}

// Scale scales the canvas by the given factors
func (c *Canvas) Scale(x, y float64) {
	// Scale implementation
}

// SetLineDash sets the line dash pattern
func (c *Canvas) SetLineDash(dash []float64) {
	// SetLineDash implementation
}

// SetTransform sets the transformation matrix
func (c *Canvas) SetTransform(a, b, c1, d, e, f float64) {
	// SetTransform implementation
}

// Transform applies the transformation matrix
func (c *Canvas) Transform(a, b, c1, d, e, f float64) {
	// Transform implementation
}

// Translate translates the canvas by the given distances
func (c *Canvas) Translate(x, y float64) {
	// Translate implementation
}

// SavePNG saves the canvas to a PNG file
func (c *Canvas) SavePNG(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, c.Context)
}

// Helper functions
func drawLine(img *image.RGBA, p1, p2 image.Point, col color.RGBA, lineWidth float64) {
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

	halfWidth := int(math.Ceil(lineWidth / 2))

	for i := -halfWidth; i <= halfWidth; i++ {
		xOffset := 0
		yOffset := 0
		if dx > dy {
			yOffset = i
		} else {
			xOffset = i
		}

		x1, y1 := p1.X+xOffset, p1.Y+yOffset
		x2, y2 := p2.X+xOffset, p2.Y+yOffset

		err := dx - dy // Reset err for each parallel line
		for {
			img.Set(x1, y1, col)
			if x1 == x2 && y1 == y2 {
				break
			}
			e2 := 2 * err // Calculate the double of err
			if e2 > -dy {
				err -= dy // Adjust err and move in x direction
				x1 += sx
			}
			if e2 < dx {
				err += dx // Adjust err and move in y direction
				y1 += sy
			}
		}
	}
}

func fillPolygon(img *image.RGBA, points []image.Point, col color.RGBA) {
	if len(points) < 3 {
		return
	}

	type Edge struct {
		yMin     int
		yMax     int
		xAtYMin  float64
		invSlope float64
	}

	var edges []Edge
	n := len(points)

	// Build the edge table
	for i := 0; i < n; i++ {
		p1 := points[i]
		p2 := points[(i+1)%n]

		if p1.Y == p2.Y {
			continue // Skip horizontal edges
		}

		if p1.Y > p2.Y {
			p1, p2 = p2, p1
		}

		invSlope := float64(p2.X-p1.X) / float64(p2.Y-p1.Y)
		edges = append(edges, Edge{p1.Y, p2.Y, float64(p1.X), invSlope})
	}

	// Sort edges by yMin, then xAtYMin
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].yMin == edges[j].yMin {
			return edges[i].xAtYMin < edges[j].xAtYMin
		}
		return edges[i].yMin < edges[j].yMin
	})

	// Scanline fill
	activeEdges := []Edge{}
	y := edges[0].yMin
	for i := 0; i < len(edges) || len(activeEdges) > 0; y++ {
		// Add edges to active edge list
		for i < len(edges) && edges[i].yMin == y {
			activeEdges = append(activeEdges, edges[i])
			i++
		}

		// Remove edges from active edge list where yMax == y
		newActiveEdges := activeEdges[:0]
		for _, edge := range activeEdges {
			if edge.yMax != y {
				newActiveEdges = append(newActiveEdges, edge)
			}
		}
		activeEdges = newActiveEdges

		// Sort active edges by xAtYMin
		sort.Slice(activeEdges, func(i, j int) bool {
			return activeEdges[i].xAtYMin < activeEdges[j].xAtYMin
		})

		// Fill pixels between pairs of intersections
		for j := 0; j < len(activeEdges); j += 2 {
			xStart := int(math.Ceil(activeEdges[j].xAtYMin))
			xEnd := int(math.Floor(activeEdges[j+1].xAtYMin))

			for x := xStart; x <= xEnd; x++ {
				img.Set(x, y, col)
			}

			activeEdges[j].xAtYMin += activeEdges[j].invSlope
			activeEdges[j+1].xAtYMin += activeEdges[j+1].invSlope
		}
	}
}
