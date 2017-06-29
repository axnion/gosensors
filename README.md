# gosensors
Go bindings for libsensors.so from the lm-sensors project via cgo. Forked from github.com/md14454/gosensors where I have later made changes so objects returned don't use methods, instead all data is stored in the object when returned. This is to make mocking of this system easier and make systems using gosensors easier to test. But because of the changes I've done, the interface for this library will not match the interface of the original.


### Example
``` go
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
```
