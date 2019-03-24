package cloud

import (
	"github.com/bidease/spl/tools"
)

type region struct {
	ID   string
	Name string
}

func getRegions() []region {
	type rawRegions struct {
		Data     []region
		NumFound uint
	}
	var rawData rawRegions
	tools.Request("cloud_computing/regions", &rawData)

	var regions []region
	for _, item := range rawData.Data {
		regions = append(regions, region{ID: item.ID, Name: item.Name})
	}

	return regions
}
