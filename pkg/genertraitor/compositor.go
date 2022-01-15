package genertraitor

import (
	"image"
	"image/draw"
	"log"
	"os"
)

// Compositor - used to arrange the variations in the final image
type Compositor struct {
	Strategy CompositionStrategy
}

// NewDefaultCompositor - creates a new compositor with default configuration
func NewDefaultCompositor() *Compositor {
	return &Compositor{
		Strategy: CompositionStrategyOverlap,
	}
}

// Compose - creates an image from the variations using the configured strategy
func (c *Compositor) Compose(variations []*Variation) (image.Image, error) {
	return c.Strategy(variations)
}

// CompositionStrategy - determines how the variations are arranged in relation to one-another
type CompositionStrategy func([]*Variation) (image.Image, error)

// CompositionStrategyOverlap - arranges the variations on top of one another
// the dimensions of the final image are the maximum width and height across all variations
var CompositionStrategyOverlap = func(variations []*Variation) (image.Image, error) {
	// open all the variation files
	files := []*os.File{}
	for _, variation := range variations {
		f, err := os.Open(variation.Path)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	// decode them
	decoded := []image.Image{}
	for _, f := range files {
		d, _, err := image.Decode(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		decoded = append(decoded, d)
	}
	// get max width and height of variants for final image
	width, height := decoded[0].Bounds().Dx(), decoded[0].Bounds().Dy()
	for _, img := range decoded {
		if w := img.Bounds().Dx(); w > width {
			width = w
		}
		if h := img.Bounds().Dy(); h > height {
			height = h
		}
	}
	// create the base image
	base := image.NewNRGBA(image.Rect(0, 0, width, height))
	// overlay the variants
	for _, img := range decoded {
		draw.Draw(base, img.Bounds(), img, image.ZP, draw.Over)
	}
	return base, nil
}

// CompositionStrategyHorizontal - arranges the images next to one another with the first variation on the left
// the dimensions of the final image are the combined width of all variations and the maximum height across all variations
var CompositionStrategyHorizontal = func(variations []*Variation) (image.Image, error) {
	// open all the variation files
	files := []*os.File{}
	for _, variation := range variations {
		f, err := os.Open(variation.Path)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	// decode them
	decoded := []image.Image{}
	for _, f := range files {
		d, _, err := image.Decode(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		decoded = append(decoded, d)
	}
	// get combined width and max height for final image
	width, height := 0, 0
	for _, img := range decoded {
		width += img.Bounds().Dx()
		if h := img.Bounds().Dy(); h > height {
			height = h
		}
	}
	log.Printf("width: %d, height: %d", width, height)
	// create the base image
	base := image.NewNRGBA(image.Rect(0, 0, width, height))
	// overlay the variants
	dx := 0
	for _, img := range decoded {
		log.Printf("dx: %d", dx)
		pt := image.Point{X: dx, Y: 0}
		r := image.Rectangle{pt, pt.Add(img.Bounds().Size())}
		draw.Draw(base, r, img, image.Point{0, 0}, draw.Src)
		dx += img.Bounds().Dx()
	}
	return base, nil
}
