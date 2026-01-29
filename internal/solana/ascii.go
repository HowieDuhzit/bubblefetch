package solana

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	defaultASCIIWidth  = 30
	defaultASCIIHeight = 15
)

// FetchAndConvertLogo downloads a token logo and converts it to ASCII art
func FetchAndConvertLogo(logoURI string) (string, error) {
	if logoURI == "" {
		return getDefaultTokenASCII(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", logoURI, nil)
	if err != nil {
		return getDefaultTokenASCII(), nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return getDefaultTokenASCII(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return getDefaultTokenASCII(), nil
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return getDefaultTokenASCII(), nil
	}

	return imageToASCII(img, defaultASCIIWidth, defaultASCIIHeight), nil
}

// imageToASCII converts an image to ASCII art
func imageToASCII(img image.Image, width, height int) string {
	bounds := img.Bounds()
	asciiChars := []rune(" .:-=+*#%@")

	var result strings.Builder

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Map ASCII coordinates to image coordinates
			imgX := bounds.Min.X + (x * bounds.Dx() / width)
			imgY := bounds.Min.Y + (y * bounds.Dy() / height)

			// Get pixel color
			r, g, b, _ := img.At(imgX, imgY).RGBA()

			// Convert to grayscale (0-255)
			gray := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0

			// Map grayscale to ASCII character
			charIndex := int(gray / 255.0 * float64(len(asciiChars)-1))
			if charIndex < 0 {
				charIndex = 0
			}
			if charIndex >= len(asciiChars) {
				charIndex = len(asciiChars) - 1
			}

			result.WriteRune(asciiChars[charIndex])
		}
		if y < height-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}

// getDefaultTokenASCII returns a default ASCII art for tokens when logo is unavailable
func getDefaultTokenASCII() string {
	return `
     ████████████████████████████████
   ████████████████████████████████
 ████████████████████████████████
████████████████████████████████


████████████████████████████████
 ████████████████████████████████
   ████████████████████████████████
     ████████████████████████████████


     ████████████████████████████████
   ████████████████████████████████
 ████████████████████████████████
████████████████████████████████
`
}

// GeneratePriceChart creates an ASCII sparkline chart from price data
func GeneratePriceChart(prices []PricePoint, width int) string {
	if len(prices) == 0 {
		return generateMockChart(width)
	}

	// Extract price values
	values := make([]float64, len(prices))
	for i, p := range prices {
		values[i] = p.Price
	}

	return generateSparkline(values, width)
}

// generateSparkline creates a sparkline chart from values
func generateSparkline(values []float64, width int) string {
	if len(values) == 0 {
		return ""
	}

	// Sparkline characters (from lowest to highest)
	chars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	// Find min and max
	min := values[0]
	max := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// Handle case where all values are the same
	if max == min {
		var result strings.Builder
		for i := 0; i < width && i < len(values); i++ {
			result.WriteRune(chars[len(chars)/2])
		}
		return result.String()
	}

	// Sample values to fit width
	sampled := sampleValues(values, width)

	// Convert to sparkline
	var result strings.Builder
	for _, v := range sampled {
		normalized := (v - min) / (max - min)
		charIndex := int(normalized * float64(len(chars)-1))
		if charIndex < 0 {
			charIndex = 0
		}
		if charIndex >= len(chars) {
			charIndex = len(chars) - 1
		}
		result.WriteRune(chars[charIndex])
	}

	return result.String()
}

// sampleValues samples an array to a target length
func sampleValues(values []float64, targetLength int) []float64 {
	if len(values) <= targetLength {
		return values
	}

	sampled := make([]float64, targetLength)
	ratio := float64(len(values)) / float64(targetLength)

	for i := 0; i < targetLength; i++ {
		index := int(float64(i) * ratio)
		if index >= len(values) {
			index = len(values) - 1
		}
		sampled[i] = values[index]
	}

	return sampled
}

// generateMockChart creates a mock price chart when real data is unavailable
func generateMockChart(width int) string {
	// Generate a sine wave pattern for demonstration
	values := make([]float64, width)
	for i := 0; i < width; i++ {
		values[i] = 50 + 30*math.Sin(float64(i)*0.3)
	}
	return generateSparkline(values, width)
}

// GeneratePriceChangeIndicator creates a visual indicator for price change
func GeneratePriceChangeIndicator(priceChange float64) string {
	if priceChange > 0 {
		return fmt.Sprintf("📈 +%.2f%%", priceChange)
	} else if priceChange < 0 {
		return fmt.Sprintf("📉 %.2f%%", priceChange)
	}
	return "📊 0.00%"
}
