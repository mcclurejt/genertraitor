package genertraitor

import (
	"fmt"
	"image"
	"log"
	"math"
)

// Genertraitor - handles selecting variants to be composited
type Genertraitor struct {
	Traits     []*Trait
	Compositor *Compositor
}

// Config - contains the configuration needed to populate a Genertraitor
type Config struct {
	TraitConfigs []*TraitConfig
	CompositionStrategy
}

// New - creates a new Genertraitor
func New(config Config) (*Genertraitor, error) {
	g := &Genertraitor{Traits: []*Trait{}}
	if config.CompositionStrategy == nil {
		g.Compositor = NewDefaultCompositor()
	} else {
		g.Compositor = &Compositor{Strategy: config.CompositionStrategy}
	}
	// populate traits
	for _, tc := range config.TraitConfigs {
		t, err := NewTrait(tc)
		if err != nil {
			return nil, err
		}
		g.Traits = append(g.Traits, t)
	}
	return g, nil
}

// Generate - creates a Compositor with the variations specified by the input integer
func (g *Genertraitor) Generate(n int) (image.Image, error) {
	variations, err := g.SelectVariants(n)
	if err != nil {
		return nil, err
	}
	return g.Compositor.Compose(variations)
}

// SelectVariants - chooses a single variation for each trait
func (g *Genertraitor) SelectVariants(n int) ([]*Variation, error) {
	nt := len(g.Traits)
	log.Printf("generating an image with %d traits...", nt)
	variations := []*Variation{}
	for _, trait := range g.Traits {
		if n <= 0 {
			return nil, fmt.Errorf("error: n's value of %d is <= 0", n)
		}
		rt := trait.RarityTable()
		mod := n % len(rt)
		variation := trait.Variations[rt[mod]]
		variations = append(variations, variation)
		log.Printf("trait: %s, variation: %s\n", trait.Name, variation.Name)
		n -= mod
		l10 := math.Log10(float64(len(rt)))
		n /= int(math.Pow10(int(math.Ceil(l10))))
	}
	return variations, nil
}

/*
	variations 9
	rarities 0, 0, 1, 2, 4, 4, 6, 8, 9

	r0 = 2 times chance of r1
	r1 = 2 times chance of r2
	r2 = 4 times chance of r4
	r4 = 4 times chance of r6
	r6 = 4 times chance of r8
	r8 = 2 times chance of r9
	r9 = r0 / 512

	probabilities 512, 512, 256, 128, 32, 32, 8, 2, 1

**/
