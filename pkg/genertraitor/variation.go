package genertraitor

import (
	"path/filepath"
	"strings"
)

// Variation - defines one variation of a trait
type Variation struct {
	Name   string // name of the variation, derived from filename
	Path   string // path to the variation file
	Rarity int    // rarity of the variation
}

// NewVariation - creates a new variation from the supplied filepath
func NewVariation(variationPath string, rarity int) (*Variation, error) {
	return &Variation{
		Name:   strings.Split(filepath.Base(variationPath), ".")[0],
		Path:   variationPath,
		Rarity: rarity,
	}, nil
}
