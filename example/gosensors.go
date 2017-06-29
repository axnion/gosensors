package main

import (
	"fmt"

	"github.com/axnion/gosensors"
)

func main() {
	gosensors.Init()
	defer gosensors.Cleanup()

	chips := gosensors.GetDetectedChips()

	for i := 0; i < len(chips); i++ {
		chip := chips[i]

		fmt.Printf("Adapter: %v\n", chip.AdapterName)

		features := chip.Features

		for j := 0; j < len(features); j++ {
			feature := features[j]

			fmt.Printf("%v ('%v'): %.1f\n", feature.Name, feature.Lable, feature.Value)

			subfeatures := feature.SubFeatures

			for k := 0; k < len(subfeatures); k++ {
				subfeature := subfeatures[k]

				fmt.Printf("  %v: %.1f\n", subfeature.Name, subfeature.Value)
			}
		}

		fmt.Printf("\n")
	}
}
