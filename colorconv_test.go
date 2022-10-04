package colorconv

import (
	"image/color"
	"testing"
)

func delta(x, y uint8) uint8 {
	if x >= y {
		return x - y
	}
	return y - x
}

// TestValueOutOfRange tests that when inputs are out of range, an error occur
func TestValueOutOfRange(t *testing.T) {
	var err error

	const (
		negativeValue1 = -0.1
		negativeValue2 = -2.1
		negativeValue3 = -10000.3333
		largeValue1    = 30000.0
		largeValue2    = 1.00000001
		largeValue3    = 2.0
	)

	_, _, _, err = HSVToRGB(negativeValue1, negativeValue2, negativeValue3)
	if err == nil {
		t.Errorf("\nerror not occur, input: h: %f, s:%f, v:%f", negativeValue1, negativeValue2, negativeValue3)
	}
	_, _, _, err = HSLToRGB(negativeValue1, negativeValue2, negativeValue3)
	if err == nil {
		t.Errorf("\nerror not occur, input: h: %f, s:%f, l:%f", negativeValue1, negativeValue2, negativeValue3)
	}
	_, err = HSVToColor(negativeValue1, negativeValue2, negativeValue3)
	if err == nil {
		t.Errorf("\nerror not occur, input: h: %f, s:%f, v:%f", negativeValue1, negativeValue2, negativeValue3)
	}
	_, err = HSLToColor(negativeValue1, negativeValue2, negativeValue3)
	if err == nil {
		t.Errorf("\nerror not occur, input: h: %f, s:%f, l:%f", negativeValue1, negativeValue2, negativeValue3)
	}

	_, _, _, err = HSVToRGB(largeValue1, largeValue2, largeValue3)
	if err == nil {
		t.Errorf("\nerror not occur, input: h: %f, s:%f, v:%f", largeValue1, largeValue2, largeValue3)
	}
	_, _, _, err = HSLToRGB(largeValue1, largeValue2, largeValue3)
	if err == nil {
		t.Errorf("\nerror not occur, input: h: %f, s:%f, l:%f", largeValue1, largeValue2, largeValue3)
	}
	_, err = HSVToColor(largeValue1, largeValue2, largeValue3)
	if err == nil {
		t.Errorf("\nerror not occur, input: h: %f, s:%f, v:%f", largeValue1, largeValue2, largeValue3)
	}
	_, err = HSLToColor(largeValue1, largeValue2, largeValue3)
	if err == nil {
		t.Errorf("\nerror not occur, input: h: %f, s:%f, l:%f", largeValue1, largeValue2, largeValue3)
	}
}

// TestInvalidHexInput tests that when inputs are invalid, an error occur
func TestInvalidHexInput(t *testing.T) {
	var err error
	tests := []struct {
		input string
	}{
		{input: "#GGGGGG"},
		{input: "0xgrewgrewg"},
		{input: "0xffxxff"},
		{input: "0xxxffff"},
		{input: "0xffffxx"},
	}

	for _, tc := range tests {
		_, _, _, err = HexToRGB(tc.input)
		if err == nil {
			t.Errorf("\nerror not occur, input: %s", tc.input)
		}
	}

	for _, tc := range tests {
		_, err = HexToColor(tc.input)
		if err == nil {
			t.Errorf("\nerror not occur, input: %s", tc.input)
		}
	}
}

// TestHSLRoundTrip tests that a subset of RGB space can be converted to HSL
// and back to within 1/256 tolerance.
func TestHSLRoundTrip(t *testing.T) {
	t.Run("HSL<->RGB", func(t *testing.T) {
		for r := 0; r < 256; r += 7 {
			for g := 0; g < 256; g += 5 {
				for b := 0; b < 256; b += 3 {
					r0, g0, b0 := uint8(r), uint8(g), uint8(b)
					h, s, l := RGBToHSL(r0, g0, b0)
					r1, g1, b1, err := HSLToRGB(h, s, l)
					if delta(r0, r1) > 1 || delta(g0, g1) > 1 || delta(b0, b1) > 1 || err != nil {
						t.Fatalf("\nr0, g0, b0 = %d, %d, %d\nh,  s, l = %f, %f, %f\nr1, g1, b1 = %d, %d, %d\nerr = %s",
							r0, g0, b0, h, s, l, r1, g1, b1, err)
					}
				}
			}
		}
	})
	t.Run("HSL<->Color", func(t *testing.T) {
		for r := 0; r < 256; r += 7 {
			for g := 0; g < 256; g += 5 {
				for b := 0; b < 256; b += 3 {
					r0, g0, b0 := uint8(r), uint8(g), uint8(b)
					h, s, l := ColorToHSL(color.RGBA{R: r0, G: g0, B: b0})
					c, err := HSLToColor(h, s, l)
					r1, g1, b1, _ := c.RGBA()
					if delta(r0, uint8(r1)) > 1 || delta(g0, uint8(g1)) > 1 || delta(b0, uint8(b1)) > 1 || err != nil {
						t.Fatalf("\nr0, g0, b0 = %d, %d, %d\nh,  s, l = %f, %f, %f\nr1, g1, b1 = %d, %d, %d\nerr = %s",
							r0, g0, b0, h, s, l, r1, g1, b1, err)
					}
				}
			}
		}
	})
}

// TestHSVRoundTrip tests that a subset of RGB space can be converted to HSV
// and back to within 1/256 tolerance.
func TestHSVRoundTrip(t *testing.T) {
	t.Run("HSV<->RGB", func(t *testing.T) {
		for r := 0; r < 256; r += 7 {
			for g := 0; g < 256; g += 5 {
				for b := 0; b < 256; b += 3 {
					r0, g0, b0 := uint8(r), uint8(g), uint8(b)
					h, s, v := RGBToHSV(r0, g0, b0)
					r1, g1, b1, err := HSVToRGB(h, s, v)
					if delta(r0, r1) > 1 || delta(g0, g1) > 1 || delta(b0, b1) > 1 || err != nil {
						t.Fatalf("\nr0, g0, b0 = %d, %d, %d\nh,  s, v = %f, %f, %f\nr1, g1, b1 = %d, %d, %d\nerr = %s",
							r0, g0, b0, h, s, v, r1, g1, b1, err)
					}
				}
			}
		}
	})
	t.Run("HSV<->Color", func(t *testing.T) {
		for r := 0; r < 256; r += 7 {
			for g := 0; g < 256; g += 5 {
				for b := 0; b < 256; b += 3 {
					r0, g0, b0 := uint8(r), uint8(g), uint8(b)
					h, s, v := ColorToHSV(color.RGBA{R: r0, G: g0, B: b0})
					c, err := HSVToColor(h, s, v)
					r1, g1, b1, _ := c.RGBA()
					if delta(r0, uint8(r1)) > 1 || delta(g0, uint8(g1)) > 1 || delta(b0, uint8(b1)) > 1 || err != nil {
						t.Fatalf("\nr0, g0, b0 = %d, %d, %d\nh,  s, v = %f, %f, %f\nr1, g1, b1 = %d, %d, %d\nerr = %s",
							r0, g0, b0, h, s, v, r1, g1, b1, err)
					}
				}
			}
		}
	})
}

// TestHexRoundTrip tests that a subset of RGB space can be converted to Hex and back
func TestHexRoundTrip(t *testing.T) {
	t.Run("Hex<->RGB", func(t *testing.T) {
		for r := 0; r < 256; r += 7 {
			for g := 0; g < 256; g += 5 {
				for b := 0; b < 256; b += 3 {
					r0, g0, b0 := uint8(r), uint8(g), uint8(b)
					hex := RGBToHex(r0, g0, b0)
					r1, g1, b1, err := HexToRGB(hex)
					if delta(r0, r1) != 0 || delta(g0, g1) != 0 || delta(b0, b1) != 0 || err != nil {
						t.Fatalf("\nr0, g0, b0 = %d, %d, %d\nhex = %s\nr1, g1, b1 = %d, %d, %d",
							r0, g0, b0, hex, r1, g1, b1)
					}
				}
			}
		}
	})

	t.Run("Hex<->Color", func(t *testing.T) {
		for r := 0; r < 256; r += 7 {
			for g := 0; g < 256; g += 5 {
				for b := 0; b < 256; b += 3 {
					r0, g0, b0 := uint8(r), uint8(g), uint8(b)
					hex := ColorToHex(color.RGBA{R: r0, G: g0, B: b0})
					c, err := HexToColor(hex)
					r1, g1, b1, _ := c.RGBA()
					if delta(r0, uint8(r1)) > 1 || delta(g0, uint8(g1)) > 1 || delta(b0, uint8(b1)) > 1 || err != nil {
						t.Fatalf("\nr0, g0, b0 = %d, %d, %d\nhex = %s\nr1, g1, b1 = %d, %d, %d",
							r0, g0, b0, hex, r1, g1, b1)
					}
				}
			}
		}
	})
}

// TestRGBGray compare RGBToGrayWithWeight result with RGB to GrayModel conversion in standard library.
func TestRGBGray(t *testing.T) {
	t.Run("RGB to Gray", func(t *testing.T) {
		for r := 0; r < 256; r += 7 {
			for g := 0; g < 256; g += 5 {
				for b := 0; b < 256; b += 3 {
					r0, g0, b0 := uint8(r), uint8(g), uint8(b)
					c := color.NRGBA{
						R: r0,
						G: g0,
						B: b0,
						A: 255,
					}
					gray := color.GrayModel.Convert(c)
					r1, g1, b1, _ := RGBToGrayWithWeight(r0, g0, b0, 299, 587, 114).RGBA()
					r2, g2, b2, _ := gray.RGBA()
					if delta(uint8(r1>>8), uint8(r2>>8)) > 1 || delta(uint8(g1>>8), uint8(g2>>8)) > 1 || delta(uint8(b1>>8), uint8(b2>>8)) > 1 {
						t.Fatalf("\nTest value:\tr0, g0, b0 = %d, %d, %d\ncolorconv value:\tr1, g1, b1 = %d, %d, %d\nGrayModel value:\tr2, g2, b2 = %d, %d, %d",
							r0, g0, b0, uint8(r1>>8), uint8(g1>>8), uint8(b1>>8), uint8(r2>>8), uint8(g2>>8), uint8(b2>>8))
					}
				}
			}
		}
	})
}

// use package level variable instead of ignore return values to avoid compiler optimization
// https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
var sinkUint8 uint8
var sinkFloat64 float64
var sinkGray color.Gray

func BenchmarkHSLToRGB(b *testing.B) {
	// Not really sure how to effectively benchmark these yet
	// so I follow BenchmarkYCbCrToRGB test (Low, Medium, High)
	b.Run("Low", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkUint8, sinkUint8, sinkUint8, _ = HSLToRGB(0, 0, 0)
		}
	})
	b.Run("Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkUint8, sinkUint8, sinkUint8, _ = HSLToRGB(180, 0.5, 0.5)
		}
	})
	b.Run("High", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkUint8, sinkUint8, sinkUint8, _ = HSLToRGB(360, 1, 1)
		}
	})
}

func BenchmarkRGBToHSL(b *testing.B) {
	// Not really sure how to effectively benchmark these yet
	// so I follow BenchmarkRGBToYCbCr test (Low, Medium, High)
	b.Run("0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkFloat64, sinkFloat64, sinkFloat64 = RGBToHSL(0, 0, 0)
		}
	})
	b.Run("128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkFloat64, sinkFloat64, sinkFloat64 = RGBToHSL(128, 128, 128)
		}
	})
	b.Run("255", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkFloat64, sinkFloat64, sinkFloat64 = RGBToHSL(255, 255, 255)
		}
	})
}

func BenchmarkHSVToRGB(b *testing.B) {
	// Not really sure how to effectively benchmark these yet
	// so I follow BenchmarkYCbCrToRGB test (Low, Medium, High)
	b.Run("Low", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkUint8, sinkUint8, sinkUint8, _ = HSVToRGB(0, 0, 0)
		}
	})
	b.Run("Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkUint8, sinkUint8, sinkUint8, _ = HSVToRGB(180, 0.5, 0.5)
		}
	})
	b.Run("High", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkUint8, sinkUint8, sinkUint8, _ = HSVToRGB(360, 1, 1)
		}
	})
}

func BenchmarkRGBToHSV(b *testing.B) {
	// Not really sure how to effectively benchmark these yet
	// so I follow BenchmarkRGBToYCbCr test (Low, Medium, High)
	b.Run("0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkFloat64, sinkFloat64, sinkFloat64 = RGBToHSV(0, 0, 0)
		}
	})
	b.Run("128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkFloat64, sinkFloat64, sinkFloat64 = RGBToHSV(128, 128, 128)
		}
	})
	b.Run("255", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkFloat64, sinkFloat64, sinkFloat64 = RGBToHSV(255, 255, 255)
		}
	})
}

func BenchmarkRGBToGrayAverage(b *testing.B) {
	// Not really sure how to effectively benchmark these yet
	// so I follow BenchmarkRGBToYCbCr test (Low, Medium, High)
	b.Run("0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkGray = RGBToGrayAverage(0, 0, 0)
		}
	})
	b.Run("128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkGray = RGBToGrayAverage(128, 128, 128)
		}
	})
	b.Run("255", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkGray = RGBToGrayAverage(255, 255, 255)
		}
	})
}
