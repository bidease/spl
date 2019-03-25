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

type cloudMachine struct {
	ID                          string
	Name                        string
	VCPUs                       uint
	RAM                         uint
	Disk                        uint
	FreeTrafficAmountGb         uint    `json:"free_traffic_amount_gb"`
	OvercommitTrafficPricePerGb float64 `json:"overcommit_traffic_price_per_gb"`
	GPUs                        uint
	DiscountedSnapshotSpace     float64 `json:"discounted_snapshot_space"`
	MonthlyPricesPerUnit        struct {
		Currency string
		Full     struct {
			Hosting struct {
				Price float64
				Tax   float64
				Total float64
			}
		}
		Original struct {
			Currency string
			Full     struct {
				Hosting struct {
					Price float64
					Tax   float64
					Total float64
				}
			}
		}
		PromoDiscountRate  float64 `json:"promo_discount_rate"`
		TaxIncludedInPrice bool    `json:"tax_included_in_price"`
	} `json:"monthly_prices_per_unit"`
	// TODO: monthly_windows_price_per_unit
	// TODO: prices_per_unit
	// TODO: windows_price_per_unit
}

func getCloudMachines(regionID uint) []cloudMachine {
	type rawCloudMachines struct {
		Data     []cloudMachine
		NumFound uint
	}
	var rawData rawCloudMachines
	tools.Request(fmt.Sprintf("cloud_computing/regions/%d/flavors", regionID), &rawData)

	var cloudMachines []cloudMachine
	for _, item := range rawData.Data {
		cloudMachineItem := cloudMachine{
			ID:                          item.ID,
			Name:                        item.Name,
			VCPUs:                       item.VCPUs,
			RAM:                         item.RAM,
			Disk:                        item.Disk,
			FreeTrafficAmountGb:         item.FreeTrafficAmountGb,
			OvercommitTrafficPricePerGb: item.OvercommitTrafficPricePerGb,
			GPUs:                        item.GPUs,
			DiscountedSnapshotSpace:     item.DiscountedSnapshotSpace,
		}
		cloudMachineItem.MonthlyPricesPerUnit.Currency = item.MonthlyPricesPerUnit.Currency
		cloudMachineItem.MonthlyPricesPerUnit.Full.Hosting.Price = item.MonthlyPricesPerUnit.Full.Hosting.Price
		cloudMachineItem.MonthlyPricesPerUnit.Full.Hosting.Tax = item.MonthlyPricesPerUnit.Full.Hosting.Tax
		cloudMachineItem.MonthlyPricesPerUnit.Full.Hosting.Total = item.MonthlyPricesPerUnit.Full.Hosting.Total
		cloudMachineItem.MonthlyPricesPerUnit.Original.Currency = item.MonthlyPricesPerUnit.Original.Currency
		cloudMachineItem.MonthlyPricesPerUnit.Original.Full.Hosting.Price = item.MonthlyPricesPerUnit.Original.Full.Hosting.Price
		cloudMachineItem.MonthlyPricesPerUnit.Original.Full.Hosting.Tax = item.MonthlyPricesPerUnit.Original.Full.Hosting.Tax
		cloudMachineItem.MonthlyPricesPerUnit.Original.Full.Hosting.Total = item.MonthlyPricesPerUnit.Original.Full.Hosting.Total
		cloudMachines = append(cloudMachines, cloudMachineItem)
	}
	return cloudMachines
}
