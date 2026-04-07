// icongen generates the default Yapp application icon at
// internal/bundle/defaulticon.png. It is intentionally stdlib-only so
// the icon asset can be reproduced from source without pulling in any
// image-processing dependencies.
//
// Run from the repo root:
//
//	go run ./tools/icongen
//
// The output is committed to the repo and embedded into the binary via
// //go:embed so `yapp-cli install` works without any external asset.
package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const (
	size   = 1024
	output = "internal/bundle/defaulticon.png"

	// Corner radius for the rounded-rect background. Not a true macOS
	// squircle; macOS will still apply its own compositing in Launchpad
	// but the icon reads cleanly in Dock/Spotlight/Cmd+Tab as-is.
	cornerRadius = 180.0

	// Stroke thickness of the Y glyph in pixels.
	strokeThickness = 98.0

	// Width of the anti-aliasing falloff band at every colored edge.
	antialiasHalfWidth = 2.0
)

// Catppuccin Mocha palette. Dark background + warm yellow foreground reads
// unambiguously at 16 px (Spotlight) and stays distinct in the Dock.
var (
	colorBackdrop = color.RGBA{0x1E, 0x1E, 0x2E, 0xFF} // base
	colorSurface  = color.RGBA{0x31, 0x32, 0x44, 0xFF} // surface0 — rounded tile
	colorGlyph    = color.RGBA{0xF9, 0xE2, 0xAF, 0xFF} // yellow
)

// Y glyph geometry, in the 1024x1024 canvas.
//  - Top fork tips sit at y=230 with 220 px of horizontal padding.
//  - Arms converge at the center point (512, 560).
//  - The stem descends from the converge point to y=830.
var ySegments = [3][4]float64{
	{220, 230, 512, 560}, // left arm
	{804, 230, 512, 560}, // right arm
	{512, 560, 512, 830}, // stem
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for y := range size {
		for x := range size {
			img.Set(x, y, pixel(float64(x)+0.5, float64(y)+0.5))
		}
	}

	if err := os.MkdirAll("internal/bundle", 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

// pixel computes the final color of the pixel whose sample point is (px, py).
// We layer three anti-aliased masks: the outer backdrop, the rounded-rect
// surface tile, and the Y glyph on top.
func pixel(px, py float64) color.RGBA {
	c := colorBackdrop

	tileDist := signedRoundRectDist(px, py, 60, 60, size-60, size-60, cornerRadius)
	c = overlay(c, colorSurface, tileDist)

	minSeg := math.Inf(1)
	for _, s := range ySegments {
		if d := distSeg(px, py, s[0], s[1], s[2], s[3]); d < minSeg {
			minSeg = d
		}
	}
	glyphDist := minSeg - strokeThickness
	c = overlay(c, colorGlyph, glyphDist)

	return c
}

// overlay blends fg over bg based on a signed distance value where <= 0 means
// "inside the shape" (fully fg), and > antialiasHalfWidth*2 means "outside"
// (fully bg). The band between the two edges is linearly interpolated.
func overlay(bg, fg color.RGBA, signedDist float64) color.RGBA {
	switch {
	case signedDist <= -antialiasHalfWidth:
		return fg
	case signedDist >= antialiasHalfWidth:
		return bg
	default:
		t := (signedDist + antialiasHalfWidth) / (2 * antialiasHalfWidth)
		return lerp(fg, bg, t)
	}
}

// distSeg returns the Euclidean distance from point (px, py) to the line
// segment (ax, ay)-(bx, by). Standard projection-clamped formulation.
func distSeg(px, py, ax, ay, bx, by float64) float64 {
	abx, aby := bx-ax, by-ay
	apx, apy := px-ax, py-ay
	lenSq := abx*abx + aby*aby
	t := 0.0
	if lenSq > 0 {
		t = (apx*abx + apy*aby) / lenSq
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
	}
	dx := px - (ax + t*abx)
	dy := py - (ay + t*aby)
	return math.Sqrt(dx*dx + dy*dy)
}

// signedRoundRectDist returns the signed distance from (px, py) to the
// rounded rectangle (x0, y0)-(x1, y1) with the given corner radius.
// Negative inside, positive outside, continuous across the boundary.
func signedRoundRectDist(px, py, x0, y0, x1, y1, r float64) float64 {
	// Distance to the axis-aligned shrunken rectangle (shrunk by r on all
	// sides) then subtract r. This is the canonical SDF trick for a
	// rounded-rect.
	cx := (x0 + x1) / 2
	cy := (y0 + y1) / 2
	halfW := (x1-x0)/2 - r
	halfH := (y1-y0)/2 - r
	dx := math.Abs(px-cx) - halfW
	dy := math.Abs(py-cy) - halfH
	outside := math.Sqrt(math.Max(dx, 0)*math.Max(dx, 0) + math.Max(dy, 0)*math.Max(dy, 0))
	inside := math.Min(math.Max(dx, dy), 0)
	return outside + inside - r
}

func lerp(a, b color.RGBA, t float64) color.RGBA {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	mix := func(x, y uint8) uint8 {
		return uint8(float64(x)*(1-t) + float64(y)*t)
	}
	return color.RGBA{mix(a.R, b.R), mix(a.G, b.G), mix(a.B, b.B), 0xFF}
}
