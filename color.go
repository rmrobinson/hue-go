package hue

import (
	"math"
)

const (
	colorPointRed   = 0
	colorPointGreen = 1
	colorPointBlue  = 2
)

// RGB is a colour represented using the red/green/blue colour model.
type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// XY is a colour represented using the CIE colour space.
type XY struct {
	X float64
	Y float64
}

// HSB is a colour represented using the hue/saturation/value representation of the RGB colour model.
type HSB struct {
	Hue        uint16
	Saturation uint8
	Brightness uint8
}

// colourPointsForModel returns the XY bounds for the specified lightbulb model.
// The returned array always has the red, then green, then blue points in that order.
func colorPointsForModel(model string) (points []XY) {
	points = make([]XY, 3)
	switch model {
	case "LCT001", "LCT002", "LCT003":
		points[colorPointRed].X = 0.674
		points[colorPointRed].Y = 0.322
		points[colorPointGreen].X = 0.408
		points[colorPointGreen].Y = 0.517
		points[colorPointBlue].X = 0.168
		points[colorPointBlue].Y = 0.041
		return

	case "LLC001", "LLC005", "LLC006", "LLC007", "LLC011", "LLC012", "LLC013", "LST001":
		points[colorPointRed].X = 0.703
		points[colorPointRed].Y = 0.296
		points[colorPointGreen].X = 0.214
		points[colorPointGreen].Y = 0.709
		points[colorPointBlue].X = 0.139
		points[colorPointBlue].Y = 0.081
		return
	}

	points[colorPointRed].X = 1.0
	points[colorPointRed].Y = 0.0
	points[colorPointGreen].X = 0.0
	points[colorPointGreen].Y = 1.0
	points[colorPointBlue].X = 0.0
	points[colorPointBlue].Y = 0.0
	return
}

func crossProduct(p1, p2 XY) float64 {
	return p1.X*p2.Y - p1.Y*p2.X
}

func getClosestPointToPoints(a, b, p XY) XY {
	ap := XY{X: p.X - a.X, Y: p.Y - a.Y}
	ab := XY{X: b.X - a.X, Y: b.Y - a.Y}

	ab2 := ab.X*ab.X + ab.Y*ab.Y
	apAB := ap.X*ab.X + ap.Y*ab.Y

	t := apAB / ab2

	if t < 0.0 {
		t = 0.0
	} else if t > 1.0 {
		t = 1.0
	}

	return XY{X: a.X + ab.X*t, Y: a.Y + ab.Y*t}
}

func getDistanceBetweenTwoPoints(p1, p2 XY) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y

	return math.Sqrt(dx*dx + dy*dy)
}

func checkPointInColorPointsReach(p XY, colorPoints []XY) bool {
	if len(colorPoints) != 3 {
		return false
	}

	red := colorPoints[colorPointRed]
	green := colorPoints[colorPointGreen]
	blue := colorPoints[colorPointBlue]

	v1 := XY{X: green.X - red.X, Y: green.Y - red.Y}
	v2 := XY{X: blue.X - red.X, Y: blue.Y - red.Y}
	q := XY{X: p.X - red.X, Y: p.Y - red.Y}

	s := crossProduct(q, v2) / crossProduct(v1, v2)
	t := crossProduct(v1, q) / crossProduct(v1, v2)

	if s >= 0.0 && t >= 0.0 && s+t <= 1.0 {
		return true
	}
	return false
}

// FromHSB converts the specified HSB value into the RGB colour space.
// This algorithm is adapted from the code at http://www.docjar.com/html/api/java/awt/Color.java.html
func (c *RGB) FromHSB(from HSB) {
	hue := float64(from.Hue / 65535)
	saturation := float64(from.Saturation / 255)
	brightness := float64(from.Brightness / 255)

	if saturation == 0 {
		c.Red = uint8(brightness*float64(255.0) + float64(0.5))
		c.Green = uint8(brightness*float64(255.0) + float64(0.5))
		c.Blue = uint8(brightness*float64(255.0) + float64(0.5))

		return
	}

	h := (hue - math.Floor(hue)) * 6.0
	f := h - math.Floor(h)
	p := brightness * (1.0 - saturation)
	q := brightness * (1.0 - saturation*f)
	t := brightness * (1.0 - (saturation * (1.0 - f)))

	var red, green, blue float64

	switch h {
	case 0:
		red = brightness*float64(255.0) + float64(0.5)
		green = t*float64(255.0) + float64(0.5)
		blue = p*float64(255.0) + float64(0.5)
		break
	case 1:
		red = q*float64(255.0) + float64(0.5)
		green = brightness*float64(255.0) + float64(0.5)
		blue = p*float64(255.0) + float64(0.5)
		break
	case 2:
		red = p*float64(255.0) + float64(0.5)
		green = brightness*float64(255.0) + float64(0.5)
		blue = t*float64(255.0) + float64(0.5)
		break
	case 3:
		red = p*float64(255.0) + float64(0.5)
		green = q*float64(255.0) + float64(0.5)
		blue = brightness*float64(255.0) + float64(0.5)
		break
	case 4:
		red = t*float64(255.0) + float64(0.5)
		green = p*float64(255.0) + float64(0.5)
		blue = brightness*float64(255.0) + float64(0.5)
		break
	case 5:
		red = brightness*float64(255.0) + float64(0.5)
		green = p + float64(255.0) + float64(0.5)
		blue = q*float64(255.0) + float64(0.5)
	}

	c.Red = uint8(red)
	c.Green = uint8(green)
	c.Blue = uint8(blue)

	return
}

// FromCT converts the specified CT value into the RGB colour space.
// This algorithm is adapted from the example at http://www.tannerhelland.com/4435/convert-temperature-rgb-algorithm-code/
func (c *RGB) FromCT(from uint16) {
	var temp float64
	temp = 1000000 / float64(from)
	temp = temp / 100

	if temp < 66 {
		c.Red = 255
	} else {
		red := temp - 60
		red = 329.698727446 * math.Pow(red, -0.1332047592)

		if red < 0 {
			c.Red = 0
		} else if red > 255 {
			c.Red = 255
		} else {
			c.Red = uint8(red)
		}
	}

	if temp <= 66 {
		green := temp
		green = 99.4708025861*math.Log(green) - 161.1195681661

		if green < 0 {
			c.Green = 0
		} else if green > 255 {
			c.Green = 255
		} else {
			c.Green = uint8(green)
		}
	} else {
		green := temp
		green = 288.1221695283 * math.Pow(green, -0.0755148492)

		if green < 0 {
			c.Green = 0
		} else if green > 255 {
			c.Green = 255
		} else {
			c.Green = uint8(green)
		}
	}

	if temp >= 66 {
		c.Blue = 255
	} else {
		if temp <= 19 {
			c.Blue = 0
		} else {
			blue := temp - 10
			blue = 138.5177312231*math.Log(blue) - 305.0447927307

			if blue < 0 {
				c.Blue = 0
			} else if blue > 255 {
				c.Blue = 255
			} else {
				c.Blue = uint8(blue)
			}
		}
	}
}

// FromXY converts the specified XY value into the RGB colour space.
// The supplied light model is used to adjust the input value accordingly.
// This algorithm is adapted from the examples at http://www.developers.meethue.com/documentation/color-conversions-rgb-xy
func (c *RGB) FromXY(from XY, model string) {
	xy := XY{X: from.X, Y: from.Y}
	colorPoints := colorPointsForModel(model)
	isColorReachable := checkPointInColorPointsReach(xy, colorPoints)

	if !isColorReachable {
		// We will have to map the requested color to the closest representable color.

		pAB := getClosestPointToPoints(colorPoints[colorPointRed], colorPoints[colorPointGreen], xy)
		pAC := getClosestPointToPoints(colorPoints[colorPointBlue], colorPoints[colorPointRed], xy)
		pBC := getClosestPointToPoints(colorPoints[colorPointGreen], colorPoints[colorPointBlue], xy)

		dAB := getDistanceBetweenTwoPoints(xy, pAB)
		dAC := getDistanceBetweenTwoPoints(xy, pAC)
		dBC := getDistanceBetweenTwoPoints(xy, pBC)

		lowest := dAB
		closestPoint := pAB

		if dAC < lowest {
			lowest = dAC
			closestPoint = pAC
		}
		if dBC < lowest {
			lowest = dBC
			closestPoint = pBC
		}

		xy.X = closestPoint.X
		xy.Y = closestPoint.Y
	}

	x := xy.X
	y := xy.Y
	z := 1.0 - x - y

	Y := 1.0
	X := (Y / y) * x
	Z := (Y / y) * z

	// sRGB D65 conversion
	// Option 1
	/*
		r := X*1.656492 - Y*0.354851 - Z*0.255038
		g := -1*X*0.707196 + Y*1.655397 + Z*0.036152
		b := X*0.051713 - Y*0.121364 + Z*1.011530
	*/
	// Option 2
	r := X*1.4628067 - Y*0.1840623 - Z*0.2743606
	g := -X*0.5217933 + Y*1.4472381 + Z*0.0677227
	b := X*0.0349342 - Y*0.0968930 + Z*1.2884099

	// Check if any color is too large and scale it down accordingly
	if r > b && r > g && r > 1.0 {
		g = g / r
		b = b / r
		r = 1.0
	} else if g > b && g > r && g > 1.0 {
		r = r / g
		b = b / g
		g = 1.0
	} else if b > r && b > g && b > 1.0 {
		r = r / b
		g = g / b
		b = 1.0
	}

	// Apply gamma correction
	if r <= 0.0031308 {
		r = r * 12.92
	} else {
		r = (1.0+0.055)*math.Pow(r, (1.0/2.4)) - 0.055
	}

	if g <= 0.0031308 {
		g = g * 12.92
	} else {
		g = (1.0+0.055)*math.Pow(g, (1.0/2.4)) - 0.055
	}

	if b <= 0.0031308 {
		b = b * 12.92
	} else {
		b = (1.0+0.055)*math.Pow(b, (1.0/2.4)) - 0.055
	}

	// Check if any color is too large and scale it down accordingly
	if r > b && r > g {
		if r > 1.0 {
			g = g / r
			b = b / r
			r = 1.0
		}
	} else if g > b && g > r {
		if g > 1.0 {
			r = r / g
			b = b / g
			g = 1.0
		}
	} else if b > r && b > g {
		if b > 1.0 {
			r = r / b
			g = g / b
			b = 1.0
		}
	}

	c.Red = uint8(r * 255)
	c.Green = uint8(g * 255)
	c.Blue = uint8(b * 255)
	return
}

// FromRGB converts the specified RGB value into the CIE colour space.
// The supplied light model is used to adjust the input value accordingly.
// This algorithm is adapted from the examples at http://www.developers.meethue.com/documentation/color-conversions-rgb-xy
func (c *XY) FromRGB(from RGB, model string) {
	red := float64(from.Red / 255)
	green := float64(from.Green / 255)
	blue := float64(from.Blue / 255)

	r := 1.0
	g := 1.0
	b := 1.0

	// Gamma correction
	if red > 0.04045 {
		r = math.Pow((red+0.055)/(1.0+0.055), 2.4)
	} else {
		r = red / 12.92
	}

	if green > 0.04045 {
		g = math.Pow((green+0.055)/(1.0+0.055), 2.4)
	} else {
		g = green / 12.92
	}

	if blue > 0.0405 {
		b = math.Pow((blue+0.055)/(1.0+0.055), 2.4)
	} else {
		b = blue / 12.92
	}

	// Convert RGB to XYZ using Wide RGB D65 conversion
	// Option 1:
	/*
		X := r*0.664511 + g*0.154324 + b*0.162028
		Y := r*0.283881 + g*0.668433 + b*0.047685
		Z := r*0.000088 + g*0.072310 + b*0.986039
	*/
	// Option 2
	X := r*0.649926 + g*0.103455 + b*0.197109
	Y := r*0.234327 + g*0.743075 + b*0.022598
	Z := r*0.000000 + g*0.053077 + b*1.035763

	cx := X / (X + Y + Z)
	cy := Y / (X + Y + Z)

	if math.IsNaN(cx) {
		cx = 0.0
	}
	if math.IsNaN(cy) {
		cy = 0.0
	}

	// Check if the requested XY value is within the color range of the light.
	xy := XY{X: cx, Y: cy}
	colorPoints := colorPointsForModel(model)
	isColorReachable := checkPointInColorPointsReach(xy, colorPoints)

	if !isColorReachable {
		// Find the closest color we can reach and send this instead

		pAB := getClosestPointToPoints(colorPoints[colorPointRed], colorPoints[colorPointGreen], xy)
		pAC := getClosestPointToPoints(colorPoints[colorPointBlue], colorPoints[colorPointRed], xy)
		pBC := getClosestPointToPoints(colorPoints[colorPointGreen], colorPoints[colorPointBlue], xy)

		dAB := getDistanceBetweenTwoPoints(xy, pAB)
		dAC := getDistanceBetweenTwoPoints(xy, pAC)
		dBC := getDistanceBetweenTwoPoints(xy, pBC)

		lowest := dAB
		closestPoint := pAB

		if dAC < lowest {
			lowest = dAC
			closestPoint = pAC
		}
		if dBC < lowest {
			lowest = dBC
			closestPoint = pBC
		}

		cx = closestPoint.X
		cy = closestPoint.Y
	}

	c.X = cx
	c.Y = cy

	return
}
