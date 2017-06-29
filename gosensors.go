package gosensors

// #cgo LDFLAGS: -lsensors
// #include <stdlib.h>
// #include <stdio.h>
// #include <sensors/sensors.h>
import "C"

import (
	"log"
	"unsafe"
)

type Chip struct {
	Prefix string
	Bus    Bus
	AdapterName string
	Addr   int32
	Path   string
	Features []Feature
}

type Feature struct {
	Name    	string
	Number  	int32
	Type    	int32
	Lable		string
	Value		float64
	SubFeatures []SubFeature
}

type SubFeature struct {
	Name    string
	Number  int32
	Type    int32
	Mapping int32
	Flags   uint32
	Value	float64
}

type Bus struct {
	Type int16
	Nr   int16
}

func Init() {
	filename := C.CString("/etc/sensors3.conf")
	defer C.free(unsafe.Pointer(filename))

	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))

	fp, err := C.fopen(filename, mode)
	defer C.fclose(fp)

	if fp == nil {
		log.Fatal(err)
	}

	C.sensors_init(fp)
}

func Cleanup() {
	C.sensors_cleanup()
}

func GetDetectedChips() []Chip {
	var chips []Chip

	var count C.int = 0

	for {
		resp := C.sensors_get_detected_chips(nil, &count)

		if resp == nil {
			break
		}

		bus := Bus{
			Type: int16(resp.bus._type),
			Nr:   int16(resp.bus.nr),
		}

		var adapterName string

		if bus.Type == -1 {
			adapterName = "*"
		} else {
			adapterName = C.GoString(C.sensors_get_adapter_name(&resp.bus))
		}

		chip := Chip{
			Prefix: C.GoString(resp.prefix),
			Bus:    bus,
			AdapterName: adapterName,
			Addr:   int32(resp.addr),
			Path:   C.GoString(resp.path),
			Features:	getFeatures(resp),
		}

		chips = append(chips, chip)
	}

	return chips
}

func getFeatures(chip *C.struct_sensors_chip_name) []Feature {
	var features []Feature

	var count C.int = 0

	for {
		resp := C.sensors_get_features(chip, &count)

		if resp == nil {
			break
		}

		subfeatures := getSubFeatures(chip, resp)

		feature := Feature{
			Name:  	C.GoString(resp.name),
			Number:	int32(resp.number),
			Type:  	int32(resp._type),
			Lable:	getLabel(chip, resp),
			Value:	subfeatures[0].Value,
			SubFeatures: getSubFeatures(chip, resp),
		}

		features = append(features, feature)
	}

	return features
}

func getSubFeatures(chip *C.struct_sensors_chip_name, feature *C.struct_sensors_feature) []SubFeature {
	var subfeatures []SubFeature

	var count C.int = 0

	for {
		resp := C.sensors_get_all_subfeatures(chip, feature, &count)

		if resp == nil {
			break
		}

		subfeature := SubFeature{
			Name:    C.GoString(resp.name),
			Number:  int32(resp.number),
			Type:    int32(resp._type),
			Mapping: int32(resp.mapping),
			Flags:   uint32(resp.flags),
			Value:	 getValue(chip, int32(resp.number)),
		}

		subfeatures = append(subfeatures, subfeature)
	}

	return subfeatures
}

func getValue(chip *C.struct_sensors_chip_name, number int32) float64 {
	var value C.double
	C.sensors_get_value(chip, C.int(number), &value)
	return float64(value)
}

func getLabel(chip *C.struct_sensors_chip_name, feature *C.struct_sensors_feature) string {
	clabel := C.sensors_get_label(chip, feature)
	golabel := C.GoString(clabel)
	C.free(unsafe.Pointer(clabel))
	return golabel
}
