package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/mcclurejt/genertraitor/pkg/genertraitor"
)

func main() {
	gt, err := genertraitor.New(genertraitor.Config{
		TraitConfigs: []*genertraitor.TraitConfig{
			{
				Path:             "./traits/background",
				RarityMultiplier: genertraitor.RarityMultiplierDefault,
			},
			{
				Path:             "./traits/tall",
				RarityMultiplier: genertraitor.RarityMultiplierDefault,
			},
			{
				Path:             "./traits/medium",
				RarityMultiplier: genertraitor.RarityMultiplierDefault,
			},
			{
				Path:             "./traits/short",
				RarityMultiplier: genertraitor.RarityMultiplierDefault,
			},
		},
		CompositionStrategy: genertraitor.CompositionStrategyOverlap,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range gt.Traits {
		fmt.Printf("Trait: %s\n", t.Name)
		for _, v := range t.Variations {
			fmt.Printf(" - Variation: %s\n", v.Name)
		}
	}
	n := 239487
	img, err := gt.Generate(n)
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create(fmt.Sprintf("./generated/%s.png", strconv.Itoa(n)))
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}
