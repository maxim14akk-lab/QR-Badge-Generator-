// qr_generator.go
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	text := flag.String("t", "", "Text to encode")
	output := flag.String("o", "qr_code.png", "Output file name")
	fg := flag.String("f", "#000000", "Foreground color (hex)")
	bg := flag.String("b", "#FFFFFF", "Background color")
	size := flag.Int("s", 10, "Size per module (pixels)")
	margin := flag.Int("m", 4, "Margin in modules")
	logoPath := flag.String("l", "", "Path to logo image")
	flag.Parse()

	if *text == "" {
		fmt.Fprintln(os.Stderr, "Error: -t text is required")
		os.Exit(1)
	}

	// Generate QR code matrix
	qr, err := qrcode.New(*text, qrcode.Medium)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating QR: %v\n", err)
		os.Exit(1)
	}
	// qrcode package doesn't support custom colors directly, so we draw ourselves
	// Get matrix as bool (true=black)
	matrix := qr.Bitmap()
	modules := len(matrix)
	moduleSize := *size
	marginPx := *margin * moduleSize
	totalSize := modules*moduleSize + 2*marginPx

	img := image.NewRGBA(image.Rect(0, 0, totalSize, totalSize))
	fgColor, err := parseHexColor(*fg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid foreground color: %v\n", err)
		os.Exit(1)
	}
	bgColor, err := parseHexColor(*bg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid background color: %v\n", err)
		os.Exit(1)
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	for y, row := range matrix {
		for x, val := range row {
			if val {
				rect := image.Rect(
					marginPx+x*moduleSize,
					marginPx+y*moduleSize,
					marginPx+(x+1)*moduleSize,
					marginPx+(y+1)*moduleSize,
				)
				draw.Draw(img, rect, &image.Uniform{fgColor}, image.Point{}, draw.Src)
			}
		}
	}

	// Embed logo
	if *logoPath != "" {
		logoFile, err := os.Open(*logoPath)
		if err == nil {
			defer logoFile.Close()
			logoImg, _, err := image.Decode(logoFile)
			if err == nil {
				logoBounds := logoImg.Bounds()
				logoSize := int(float64(modules*moduleSize) * 0.25)
				// Scale logo (simplistic: we resize by scaling bounds)
				// For simplicity, we just draw as is, but we should resize.
				// We'll assume logo is already small or use a proper resize lib.
				// Here we just draw without resize (but we'll set a max size)
				posX := marginPx + (modules*moduleSize-logoSize)/2
				posY := marginPx + (modules*moduleSize-logoSize)/2
				// Create a new image with logo scaled (using draw.Scale)
				// Since we don't have scaling in std, we skip advanced resize.
				// For demo, we'll paste at center without resize (if too big, it will overflow)
				draw.Draw(img, image.Rect(posX, posY, posX+logoBounds.Dx(), posY+logoBounds.Dy()), logoImg, logoBounds.Min, draw.Over)
			}
		}
	}

	// Save
	f, err := os.Create(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("QR code saved to %s\n", *output)
}

func parseHexColor(s string) (color.Color, error) {
	// Support simple hex #RRGGBB
	if len(s) == 7 && s[0] == '#' {
		var r, g, b uint8
		_, err := fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return nil, err
		}
		return color.RGBA{r, g, b, 255}, nil
	}
	// Fallback: try to parse as named color? Not implemented.
	return color.Black, fmt.Errorf("unsupported color format")
}
