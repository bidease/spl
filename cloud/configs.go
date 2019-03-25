package cloud

import (
	"fmt"

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

type image struct {
	CreatedAt       string `json:"created_at"`
	ID              string
	Name            string
	DisplayPriority uint     `json:"display_priority"`
	RequiresSSHKey  bool     `json:"requires_ssh_key"`
	MinDisk         uint     `json:"min_disk"`
	IsWindows       bool     `json:"is_windows"`
	AllowedFlavors  []string `json:"allowed_flavors"`
}

func getImages(regionID uint) []image {
	type rawImages struct {
		Data     []image
		NumFound uint
	}
	var rawData rawImages
	tools.Request(fmt.Sprintf("cloud_computing/regions/%d/images", regionID), &rawData)

	var images []image
	for _, item := range rawData.Data {
		images = append(images, image{
			CreatedAt:       item.CreatedAt,
			ID:              item.ID,
			Name:            item.Name,
			DisplayPriority: item.DisplayPriority,
			RequiresSSHKey:  item.RequiresSSHKey,
			MinDisk:         item.MinDisk,
			IsWindows:       item.IsWindows,
			AllowedFlavors:  item.AllowedFlavors,
		})
	}
	return images
}
