// Package colorconv provide conversion of color to HSL, HSV and hex value.
// All the conversion methods is based on the website: https://www.rapidtables.com/convert/color/index.html
package colorconv

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

//ColorToHSL convert Color into HSL triple, ignoring the alpha channel.
func ColorToHSL(c color.Color) (h, s, l float64) {
	r, g, b, _ := c.RGBA()
	return RGBToHSL(uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

//ColorToHSV convert Color into HSV triple, ignoring the alpha channel.
func ColorToHSV(c color.Color) (h, s, v float64) {
	r, g, b, _ := c.RGBA()
	return RGBToHSV(uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

//ColorToHex convert Color into Hex string, ignoring the alpha channel.
func ColorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return RGBToHex(uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

//RGBToHSL converts a RGB triple to an HSL triple.
func RGBToHSL(r, g, b uint8) (h, s, l float64) {
	// TODO add in value-out-of-range error
	// convert uint32 pre-multiplied value to uint8
	// The r,g,b values are divided by 255 to change the range from 0..255 to 0..1:
	Rnot := float64(r) / 255
	Gnot := float64(g) / 255
	Bnot := float64(b) / 255
	Cmax, Cmin := getMaxMin(Rnot, Gnot, Bnot)
	delta := Cmax - Cmin
	// Lightness calculation:
	l = (Cmax + Cmin) / 2
	// Hue and Saturation Calculation:
	if delta == 0 {
		h = 0
		s = 0
	} else {
		switch Cmax {
		case Rnot:
			h = 60 * (math.Mod((Gnot-Bnot)/delta, 6))
		case Gnot:
			h = 60 * (((Bnot - Rnot) / delta) + 2)
		case Bnot:
			h = 60 * (((Rnot - Gnot) / delta) + 4)
		}
		if h < 0 {
			h += 360
		}

		s = delta / (1 - math.Abs((2*l)-1))
	}

	return h, s, l
}

func getMaxMin(a, b, c float64) (max, min float64) {
	if a > b {
		max = a
		min = b
	} else {
		max = b
		min = a
	}
	if c > max {
		max = c
	} else if c < min {
		min = c
	}
	return max, min
}

//HSLToRGB converts a HSL triple to an RGB triple.
func HSLToRGB(h, s, l float64) (r, g, b uint8) {
	// TODO add in value-out-of-range error
	// When 0 ≤ h < 360, 0 ≤ s ≤ 1 and 0 ≤ l ≤ 1:
	C := (1 - math.Abs((2*l)-1)) * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - (C / 2)
	var Rnot, Gnot, Bnot float64
	if h >= 360 {
		h -= 360
	}
	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r = uint8(math.Round((Rnot + m) * 255))
	g = uint8(math.Round((Gnot + m) * 255))
	b = uint8(math.Round((Bnot + m) * 255))
	return r, g, b
}

//RGBToHSV converts a RGB triple to an HSV triple.
func RGBToHSV(r, g, b uint8) (h, s, v float64) {
	// TODO add in value-out-of-range error
	// convert uint32 pre-multiplied value to uint8
	// The r,g,b values are divided by 255 to change the range from 0..255 to 0..1:
	Rnot := float64(r) / 255
	Gnot := float64(g) / 255
	Bnot := float64(b) / 255
	Cmax, Cmin := getMaxMin(Rnot, Gnot, Bnot)
	delta := Cmax - Cmin

	// Hue calculation:
	if delta == 0 {
		h = 0
	} else {
		switch Cmax {
		case Rnot:
			h = 60 * (math.Mod((Gnot-Bnot)/delta, 6))
		case Gnot:
			h = 60 * (((Bnot - Rnot) / delta) + 2)
		case Bnot:
			h = 60 * (((Rnot - Gnot) / delta) + 4)
		}
		if h < 0 {
			h += 360
		}

	}
	// Saturation calculation:
	if Cmax == 0 {
		s = 0
	} else {
		s = delta / Cmax
	}
	// Value calculation:
	v = Cmax

	return h, s, v
}

//HSVToRGB converts a HSV triple to an RGB triple.
func HSVToRGB(h, s, v float64) (r, g, b uint8) {
	// TODO add in value-out-of-range error
	// When 0 ≤ h < 360, 0 ≤ s ≤ 1 and 0 ≤ v ≤ 1:
	C := v * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - C
	var Rnot, Gnot, Bnot float64
	if h >= 360 {
		h -= 360
	}
	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r = uint8(math.Round((Rnot + m) * 255))
	g = uint8(math.Round((Gnot + m) * 255))
	b = uint8(math.Round((Bnot + m) * 255))
	return r, g, b
}

//RGBToHex converts a RGB triple to an Hex string in the format of 0xffff.
func RGBToHex(r, g, b uint8) string {
	// TODO add in value-out-of-range error
	return fmt.Sprintf("0x%02x%02x%02x", r, g, b)
}

//HexToRGB converts a Hex string to an RGB triple.
func HexToRGB(hex string) (r, g, b uint8) {
	// remove prefixes if found in the input string
	hex = strings.Replace(hex, "0x", "", -1)
	hex = strings.Replace(hex, "#", "", -1)
	//TODO check range
	if len(hex) != 6 {
		panic("not a valid input")
	}

	hex2uint8 := func(hexStr string) uint8 {
		// base 16 for hexadecimal
		result, err := strconv.ParseUint(hexStr, 16, 8)
		if err != nil {
			panic(err)
		}
		return uint8(result)
	}
	r = hex2uint8(hex[0:2])
	g = hex2uint8(hex[2:4])
	b = hex2uint8(hex[4:6])
	// TODO add in value-out-of-range error
	return r, g, b
}
