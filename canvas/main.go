package canvas

import (
	"fmt"
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
	Path                     [][]image.Point
	Font                     *truetype.Font
	LineWidth                float64
	direction                string
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
	textAlign                string
	textBaseline             string
	textRendering            string
	wordSpacing              float64
	transforms               []map[string]float64
}

type Point struct {
	X float64
	Y float64
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
	c.Path[len(c.Path)-1] = append(c.Path[len(c.Path)-1], image.Point{x, y})
}

// LineTo adds a line to the current path
func (c *Canvas) LineTo(x, y int) {
	if len(c.Path[len(c.Path)-1]) > 0 {
		p1 := c.Path[len(c.Path)-1][len(c.Path[len(c.Path)-1])-1]
		p2 := image.Point{x, y}

		points := generatePoints(p1, p2)
		c.Path[len(c.Path)-1] = append(c.Path[len(c.Path)-1], points...)
	}
}

func generatePoints(p1, p2 image.Point) []image.Point {
	var points []image.Point

	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy))))

	xIncrement := float64(dx) / float64(steps)
	yIncrement := float64(dy) / float64(steps)

	x := float64(p1.X)
	y := float64(p1.Y)

	for i := 0; i <= steps; i++ {
		points = append(points, image.Point{int(math.Round(x)), int(math.Round(y))})
		x += xIncrement
		y += yIncrement
	}

	return points
}

// BeginPath starts a new path
func (c *Canvas) BeginPath() {
	// c.path = [][]image.Point{}
	c.Path = append(c.Path, []image.Point{})
}

// Stroke draws the current path
func (c *Canvas) Stroke() {
	c.runTransforms(len(c.Path) - 1)
	color := c.StrokeStyle
	// for i := 1; i < len(c.Path[len(c.Path)-1]); i++ {
	points := c.Path[len(c.Path)-1]
	for i := 0; i < len(points); i++ {
		// color.R = (c.StrokeStyle.R / uint8(len(c.path))) * uint8(i)
		// color.G = (c.StrokeStyle.G / uint8(len(c.path))) * uint8(i)
		// color.B = (c.StrokeStyle.B / uint8(len(c.path))) * uint8(i)
		xy := points[i]
		c.Context.Set(xy.X, xy.Y, color)
		// drawLine(c.Context, c.Path[len(c.Path)-1][i-1], c.Path[len(c.Path)-1][i], color, c.LineWidth)
	}
	// c.BeginPath()
}

// Fill fills the current path
func (c *Canvas) Fill() {
	// Idea is to stroke the shape then flood fill it but we need to find a start point
	// So go to the middle key in grid and find the middle value.
	// Are there issue with this method, yes
	c.Stroke()
	points := c.Path[len(c.Path)-1]
	img := c.Context
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	grid := make(map[int][]int)

	for i := 0; i < len(points); i++ {
		y := points[i].Y
		x := points[i].X
		if y > 0 && y < h && x > 0 && x < w {
			grid[y] = append(grid[y], x)
		}
	}
	keys := make([]int, 0, len(grid))
	for key := range grid {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Ints(keys)

	// Find the middle key
	centerY := int(len(keys) / 2)
	centerX := int(((grid[centerY][len(grid[centerY])-1] - grid[centerY][0]) / 2) + grid[centerY][0])

	if len(c.Path) > 0 {
		floodFill(img, image.Point{X: centerX, Y: centerY}, c.FillStyle)
	}
	c.BeginPath()
}

func (c *Canvas) Arc(x, y, radius, startAngle, endAngle float64, clockwise bool) {
	for ri := 0; ri <= int(c.LineWidth); ri++ {
		r := radius - float64(ri)
		var angleStep float64
		if clockwise {
			angleStep = (endAngle - startAngle) / float64(EstimateArcPixels(startAngle, endAngle, r))
		} else {
			angleStep = (startAngle - endAngle) / float64(EstimateArcPixels(startAngle, endAngle, r))
		}

		var lastX, lastY int
		for angle := 0.0; math.Abs(angle) < math.Abs(endAngle-startAngle); angle += angleStep {
			var currentAngle float64
			if clockwise {
				currentAngle = startAngle + angle
			} else {
				currentAngle = startAngle - angle
			}
			p := CalculatePoint(x, y, r, currentAngle)
			ix := int(p.X)
			iy := int(p.Y)
			if ix != lastX || iy != lastY {
				c.Path[len(c.Path)-1] = append(c.Path[len(c.Path)-1], image.Point{X: ix, Y: iy})
				lastX = ix
				lastY = iy
			}
		}
	}
}

func EstimateArcPixels(startAngle, stopAngle, radius float64) int {
	// Calculate the central angle in radians
	centralAngle := math.Abs(stopAngle - startAngle)

	// Ensure the central angle is within the range [0, 2π]
	if centralAngle < 0 {
		centralAngle += 2 * math.Pi
	}

	// Calculate the arc length
	arcLength := radius * centralAngle

	// Assuming each pixel approximately covers 1 unit length
	// Adjust the pixel density factor as needed for different resolutions
	pixelDensityFactor := 0.3 // This can be adjusted based on resolution

	// Estimate the number of pixels along the arc length
	numPixels := int(math.Round(arcLength / pixelDensityFactor))

	return numPixels
}

// EstimateBezierPoints estimates the required number of points for a cubic Bezier curve
func EstimateCubicBezierPoints(p0, p1, p2, p3 image.Point) int {
	// Helper function to calculate distance between two points
	distance := func(a, b image.Point) float64 {
		dx := float64(b.X - a.X)
		dy := float64(b.Y - a.Y)
		return math.Sqrt(dx*dx + dy*dy)
	}

	// Estimate the length of the cubic Bezier curve
	// Sum up the distances between successive control points
	approxCurveLength := distance(p0, p1) + distance(p1, p2) + distance(p2, p3)

	// Adjust the pixel density factor as needed for different resolutions
	// The higher the pixelDensityFactor, the fewer points will be used
	pixelDensityFactor := 0.3 // Default value if the input is invalid

	// Estimate the number of points along the curve length
	numPoints := int(math.Round(approxCurveLength / pixelDensityFactor))

	return numPoints
}

// EstimateQuadraticBezierPoints estimates the required number of points for a quadratic Bezier curve
func EstimateQuadraticBezierPoints(p0, p1, p2 image.Point) int {
	// Helper function to calculate distance between two points
	distance := func(a, b image.Point) float64 {
		dx := float64(b.X - a.X)
		dy := float64(b.Y - a.Y)
		return math.Sqrt(dx*dx + dy*dy)
	}

	// Estimate the length of the quadratic Bezier curve
	// Sum up the distances between successive control points
	approxCurveLength := distance(p0, p1) + distance(p1, p2)

	// Adjust the pixel density factor as needed for different resolutions
	// The higher the pixelDensityFactor, the fewer points will be used
	pixelDensityFactor := 0.3 // Default value if the input is invalid

	// Estimate the number of points along the curve length
	numPoints := int(math.Round(approxCurveLength / pixelDensityFactor))

	return numPoints
}

func CalculatePoint(cx, cy, radius, angle float64) Point {
	// Calculate x, y coordinates
	x := cx + radius*math.Cos(angle)
	y := cy + radius*math.Sin(angle)

	return Point{X: x, Y: y}
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
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
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
	c.Arc(float64(x+width-radius), float64(y+radius), float64(radius), 1.5*math.Pi, 2*math.Pi, true)
	c.LineTo(x+width, y+height-radius)
	c.Arc(float64(x+width-radius), float64(y+height-radius), float64(radius), 0, 0.5*math.Pi, true)
	c.LineTo(x+radius, y+height)
	c.Arc(float64(x+radius), float64(y+height-radius), float64(radius), 0.5*math.Pi, math.Pi, true)
	c.LineTo(x, y+radius)
	c.Arc(float64(x+radius), float64(y+radius), float64(radius), math.Pi, 1.5*math.Pi, true)
	c.ClosePath()
}

// ClosePath closes the current path
func (c *Canvas) ClosePath() {
	if len(c.Path) > 0 {
		c.Path = append(c.Path, c.Path[0])
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
func (c *Canvas) ArcTo(x1, y1, x2, y2, r float64, counterclockwise bool) {
	// Calculate the midpoint
	mx := (x1 + x2) / 2
	my := (y1 + y2) / 2

	// Calculate the distance between the points
	d := math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))

	// Check if the given radius is sufficient
	if d > 2*r {
		return
	} else {
		// Calculate the distance from the midpoint to the circle center
		h := math.Sqrt(math.Pow(r, 2) - math.Pow(d/2, 2))

		// Calculate the direction vector perpendicular to AB
		vx := -(y2 - y1)
		vy := x2 - x1

		// Normalize the direction vector
		length := math.Sqrt(vx*vx + vy*vy)
		nx := vx / length
		ny := vy / length

		// Calculate the two possible centers
		center1 := Point{X: mx + h*nx, Y: my + h*ny}
		center2 := Point{X: mx - h*nx, Y: my - h*ny}

		// Calculate the angles for the points A and B relative to the first center
		angleA1 := math.Atan2(y1-center1.Y, x1-center1.X)
		angleB1 := math.Atan2(y2-center1.Y, x2-center1.X)

		// Calculate the angles for the points A and B relative to the second center
		angleA2 := math.Atan2(y1-center2.Y, x1-center2.X)
		angleB2 := math.Atan2(y2-center2.Y, x2-center2.X)

		if counterclockwise {
			c.Arc(center1.X, center1.Y, r, angleA1, angleB1, true)
		} else {
			c.Arc(center2.X, center2.Y, r, angleA2, angleB2, true)

		}

	}

}
func degToRad(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}
func radToDeg(radians float64) float64 {
	return radians * (180.0 / math.Pi)
}

// BezierCurveTo adds a cubic Bezier curve to the path
func (c *Canvas) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	// Calculate Bezier curve points and add them to the path
	steps := 100
	p0 := c.Path[len(c.Path)-1][len(c.Path[len(c.Path)-1])-1]
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		mt := 1 - t
		xPos := math.Pow(mt, 3)*float64(p0.X) + 3*math.Pow(mt, 2)*t*cp1x + 3*mt*math.Pow(t, 2)*cp2x + math.Pow(t, 3)*x
		yPos := math.Pow(mt, 3)*float64(p0.Y) + 3*math.Pow(mt, 2)*t*cp1y + 3*mt*math.Pow(t, 2)*cp2y + math.Pow(t, 3)*y
		c.LineTo(int(xPos), int(yPos))
	}
}

// QuadraticBezier returns a slice of points on a quadratic Bezier curve
func QuadraticBezier(p0, p1, p2 image.Point) []image.Point {
	steps := EstimateQuadraticBezierPoints(p0, p1, p2)
	points := make([]image.Point, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := math.Pow(1-t, 2)*float64(p0.X) + 2*(1-t)*t*float64(p1.X) + math.Pow(t, 2)*float64(p2.X)
		y := math.Pow(1-t, 2)*float64(p0.Y) + 2*(1-t)*t*float64(p1.Y) + math.Pow(t, 2)*float64(p2.Y)
		points = append(points, image.Point{X: int(x), Y: int(y)})
	}
	return points
}

// CubicBezier returns a slice of points on a cubic Bezier curve
func CubicBezier(p0, p1, p2, p3 image.Point) []image.Point {
	steps := EstimateCubicBezierPoints(p0, p1, p2, p3)
	points := make([]image.Point, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := math.Pow(1-t, 3)*float64(p0.X) + 3*math.Pow(1-t, 2)*t*float64(p1.X) + 3*(1-t)*math.Pow(t, 2)*float64(p2.X) + math.Pow(t, 3)*float64(p3.X)
		y := math.Pow(1-t, 3)*float64(p0.Y) + 3*math.Pow(1-t, 2)*t*float64(p1.Y) + 3*(1-t)*math.Pow(t, 2)*float64(p2.Y) + math.Pow(t, 3)*float64(p3.Y)
		points = append(points, image.Point{X: int(x), Y: int(y)})
	}
	return points
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
		// Dst:  c.Context,
		// Src:  &image.Uniform{c.FillStyle},
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
	c.transforms = []map[string]float64{}
}

// ResetTransform resets the transformation matrix
func (c *Canvas) ResetTransform() {
	// ResetTransform implementation
	c.transforms = []map[string]float64{}
}

// Restore restores the most recently saved canvas state
func (c *Canvas) Restore() {
	// Restore implementation
}

// Rotate rotates the canvas around the given angle
func (c *Canvas) Rotate(angle float64) {
	// Rotate implementation
	c.transforms = append(c.transforms, map[string]float64{
		"type":    0, // 0 == rotate
		"radians": angle,
	})
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
	c.transforms = append(c.transforms, map[string]float64{
		"type": 1, // 1 == translate
		"x":    x,
		"y":    y,
	})
}

func (c *Canvas) runTransforms(path int) {
	contextPoint := image.Point{X: 0, Y: 0}
	for _, t := range c.transforms {
		switch t["type"] {
		case 0:
			// Rotate
			fmt.Println(contextPoint)
			c.Path[path] = rotatePoints(c.Path[path], contextPoint, t["radians"])
		case 1:
			// Translate
			for i := 0; i < len(c.Path[path]); i++ {
				c.Path[path][i].X += int(t["x"])
				c.Path[path][i].Y += int(t["y"])
			}
			contextPoint.X += int(t["x"])
			contextPoint.Y += int(t["y"])
		}
	}
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

// Function to check if a point is within the bounds of the image
func inBounds(img *image.RGBA, p image.Point) bool {
	return p.X >= 0 && p.X < img.Bounds().Dx() && p.Y >= 0 && p.Y < img.Bounds().Dy()
}

// Flood fill algorithm to fill a polygon from a starting point
func floodFill(img *image.RGBA, start image.Point, fillColor color.RGBA) {
	// Get the original color at the starting point
	targetColor := img.At(start.X, start.Y).(color.RGBA)

	// If the target color is the same as the fill color, do nothing
	if targetColor == fillColor {
		return
	}

	// Stack-based flood fill (to avoid stack overflow with recursive approach)
	stack := []image.Point{start}

	for len(stack) > 0 {
		// Pop the last point from the stack
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Check if the point is within the image bounds
		if !inBounds(img, p) {
			continue
		}

		// Get the color of the current point
		currentColor := img.At(p.X, p.Y).(color.RGBA)

		// If the color matches the target color, fill it with the new color
		if currentColor == targetColor {
			img.SetRGBA(p.X, p.Y, fillColor)

			// Add neighboring points to the stack
			stack = append(stack, image.Point{X: p.X + 1, Y: p.Y})
			stack = append(stack, image.Point{X: p.X - 1, Y: p.Y})
			stack = append(stack, image.Point{X: p.X, Y: p.Y + 1})
			stack = append(stack, image.Point{X: p.X, Y: p.Y - 1})
		}
	}
}

func rotatePoints(points []image.Point, center image.Point, rad float64) []image.Point {
	rotatedPoints := make([]image.Point, len(points))

	// Precompute sine and cosine of the angle
	cosTheta := math.Cos(rad)
	sinTheta := math.Sin(rad)

	for i, p := range points {
		// Translate point to origin
		translatedX := float64(p.X - center.X)
		translatedY := float64(p.Y - center.Y)

		// Apply rotation
		rotatedX := translatedX*cosTheta - translatedY*sinTheta
		rotatedY := translatedX*sinTheta + translatedY*cosTheta

		// Translate back to original location
		rotatedPoints[i] = image.Point{
			X: int(math.Round(rotatedX + float64(center.X))),
			Y: int(math.Round(rotatedY + float64(center.Y))),
		}
	}

	return rotatedPoints
}
