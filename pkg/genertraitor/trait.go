package genertraitor

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

// RarityMultiplierDefault - Default rarity multiplier (2)
var RarityMultiplierDefault = 2

// Trait - a trait characteristic of the resulting image. Only one variation of a trait can be present at any time
type Trait struct {
	Name       string       // name of the trait, derived from folder name
	Variations []*Variation // list of variations for the trait
	Config     *TraitConfig // trait-specific configuration
}

// TraitConfig - contains the configuration necessary to populate and composite a trait
type TraitConfig struct {
	Path             string
	RarityMultiplier int // the difference in rarity between levels, a rarity multiplier of 2 means that rarity 0 is 2x more likely to appear than rarity 1
}

// Validate - ensures a correct TraitConfig
func (tc *TraitConfig) Validate() error {
	if tc.RarityMultiplier <= 0 {
		return errors.New("must provide positive rarity multiplier")
	}
	return nil
}

// NewDefaultTraitConfig - returns a TraitConfig with the default rarity multiplier
func NewDefaultTraitConfig(traitPath string) *TraitConfig {
	return &TraitConfig{
		Path:             traitPath,
		RarityMultiplier: RarityMultiplierDefault,
	}
}

// NewTrait - creates a new trait from the provided path and populates it with variants
func NewTrait(tc *TraitConfig) (*Trait, error) {
	// validate config
	if err := tc.Validate(); err != nil {
		return nil, err
	}
	// path must exist
	fileInfo, err := os.Stat(tc.Path)
	if err != nil {
		return nil, err
	}
	// path must be a directory
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("trait path %s is not a directory", tc.Path)
	}
	// walk through each subdirectory containing variations and add them to the trait
	t := &Trait{Name: fileInfo.Name(), Variations: []*Variation{}, Config: tc}
	err = filepath.Walk(tc.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path != tc.Path && info.IsDir() {
			return t.addVariations(path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Trait) addVariations(variationDir string) error {
	// infer rarity from the variation directory name
	rarity, err := strconv.Atoi(filepath.Base(variationDir))
	if err != nil {
		return err
	}
	// walk through variations and add to trait
	pattern := "*.png"
	err = filepath.Walk(variationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			v, err := NewVariation(path, rarity)
			if err != nil {
				return err
			}
			t.Variations = append(t.Variations, v)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// MaxRarity - returns the highest rarity value of any variation for the trait
func (t *Trait) MaxRarity() int {
	max := 0
	for _, v := range t.Variations {
		if v.Rarity > max {
			max = v.Rarity
		}
	}
	return max
}

// RarityTable - returns a list of variation indices where the frequency of each indicy corresponds to the variation's rarity
func (t *Trait) RarityTable() []int {
	rt := []int{}
	mr := t.MaxRarity()
	baseFrequency := int(math.Pow(float64(t.Config.RarityMultiplier), float64(mr)))
	for i, v := range t.Variations {
		for j := 0; j < baseFrequency/int(math.Pow(float64(t.Config.RarityMultiplier), float64(v.Rarity))); j++ {
			rt = append(rt, i)
		}
	}
	return rt
}
