package border

import (
	"gui/canvas"
	"gui/color"
	"gui/element"
	"gui/library"
	"gui/utils"
	"image"
	ic "image/color"
	"math"
	"strconv"
	"strings"
)

func Parse(cssProperties map[string]string, self, parent element.State) (element.Border, error) {
	// Define default values
	defaultWidth := "0px"
	defaultStyle := "solid"
	defaultColor := "#000000"
	defaultRadius := "0px"

	// Helper function to parse border component
	parseBorderComponent := func(value string) (width, style, color string) {
		components := strings.Fields(value)
		width, style, color = defaultWidth, defaultStyle, defaultColor
		widthSuffixes := []string{"px", "em", "pt", "pc", "%", "vw", "vh", "cm", "in"}

		for _, component := range components {
			if isWidthComponent(component, widthSuffixes) {
				width = component
			} else {
				switch component {
				case "thin", "medium", "thick":
					width = component
				case "none", "hidden", "dotted", "dashed", "solid", "double", "groove", "ridge", "inset", "outset":
					style = component
				default:
					color = component
				}
			}
		}

		return
	}

	// Helper function to parse border radius component
	parseBorderRadiusComponent := func(value string) float32 {
		if value == "" {
			value = defaultRadius
		}
		return utils.ConvertToPixels(value, self.EM, parent.Width)
	}

	// Parse individual border sides
	topWidth, topStyle, topColor := parseBorderComponent(cssProperties["border-top"])
	rightWidth, rightStyle, rightColor := parseBorderComponent(cssProperties["border-right"])
	bottomWidth, bottomStyle, bottomColor := parseBorderComponent(cssProperties["border-bottom"])
	leftWidth, leftStyle, leftColor := parseBorderComponent(cssProperties["border-left"])

	// Parse shorthand border property
	if border, exists := cssProperties["border"]; exists {
		width, style, color := parseBorderComponent(border)
		if _, exists := cssProperties["border-top"]; !exists {
			topWidth, topStyle, topColor = width, style, color
		}
		if _, exists := cssProperties["border-right"]; !exists {
			rightWidth, rightStyle, rightColor = width, style, color
		}
		if _, exists := cssProperties["border-bottom"]; !exists {
			bottomWidth, bottomStyle, bottomColor = width, style, color
		}
		if _, exists := cssProperties["border-left"]; !exists {
			leftWidth, leftStyle, leftColor = width, style, color
		}
	}

	var topLeftRadius,
		topRightRadius,
		bottomLeftRadius,
		bottomRightRadius float32

	if cssProperties["border-radius"] != "" {
		rad := parseBorderRadiusComponent(cssProperties["border-radius"])
		topLeftRadius = rad
		topRightRadius = rad
		bottomLeftRadius = rad
		bottomRightRadius = rad
	}

	// Parse border-radius
	if cssProperties["border-top-left-radius"] != "" {
		topLeftRadius = parseBorderRadiusComponent(cssProperties["border-top-left-radius"])
	}
	if cssProperties["border-top-right-radius"] != "" {
		topRightRadius = parseBorderRadiusComponent(cssProperties["border-top-right-radius"])
	}
	if cssProperties["border-bottom-left-radius"] != "" {
		bottomLeftRadius = parseBorderRadiusComponent(cssProperties["border-bottom-left-radius"])
	}
	if cssProperties["border-bottom-right-radius"] != "" {
		bottomRightRadius = parseBorderRadiusComponent(cssProperties["border-bottom-right-radius"])
	}

	// Convert to pixels
	topWidthPx := utils.ConvertToPixels(topWidth, self.EM, parent.Width)
	rightWidthPx := utils.ConvertToPixels(rightWidth, self.EM, parent.Width)
	bottomWidthPx := utils.ConvertToPixels(bottomWidth, self.EM, parent.Width)
	leftWidthPx := utils.ConvertToPixels(leftWidth, self.EM, parent.Width)

	// Parse colors
	topParsedColor := color.Color(topColor)
	rightParsedColor := color.Color(rightColor)
	bottomParsedColor := color.Color(bottomColor)
	leftParsedColor := color.Color(leftColor)

	width := self.Width + self.Border.Left.Width + self.Border.Right.Width
	height := self.Height + self.Border.Top.Width + self.Border.Bottom.Width

	if topLeftRadius+topRightRadius > width {
		topLeftRadius = width / 2
		topRightRadius = width / 2
	}
	if bottomLeftRadius+bottomRightRadius > width {
		bottomLeftRadius = width / 2
		bottomRightRadius = width / 2
	}
	if topLeftRadius+bottomLeftRadius > height {
		topLeftRadius = height / 2
		bottomLeftRadius = height / 2
	}
	if topRightRadius+bottomRightRadius > height {
		topRightRadius = height / 2
		bottomRightRadius = height / 2
	}

	return element.Border{
		Top: element.BorderSide{
			Width: topWidthPx,
			Style: topStyle,
			Color: topParsedColor,
		},
		Right: element.BorderSide{
			Width: rightWidthPx,
			Style: rightStyle,
			Color: rightParsedColor,
		},
		Bottom: element.BorderSide{
			Width: bottomWidthPx,
			Style: bottomStyle,
			Color: bottomParsedColor,
		},
		Left: element.BorderSide{
			Width: leftWidthPx,
			Style: leftStyle,
			Color: leftParsedColor,
		},
		Radius: element.BorderRadius{
			TopLeft:     topLeftRadius,
			TopRight:    topRightRadius,
			BottomLeft:  bottomLeftRadius,
			BottomRight: bottomRightRadius,
		},
	}, nil
}

func Draw(n *element.State, shelf *library.Shelf) {
	// lastChange := time.Now()
	if n.Border.Top.Width > 0 ||
		n.Border.Right.Width > 0 ||
		n.Border.Bottom.Width > 0 ||
		n.Border.Left.Width > 0 {

		// Format: widthheightborderdatatopleftbottomright
		// borderdata: widthstylecolorradius
		// 50020020solid#fff520solid#fff520solid#fff520solid#fff520solid#fff
		key := strconv.Itoa(int(n.Width)) + strconv.Itoa(int(n.Height)) + (strconv.Itoa(int(n.Border.Top.Width)) + n.Border.Top.Style + utils.RGBAtoString(n.Border.Top.Color) + strconv.Itoa(int(n.Border.Radius.TopLeft))) + (strconv.Itoa(int(n.Border.Left.Width)) + n.Border.Left.Style + utils.RGBAtoString(n.Border.Left.Color) + strconv.Itoa(int(n.Border.Radius.BottomLeft))) + (strconv.Itoa(int(n.Border.Bottom.Width)) + n.Border.Bottom.Style + utils.RGBAtoString(n.Border.Bottom.Color) + strconv.Itoa(int(n.Border.Radius.BottomRight))) + (strconv.Itoa(int(n.Border.Right.Width)) + n.Border.Right.Style + utils.RGBAtoString(n.Border.Right.Color) + strconv.Itoa(int(n.Border.Radius.TopRight)))
		exists := shelf.Check(key)

		if exists {
			// Convert slice to a map for faster lookup
			lookup := make(map[string]struct{}, len(n.Textures))
			for _, v := range n.Textures {
				lookup[v] = struct{}{}
			}

			if _, found := lookup[key]; !found {
				n.Textures = append(n.Textures, key)
			}
		} else {
			ctx := canvas.NewCanvas(int(n.X+
				n.Width+n.Border.Left.Width+n.Border.Right.Width),
				int(n.Y+n.Height+n.Border.Top.Width+n.Border.Bottom.Width))
			ctx.StrokeStyle = ic.RGBA{0, 0, 0, 255}
			if n.Border.Top.Width > 0 {
				drawBorderSide(ctx, "top", n.Border.Top, n)
			}
			if n.Border.Right.Width > 0 {
				drawBorderSide(ctx, "right", n.Border.Right, n)
			}
			if n.Border.Bottom.Width > 0 {
				drawBorderSide(ctx, "bottom", n.Border.Bottom, n)
			}
			if n.Border.Left.Width > 0 {
				drawBorderSide(ctx, "left", n.Border.Left, n)
			}
			n.Textures = append(n.Textures, shelf.Set(key, ctx.Context))
		}

	}

	// fmt.Println(time.Since(lastChange))
}

func drawBorderSide(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	switch border.Style {
	case "solid":
		drawSolidBorder(ctx, side, border, s)
	case "dashed":
		drawDashedBorder(ctx, side, border, s)
	case "dotted":
		drawDottedBorder(ctx, side, border, s)
	case "double":
		drawDoubleBorder(ctx, side, border, s)
	case "groove":
		drawGrooveBorder(ctx, side, border, s)
	case "ridge":
		drawRidgeBorder(ctx, side, border, s)
	case "inset":
		drawInsetBorder(ctx, side, border, s)
	case "outset":
		drawOutsetBorder(ctx, side, border, s)
	default:
		drawSolidBorder(ctx, side, border, s)
	}
}

// Helper function to determine if a component is a width value
func isWidthComponent(component string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(component, suffix) {
			return true
		}
	}
	return false
}
func degToRad(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func drawSolidBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	ctx.FillStyle = border.Color
	ctx.StrokeStyle = border.Color

	width := s.Width + s.Border.Left.Width + s.Border.Right.Width
	height := s.Height + s.Border.Top.Width + s.Border.Bottom.Width
	ctx.BeginPath()
	switch side {
	case "top":
		// fmt.Println(1)
		v1 := math.Max(float64(s.Border.Radius.TopLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.TopRight), 1)
		genSolidBorder(ctx, width, v1, v2, border, s.Border.Left, s.Border.Right)
	case "right":
		// fmt.Println(2)
		v1 := math.Max(float64(s.Border.Radius.TopRight), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomRight), 1)
		genSolidBorder(ctx, height, v1, v2, border, s.Border.Top, s.Border.Bottom)

		ctx.Translate(float64(width), 0)
		ctx.Rotate(math.Pi / 2)
	case "bottom":
		// fmt.Println(3)
		v1 := math.Max(float64(s.Border.Radius.BottomLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomRight), 1)
		genSolidBorder(ctx, width, v2, v1, border, s.Border.Right, s.Border.Left)

		ctx.Translate(float64(width), float64(height))
		ctx.Rotate(math.Pi)
	case "left":
		// fmt.Println(4)
		v1 := math.Max(float64(s.Border.Radius.TopLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomLeft), 1)
		genSolidBorder(ctx, height, v2, v1, border, s.Border.Bottom, s.Border.Top)

		ctx.Translate(0, float64(height))
		ctx.Rotate((-math.Pi / 2))
	}
	ctx.Fill()
	// ctx.Stroke()
	ctx.Reset()
	ctx.ClosePath()
}

func genSolidBorder(ctx *canvas.Canvas, width float32, v1, v2 float64, border, side1, side2 element.BorderSide) {
	s1w, s2w := float32(side1.Width), float32(side2.Width)
	if s1w < 1 {
		s1w = 1
	}
	if s2w < 1 {
		s2w = 1
	}
	startAngleLeft := FindBorderStopAngle(image.Point{X: 0, Y: 0}, image.Point{X: int(s1w), Y: int(border.Width)}, image.Point{X: int(v1), Y: int(v1)}, v1)

	ctx.Arc(v1, v1, v1, -math.Pi/2, startAngleLeft[0]-math.Pi, false) // top arc left
	lineStart := ctx.Path[len(ctx.Path)-1]

	anglePercent := ((startAngleLeft[0] - math.Pi) - (-math.Pi / 2)) / ((-math.Pi) - (-math.Pi / 2))
	midBorderWidth := math.Max(math.Abs(float64(border.Width)-float64(s1w)), math.Min(float64(border.Width), float64(s1w)))

	bottomBorderEnd := FindPointOnLine(image.Point{X: 0, Y: 0}, lineStart, midBorderWidth)

	var lineEnd, parallelLineStart, parallelLineEnd image.Point
	if v1 <= float64(border.Width) {
		lineEnd = image.Point{X: int(s1w), Y: int(border.Width)}
		parallelLineStart = lineEnd

	} else {
		controlPoint1 := FindQuadraticControlPoint(image.Point{X: int(v1), Y: int(border.Width)}, bottomBorderEnd, image.Point{X: int(s1w), Y: int(v1)}, anglePercent)

		endPoints := canvas.QuadraticBezier(image.Point{X: int(v1), Y: int(border.Width)}, controlPoint1, image.Point{X: int(s1w), Y: int(v1)})

		var index int
		lineEnd, index, _ = FindClosestPoint(endPoints, bottomBorderEnd)
		parallelLineStart = endPoints[0]
		ctx.Path = append(ctx.Path, endPoints[:index]...)
	}

	ctx.MoveTo(lineStart.X, lineStart.Y)
	ctx.LineTo(lineEnd.X, lineEnd.Y) // cap the end of the arc (other side)

	// These are the paralle lines pt 1
	ctx.MoveTo(int(float64(width)-(v2)), 0)
	ctx.LineTo(int(v1), 0)

	// Move to start the second corner
	ctx.MoveTo(int(float64(width)-(v2)), 0)

	startAngleRight := FindBorderStopAngle(image.Point{X: int(width), Y: 0}, image.Point{X: int(width - s2w), Y: int(border.Width)}, image.Point{X: int(float64(width) - v2), Y: int(v2)}, v2)
	ctx.Arc(float64(width)-v2, v2, v2, -math.Pi/2, startAngleRight[0]-math.Pi, false) // top arc Right
	lineStart = ctx.Path[len(ctx.Path)-1]
	anglePercent = math.Abs(((startAngleRight[0] - math.Pi) - (-math.Pi / 2)) / ((-math.Pi) - (-math.Pi / 2)))
	midBorderWidth = math.Max(math.Abs(float64(border.Width)-float64(s2w)), math.Min(float64(border.Width), float64(s2w)))

	bottomBorderEnd = FindPointOnLine(image.Point{X: int(width), Y: 0}, lineStart, midBorderWidth)

	if v2 <= float64(border.Width) {
		lineEnd = image.Point{X: int(width - s2w), Y: int(border.Width)}
		parallelLineEnd = lineEnd
	} else {
		controlPoint := FindQuadraticControlPoint(image.Point{X: int(float64(width) - v2), Y: int(border.Width)}, bottomBorderEnd, image.Point{X: int(width - s2w), Y: int(v2)}, anglePercent)
		endPoints := canvas.QuadraticBezier(image.Point{X: int(float64(width) - v2), Y: int(border.Width)}, controlPoint, image.Point{X: int(width - s2w), Y: int(v2)})
		var index int
		lineEnd, index, _ = FindClosestPoint(endPoints, bottomBorderEnd)
		parallelLineEnd = endPoints[0]
		ctx.Path = append(ctx.Path, endPoints[:index]...)
	}

	// These are the paralle lines pt 2
	ctx.MoveTo(parallelLineStart.X, parallelLineStart.Y)
	ctx.LineTo(parallelLineEnd.X, parallelLineEnd.Y)

	ctx.MoveTo(lineStart.X, lineStart.Y)
	ctx.LineTo(lineEnd.X, lineEnd.Y) // cap the end of the arc (other side)
}

func drawDashedBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	ctx.FillStyle = border.Color
	ctx.StrokeStyle = border.Color

	width := s.Width + s.Border.Left.Width + s.Border.Right.Width
	height := s.Height + s.Border.Top.Width + s.Border.Bottom.Width
	ctx.BeginPath()
	switch side {
	case "top":
		v1 := math.Max(float64(s.Border.Radius.TopLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.TopRight), 1)
		genSolidBorder(ctx, width, v1, v2, border, s.Border.Left, s.Border.Right)
	case "right":
		v1 := math.Max(float64(s.Border.Radius.TopRight), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomRight), 1)
		genSolidBorder(ctx, height, v1, v2, border, s.Border.Top, s.Border.Bottom)

		ctx.Translate(float64(width), 0)
		ctx.Rotate(math.Pi / 2)
	case "bottom":
		v1 := math.Max(float64(s.Border.Radius.BottomLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomRight), 1)
		genSolidBorder(ctx, width, v2, v1, border, s.Border.Right, s.Border.Left)

		ctx.Translate(float64(width), float64(height))
		ctx.Rotate(math.Pi)
	case "left":
		// Top Right
		v1 := math.Max(float64(s.Border.Radius.TopLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomLeft), 1)
		genSolidBorder(ctx, height, v2, v1, border, s.Border.Top, s.Border.Bottom)

		ctx.Translate(0, float64(height))
		ctx.Rotate(-math.Pi / 2)
	}
	ctx.Fill()
	// ctx.Stroke()
	ctx.Reset()
	ctx.ClosePath()
}

func drawDottedBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// dotSize := int(border.Width)
	// drawDot := func(x, y int) {
	// 	for i := 0; i < dotSize; i++ {
	// 		for j := 0; j < dotSize; j++ {
	// 			img.Set(x+i, y+j, border.Color)
	// 		}
	// 	}
	// }
	// switch side {
	// case "top":
	// 	for i := rect.Min.X; i < rect.Max.X; i += dotSize * 2 {
	// 		drawDot(i, rect.Min.Y)
	// 	}
	// case "right":
	// 	for i := rect.Min.Y; i < rect.Max.Y; i += dotSize * 2 {
	// 		drawDot(rect.Max.X-dotSize, i)
	// 	}
	// case "bottom":
	// 	for i := rect.Min.X; i < rect.Max.X; i += dotSize * 2 {
	// 		drawDot(i, rect.Max.Y-dotSize)
	// 	}
	// case "left":
	// 	for i := rect.Min.Y; i < rect.Max.Y; i += dotSize * 2 {
	// 		drawDot(rect.Min.X, i)
	// 	}
	// }
}

func drawDoubleBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// innerOffset := int(border.Width / 3)
	// outerOffset := innerOffset * 2
	// drawLine := func(x1, y1, x2, y2, offset int) {
	// 	for i := 0; i < innerOffset; i++ {
	// 		drawLineHelper(img, x1, y1+i+offset, x2, y2+i+offset, border.Color)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, 0)
	// 	drawLine(rect.Min.X, rect.Min.Y+outerOffset, rect.Max.X, rect.Min.Y+outerOffset, 0)
	// case "right":
	// 	drawLine(rect.Max.X-innerOffset, rect.Min.Y, rect.Max.X-innerOffset, rect.Max.Y, 0)
	// 	drawLine(rect.Max.X-outerOffset, rect.Min.Y, rect.Max.X-outerOffset, rect.Max.Y, 0)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-innerOffset, rect.Max.X, rect.Max.Y-innerOffset, 0)
	// 	drawLine(rect.Min.X, rect.Max.Y-outerOffset, rect.Max.X, rect.Max.Y-outerOffset, 0)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, 0)
	// 	drawLine(rect.Min.X+outerOffset, rect.Min.Y, rect.Min.X+outerOffset, rect.Max.Y, 0)
	// }
}

func drawGrooveBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// shadowColor := ic.RGBA{border.Color.R / 2, border.Color.G / 2, border.Color.B / 2, border.Color.A}
	// highlightColor := ic.RGBA{border.Color.R * 2 / 3, border.Color.G * 2 / 3, border.Color.B * 2 / 3, border.Color.A}
	// drawLine := func(x1, y1, x2, y2 int, col ic.RGBA) {
	// 	for i := 0; i < int(border.Width/2); i++ {
	// 		drawLineHelper(img, x1, y1+i, x2, y2+i, col)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, shadowColor)
	// 	drawLine(rect.Min.X, rect.Min.Y+int(border.Width/2), rect.Max.X, rect.Min.Y+int(border.Width/2), highlightColor)
	// case "right":
	// 	drawLine(rect.Max.X-int(border.Width/2), rect.Min.Y, rect.Max.X-int(border.Width/2), rect.Max.Y, shadowColor)
	// 	drawLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y, highlightColor)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width/2), rect.Max.X, rect.Max.Y-int(border.Width/2), highlightColor)
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width), shadowColor)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, highlightColor)
	// 	drawLine(rect.Min.X+int(border.Width/2), rect.Min.Y, rect.Min.X+int(border.Width/2), rect.Max.Y, shadowColor)
	// }
}

func drawRidgeBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// shadowColor := ic.RGBA{border.Color.R / 2, border.Color.G / 2, border.Color.B / 2, border.Color.A}
	// highlightColor := ic.RGBA{border.Color.R * 2 / 3, border.Color.G * 2 / 3, border.Color.B * 2 / 3, border.Color.A}
	// drawLine := func(x1, y1, x2, y2 int, col ic.RGBA) {
	// 	for i := 0; i < int(border.Width/2); i++ {
	// 		drawLineHelper(img, x1, y1+i, x2, y2+i, col)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, highlightColor)
	// 	drawLine(rect.Min.X, rect.Min.Y+int(border.Width/2), rect.Max.X, rect.Min.Y+int(border.Width/2), shadowColor)
	// case "right":
	// 	drawLine(rect.Max.X-int(border.Width/2), rect.Min.Y, rect.Max.X-int(border.Width/2), rect.Max.Y, highlightColor)
	// 	drawLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y, shadowColor)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width/2), rect.Max.X, rect.Max.Y-int(border.Width/2), shadowColor)
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width), highlightColor)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, shadowColor)
	// 	drawLine(rect.Min.X+int(border.Width/2), rect.Min.Y, rect.Min.X+int(border.Width/2), rect.Max.Y, highlightColor)
	// }
}

func drawInsetBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// shadowColor := ic.RGBA{border.Color.R / 2, border.Color.G / 2, border.Color.B / 2, border.Color.A}
	// highlightColor := ic.RGBA{border.Color.R * 2 / 3, border.Color.G * 2 / 3, border.Color.B * 2 / 3, border.Color.A}
	// drawLine := func(x1, y1, x2, y2 int, col ic.RGBA) {
	// 	for i := 0; i < int(border.Width/2); i++ {
	// 		drawLineHelper(img, x1, y1+i, x2, y2+i, col)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, shadowColor)
	// 	drawLine(rect.Min.X, rect.Min.Y+int(border.Width/2), rect.Max.X, rect.Min.Y+int(border.Width/2), highlightColor)
	// case "right":
	// 	drawLine(rect.Max.X-int(border.Width/2), rect.Min.Y, rect.Max.X-int(border.Width/2), rect.Max.Y, shadowColor)
	// 	drawLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y, highlightColor)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width/2), rect.Max.X, rect.Max.Y-int(border.Width/2), highlightColor)
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width), shadowColor)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, highlightColor)
	// 	drawLine(rect.Min.X+int(border.Width/2), rect.Min.Y, rect.Min.X+int(border.Width/2), rect.Max.Y, shadowColor)
	// }
}

func drawOutsetBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// shadowColor := ic.RGBA{border.Color.R * 2 / 3, border.Color.G * 2 / 3, border.Color.B * 2 / 3, border.Color.A}
	// highlightColor := ic.RGBA{border.Color.R / 2, border.Color.G / 2, border.Color.B / 2, border.Color.A}
	// drawLine := func(x1, y1, x2, y2 int, col ic.RGBA) {
	// 	for i := 0; i < int(border.Width/2); i++ {
	// 		drawLineHelper(img, x1, y1+i, x2, y2+i, col)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, highlightColor)
	// 	drawLine(rect.Min.X, rect.Min.Y+int(border.Width/2), rect.Max.X, rect.Min.Y+int(border.Width/2), shadowColor)
	// case "right":
	// 	drawLine(rect.Max.X-int(border.Width/2), rect.Min.Y, rect.Max.X-int(border.Width/2), rect.Max.Y, highlightColor)
	// 	drawLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y, shadowColor)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width/2), rect.Max.X, rect.Max.Y-int(border.Width/2), shadowColor)
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width), highlightColor)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, shadowColor)
	// 	drawLine(rect.Min.X+int(border.Width/2), rect.Min.Y, rect.Min.X+int(border.Width/2), rect.Max.Y, highlightColor)
	// }
}

// FindPointOnLine calculates the coordinates of a point that is at a specified distance
// from the second point (P2) along the line defined by two points (P1 and P2).
func FindPointOnLine(P1, P2 image.Point, distance float64) image.Point {
	// Calculate the difference in x and y between P1 and P2
	dx := float64(P2.X - P1.X)
	dy := float64(P2.Y - P1.Y)

	// Calculate the length of the line segment between P1 and P2
	lineLength := math.Sqrt(dx*dx + dy*dy)

	// Normalize the direction vector (dx, dy) to unit length
	ux := dx / lineLength
	uy := dy / lineLength

	// Calculate the coordinates of the new point at the given distance from P2
	newX := float64(P2.X) + distance*ux
	newY := float64(P2.Y) + distance*uy

	// Return the new point as an image.Point (rounded to the nearest integer)
	return image.Point{
		X: int(math.Round(newX)),
		Y: int(math.Round(newY)),
	}
}

// Distance calculates the Euclidean distance between two points.
func Distance(p1, p2 image.Point) float64 {
	return math.Sqrt(float64((p2.X-p1.X)*(p2.X-p1.X) + (p2.Y-p1.Y)*(p2.Y-p1.Y)))
}

func FindClosestPoint(points []image.Point, target image.Point) (closestPoint image.Point, closestIndex int, minDistance float64) {
	if len(points) == 0 {
		// Return a zero-value point, -1 index, and infinity if the slice is empty
		return image.Point{}, -1, math.Inf(1)
	}

	closestPoint = points[0]
	closestIndex = 0
	minDistance = Distance(points[0], target)

	for i, point := range points[1:] {
		distance := Distance(point, target)
		if distance < minDistance {
			minDistance = distance
			closestPoint = point
			closestIndex = i + 1 // Adjust for the slice offset by 1
		}
	}

	return closestPoint, closestIndex, minDistance
}

// FindPointGivenAngleDistance calculates the coordinates of a point given the starting point,
// an angle in radians, and a distance from the starting point.
func FindPointGivenAngleDistance(start image.Point, angle float64, distance float64) image.Point {
	// Calculate the change in x and y based on the angle and distance
	deltaX := distance * math.Cos(angle)
	deltaY := distance * math.Sin(angle)

	// Calculate the new point coordinates
	newX := float64(start.X) + deltaX
	newY := float64(start.Y) + deltaY

	// Return the new point as an image.Point (rounded to the nearest integer)
	return image.Point{
		X: int(math.Round(newX)),
		Y: int(math.Round(newY)),
	}
}

func FindQuadraticControlPoint(P0, P2, P3 image.Point, t1 float64) image.Point {
	// Calculate (1-t1)^2, 2(1-t1)t1, and t1^2
	oneMinusT1 := 1 - t1
	oneMinusT1Squared := oneMinusT1 * oneMinusT1
	t1Squared := t1 * t1

	// Calculate the control point P1
	P1 := image.Point{
		X: int((float64(P2.X) - oneMinusT1Squared*float64(P0.X) - t1Squared*float64(P3.X)) / (2 * oneMinusT1 * t1)),
		Y: int((float64(P2.Y) - oneMinusT1Squared*float64(P0.Y) - t1Squared*float64(P3.Y)) / (2 * oneMinusT1 * t1)),
	}

	return P1
}

func FindCubicControlPoints(P0, P3, PA, PB image.Point, tA, tB float64) (P1, P2 image.Point) {
	// Coefficients for the equations
	A1 := 3 * tA * (1 - tA) * (1 - tA)
	A2 := 3 * tB * (1 - tB) * (1 - tB)
	B1 := 3 * (1 - tA) * tA * tA
	B2 := 3 * (1 - tB) * tB * tB

	// Solve for P1 and P2 using the given equations for PA and PB
	P1X := (float64(PA.X) - (1-tA)*(1-tA)*(1-tA)*float64(P0.X) - tA*tA*tA*float64(P3.X) -
		B1/B2*(float64(PB.X)-(1-tB)*(1-tB)*(1-tB)*float64(P0.X)-tB*tB*tB*float64(P3.X))) /
		(A1 - B1*B1/B2)
	P1Y := (float64(PA.Y) - (1-tA)*(1-tA)*(1-tA)*float64(P0.Y) - tA*tA*tA*float64(P3.Y) -
		B1/B2*(float64(PB.Y)-(1-tB)*(1-tB)*(1-tB)*float64(P0.Y)-tB*tB*tB*float64(P3.Y))) /
		(A1 - B1*B1/B2)

	P2X := (float64(PB.X) - (1-tB)*(1-tB)*(1-tB)*float64(P0.X) - tB*tB*tB*float64(P3.X) - A2*P1X) / B2
	P2Y := (float64(PB.Y) - (1-tB)*(1-tB)*(1-tB)*float64(P0.Y) - tB*tB*tB*float64(P3.Y) - A2*P1Y) / B2

	P1 = image.Point{X: int(P1X), Y: int(P1Y)}
	P2 = image.Point{X: int(P2X), Y: int(P2Y)}

	return P1, P2
}

func FindBorderStopAngle(origin, crossPoint, circleCenter image.Point, radius float64) []float64 {
	// Calculate the difference between the points
	dx := float64(origin.X - crossPoint.X)
	dy := float64(origin.Y - crossPoint.Y)

	// Calculate the angle using Atan2
	angle := math.Atan2(dy, dx)

	points := LineCircleIntersection(origin, angle, circleCenter, radius)

	if len(points) < 2 {
		return []float64{0, 0}
	}

	// Convert the angle from radians to degrees if needed

	dx = float64(circleCenter.X - points[0].X)
	dy = float64(circleCenter.Y - points[0].Y)

	angle2 := math.Atan2(dy, dx)

	dx = float64(circleCenter.X - points[1].X)
	dy = float64(circleCenter.Y - points[1].Y)

	angle3 := math.Atan2(dy, dx)
	return []float64{angle2, angle3}
}

func LineCircleIntersection(lineStart image.Point, angle float64, circleCenter image.Point, radius float64) []image.Point {
	// Parametric equations for the line: x = x0 + t * cos(theta), y = y0 + t * sin(theta)
	// Substitute these into the circle equation: (x - h)^2 + (y - k)^2 = r^2

	cosTheta := math.Cos(angle)
	sinTheta := math.Sin(angle)

	// Convert image.Point to float64 for calculations
	x0, y0 := float64(lineStart.X), float64(lineStart.Y)
	h, k := float64(circleCenter.X), float64(circleCenter.Y)

	// Coefficients for the quadratic equation At^2 + Bt + C = 0
	A := cosTheta*cosTheta + sinTheta*sinTheta
	B := 2 * (cosTheta*(x0-h) + sinTheta*(y0-k))
	C := (x0-h)*(x0-h) + (y0-k)*(y0-k) - radius*radius

	// Discriminant
	discriminant := B*B - 4*A*C

	// No intersection
	if discriminant < 0 {
		return nil
	}

	// Calculate the two solutions for t
	t1 := (-B + math.Sqrt(discriminant)) / (2 * A)
	t2 := (-B - math.Sqrt(discriminant)) / (2 * A)

	// Calculate the intersection points
	intersectionPoints := []image.Point{
		{
			X: int(math.Round(x0 + t1*cosTheta)),
			Y: int(math.Round(y0 + t1*sinTheta)),
		},
		{
			X: int(math.Round(x0 + t2*cosTheta)),
			Y: int(math.Round(y0 + t2*sinTheta)),
		},
	}

	// If the discriminant is zero, the line is tangent to the circle, and there's only one intersection point.
	if discriminant == 0 {
		return []image.Point{intersectionPoints[0]}
	}

	return intersectionPoints
}
