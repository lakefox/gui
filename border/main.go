package border

import (
	"gui/canvas"
	"gui/color"
	"gui/element"
	"gui/library"
	"gui/utils"
	"image"
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
			ctx.SetStrokeStyle(0, 0, 0, 255)
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
			n.Textures = append(n.Textures, shelf.Set(key, ctx.RGBA))
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
		// drawDottedBorder(ctx, side, border, s)
	case "double":
		// drawDoubleBorder(ctx, side, border, s)
	case "groove":
		// drawGrooveBorder(ctx, side, border, s)
	case "ridge":
		// drawRidgeBorder(ctx, side, border, s)
	case "inset":
		// drawInsetBorder(ctx, side, border, s)
	case "outset":
		// drawOutsetBorder(ctx, side, border, s)
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

func drawSolidBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	ctx.SetFillStyle(border.Color.R, border.Color.G, border.Color.B, border.Color.A)

	width := float64(s.Width + s.Border.Left.Width + s.Border.Right.Width)
	height := float64(s.Height + s.Border.Top.Width + s.Border.Bottom.Width)

	ctx.BeginPath()
	ctx.Save()

	switch side {
	case "top":
		v1 := math.Max(float64(s.Border.Radius.TopLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.TopRight), 1)
		genSolidBorder(ctx, width, v1, v2, border, s.Border.Left, s.Border.Right)

	case "right":
		v1 := math.Max(float64(s.Border.Radius.TopRight), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomRight), 1)
		ctx.Translate(width, 0)
		ctx.Rotate(math.Pi / 2)
		genSolidBorder(ctx, height, v1, v2, border, s.Border.Top, s.Border.Bottom)

	case "bottom":
		v1 := math.Max(float64(s.Border.Radius.BottomLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomRight), 1)
		ctx.Translate(float64(width), float64(height))
		ctx.Rotate(math.Pi)
		genSolidBorder(ctx, width, v2, v1, border, s.Border.Right, s.Border.Left)

	case "left":
		v1 := math.Max(float64(s.Border.Radius.TopLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomLeft), 1)
		ctx.Translate(0, float64(height))
		ctx.Rotate(-math.Pi / 2)
		genSolidBorder(ctx, height, v2, v1, border, s.Border.Bottom, s.Border.Top)
	}

	if border.Width == 1 {
		ctx.Stroke()
	} else {
		ctx.Fill()
	}
	ctx.Reset()
	ctx.ClosePath()
}

func drawDashedBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	ctx.SetFillStyle(border.Color.R, border.Color.G, border.Color.B, border.Color.A)

	width := float64(s.Width + s.Border.Left.Width + s.Border.Right.Width)
	height := float64(s.Height + s.Border.Top.Width + s.Border.Bottom.Width)

	ctx.BeginPath()
	ctx.Save()

	switch side {
	case "top":
		v1 := math.Max(float64(s.Border.Radius.TopLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.TopRight), 1)
		genDashedBorder(ctx, width, v1, v2, border, s.Border.Left, s.Border.Right)

	case "right":
		v1 := math.Max(float64(s.Border.Radius.TopRight), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomRight), 1)
		ctx.Translate(width, 0)
		ctx.Rotate(math.Pi / 2)
		genDashedBorder(ctx, height, v1, v2, border, s.Border.Top, s.Border.Bottom)

	case "bottom":
		v1 := math.Max(float64(s.Border.Radius.BottomLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomRight), 1)
		ctx.Translate(float64(width), float64(height))
		ctx.Rotate(math.Pi)
		genDashedBorder(ctx, width, v2, v1, border, s.Border.Right, s.Border.Left)

	case "left":
		v1 := math.Max(float64(s.Border.Radius.TopLeft), 1)
		v2 := math.Max(float64(s.Border.Radius.BottomLeft), 1)
		ctx.Translate(0, float64(height))
		ctx.Rotate(-math.Pi / 2)
		genDashedBorder(ctx, height, v2, v1, border, s.Border.Bottom, s.Border.Top)
	}

	// ctx.Stroke()
	ctx.Reset()
	ctx.ClosePath()
}

func genSolidBorder(ctx *canvas.Canvas, width float64, v1, v2 float64, border, side1, side2 element.BorderSide) {
	s1w := math.Max(float64(side1.Width), 1)
	s2w := math.Max(float64(side2.Width), 1)

	ctx.SetLineWidth(1)
	splX, splY, sprX, sprY := genTopLine(ctx, s1w, s2w, float64(border.Width), v1, v2, width, 0, 0, 0)

	if border.Width == 1 {
		return
	}

	// Right flat
	min := float64(border.Width)
	if min > s2w {
		min = s2w
	}

	d := (math.Abs(float64(border.Width)-s2w) / 2) + min

	xe, ye := EndPointFromMidpoint(width-(s2w/10), float64(border.Width/10), sprX, sprY, d, true)

	ctx.LineTo(xe, ye)

	// Ellipse right

	xr, yr := CalculateInnerRadii(v2, s2w, float64(border.Width))

	sa := AngleInRadians(width-(s2w+xr), float64(border.Width)+yr, xe, ye)

	ctx.Ellipse(width-(s2w+xr), float64(border.Width)+yr, xr, yr, 0, sa, -math.Pi/2, false)

	// Find top of left ellipse to draw line to

	xr, yr = CalculateInnerRadii(v1, s1w, float64(border.Width))

	ellipseEndX, ellipseEndY := EllipsePoint(s1w+xr, float64(border.Width)+yr, xr, yr, 0, -math.Pi/2)

	// Bottom Line
	ctx.LineTo(ellipseEndX, ellipseEndY)

	// Left flat
	min = float64(border.Width)
	if min > s1w {
		min = s1w
	}

	d = (math.Abs(float64(border.Width)-s1w) / 2) + min
	xe, ye = EndPointFromMidpoint((s1w / 10), float64(border.Width/10), splX, splY, d, true)

	sa = AngleInRadians(s1w+xr, float64(border.Width)+yr, xe, ye)

	// Ellipse left
	ctx.Ellipse(s1w+xr, float64(border.Width)+yr, xr, yr, 0, -math.Pi/2, sa, false)
	ctx.LineTo(xe, ye)

	// Left flat line
	ctx.ClosePath()

}

func genDashedBorder(ctx *canvas.Canvas, width float64, v1, v2 float64, border, side1, side2 element.BorderSide) {
	s1w := math.Max(float64(side1.Width), 1)
	s2w := math.Max(float64(side2.Width), 1)

	ctx.BeginPath()
	genSolidBorder(ctx, width, v1, v2, border, side1, side2)
	ctx.Clip()

	ctx.Context.SetLineCapSquare()
	// !ISSUE: Doesn't fill up the clipping region, also this is a CSS prop
	w := math.Max(float64(border.Width), math.Max(s1w, s2w))
	ctx.SetLineWidth(w)
	ctx.SetLineDash(50)
	genTopLine(ctx, s1w, s2w, float64(border.Width), v1, v2, width, s1w/2, float64(border.Width)/2, s2w/2)
	ctx.Stroke()
}

func genTopLine(ctx *canvas.Canvas, s1w, s2w, borderWidth, v1, v2, width, o1, o2, o3 float64) (float64, float64, float64, float64) {
	// Top-left corner arc
	startAngleLeft := FindBorderStopAngle(
		image.Point{X: 0, Y: 0},
		image.Point{X: int(s1w), Y: int(borderWidth)},
		image.Point{X: int(v1), Y: int(v1)},
		v1-o1,
	)
	ctx.Arc(v1, v1, v1-o1, startAngleLeft[0]-math.Pi, -math.Pi/2)
	// Reversed to get startpoint
	splX, splY := PointAtAngle(v1, v1, v1-o1, startAngleLeft[0]-math.Pi)
	// top line
	ctx.LineTo(width-v2-o3, o3)

	// Top-right corner arc
	startAngleRight := FindBorderStopAngle(
		image.Point{X: int(width), Y: 0},
		image.Point{X: int(width - s2w), Y: int(borderWidth)},
		image.Point{X: int(width - v2), Y: int(v2)},
		v2-o3,
	)
	ctx.Arc(width-v2, v2, v2-o3, -math.Pi/2, startAngleRight[0]-math.Pi)
	sprX, sprY := PointAtAngle(width-v2, v2, v2-o3, startAngleRight[0]-math.Pi)
	return splX, splY, sprX, sprY
}

func FindBorderStopAngle(origin, crossPoint, circleCenter image.Point, radius float64) []float64 {
	dx := float64(origin.X - crossPoint.X)
	dy := float64(origin.Y - crossPoint.Y)

	angle := math.Atan2(dy, dx)
	points := LineCircleIntersection(origin, angle, circleCenter, radius)

	// Validate points
	if len(points) < 2 {
		return []float64{0, 0}
	}

	var result []float64
	for _, p := range points {
		// Check if the point lies within the arc's visible range
		dx = float64(circleCenter.X - p.X)
		dy = float64(circleCenter.Y - p.Y)
		angle = math.Atan2(dy, dx)
		result = append(result, angle)
	}

	return result
}

func LineCircleIntersection(lineStart image.Point, angle float64, circleCenter image.Point, radius float64) []image.Point {
	cosTheta := math.Cos(angle)
	sinTheta := math.Sin(angle)

	x0, y0 := float64(lineStart.X), float64(lineStart.Y)
	h, k := float64(circleCenter.X), float64(circleCenter.Y)

	A := cosTheta*cosTheta + sinTheta*sinTheta
	B := 2 * (cosTheta*(x0-h) + sinTheta*(y0-k))
	C := (x0-h)*(x0-h) + (y0-k)*(y0-k) - radius*radius

	// Use a small epsilon to handle near-tangent cases
	const epsilon = 1e-10
	discriminant := B*B - 4*A*C

	if discriminant < -epsilon {
		return nil // No intersection
	}

	if math.Abs(discriminant) < epsilon {
		// Tangent case: return one intersection point
		t := -B / (2 * A)
		return []image.Point{
			{
				X: int(math.Round(x0 + t*cosTheta)),
				Y: int(math.Round(y0 + t*sinTheta)),
			},
		}
	}

	// Two intersection points
	t1 := (-B + math.Sqrt(discriminant)) / (2 * A)
	t2 := (-B - math.Sqrt(discriminant)) / (2 * A)

	return []image.Point{
		{
			X: int(math.Round(x0 + t1*cosTheta)),
			Y: int(math.Round(y0 + t1*sinTheta)),
		},
		{
			X: int(math.Round(x0 + t2*cosTheta)),
			Y: int(math.Round(y0 + t2*sinTheta)),
		},
	}
}

func EndPointFromMidpoint(x1, y1, xm, ym, distance float64, reverse bool) (xe, ye float64) {
	// Compute direction vector from midpoint to the reference point
	dx := (x1 - xm)
	dy := (y1 - ym)

	// Handle the case where midpoint and reference point are the same
	if dx == 0 && dy == 0 {
		// Return a point directly along the x-axis for simplicity
		if reverse {
			return xm - distance, ym
		}
		return xm + distance, ym
	}

	// Normalize the direction vector
	length := math.Sqrt(dx*dx + dy*dy)
	ux := dx / length
	uy := dy / length

	// Reverse the direction if needed
	if reverse {
		ux = -ux
		uy = -uy
	}

	// Calculate the endpoint
	xe = xm + ux*distance
	ye = ym + uy*distance
	return
}

// CalculateInnerRadii calculates the inner radii of a box corner given the outer radii and border widths.
func CalculateInnerRadii(outerRadius, borderWidthLeft, borderWidthTop float64) (float64, float64) {
	innerRadiusX := outerRadius - borderWidthLeft
	innerRadiusY := outerRadius - borderWidthTop

	// Clamp the inner radii to zero if they go negative
	if innerRadiusX < 0 {
		innerRadiusX = 0
	}
	if innerRadiusY < 0 {
		innerRadiusY = 0
	}

	return innerRadiusX, innerRadiusY
}

// AngleInRadians calculates the angle from one point to another in radians
func AngleInRadians(p1X, p1Y, p2X, p2Y float64) float64 {
	// Calculate the differences in x and y
	dx := p2X - p1X
	dy := p2Y - p1Y

	// Use math.Atan2 to calculate the angle in radians
	return math.Atan2(dy, dx)
}

func EllipsePoint(x, y, radiusX, radiusY, rotation, angle float64) (float64, float64) {
	// Calculate the point on the ellipse without rotation
	ellipseX := radiusX * math.Cos(angle)
	ellipseY := radiusY * math.Sin(angle)

	// Apply rotation to the point
	rotatedX := ellipseX*math.Cos(rotation) - ellipseY*math.Sin(rotation)
	rotatedY := ellipseX*math.Sin(rotation) + ellipseY*math.Cos(rotation)

	// Translate the point to the center of the ellipse
	finalX := rotatedX + x
	finalY := rotatedY + y

	return finalX, finalY
}

func PointAtAngle(x, y, radius, endAngle float64) (float64, float64) {
	// Calculate the endpoint coordinates
	endX := x + radius*math.Cos(endAngle)
	endY := y + radius*math.Sin(endAngle)

	return endX, endY
}
