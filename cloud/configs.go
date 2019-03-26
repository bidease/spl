package cloud

import (
	"fmt"
	"log"

	"github.com/bidease/spl/common"
	"github.com/bidease/spl/tools"

	"github.com/urfave/cli"
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
	tools.GetRequest("cloud_computing/regions", &rawData)

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
	tools.GetRequest(fmt.Sprintf("cloud_computing/regions/%d/images", regionID), &rawData)

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

type cloudServer struct {
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

func getCloudServers(regionID uint) []cloudServer {
	type rawCloudServers struct {
		Data     []cloudServer
		NumFound uint
	}
	var rawData rawCloudServers
	tools.GetRequest(fmt.Sprintf("cloud_computing/regions/%d/flavors", regionID), &rawData)

	var cloudServers []cloudServer
	for _, item := range rawData.Data {
		cloudServerItem := cloudServer{
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
		cloudServerItem.MonthlyPricesPerUnit.Currency = item.MonthlyPricesPerUnit.Currency
		cloudServerItem.MonthlyPricesPerUnit.Full.Hosting.Price = item.MonthlyPricesPerUnit.Full.Hosting.Price
		cloudServerItem.MonthlyPricesPerUnit.Full.Hosting.Tax = item.MonthlyPricesPerUnit.Full.Hosting.Tax
		cloudServerItem.MonthlyPricesPerUnit.Full.Hosting.Total = item.MonthlyPricesPerUnit.Full.Hosting.Total
		cloudServerItem.MonthlyPricesPerUnit.Original.Currency = item.MonthlyPricesPerUnit.Original.Currency
		cloudServerItem.MonthlyPricesPerUnit.Original.Full.Hosting.Price = item.MonthlyPricesPerUnit.Original.Full.Hosting.Price
		cloudServerItem.MonthlyPricesPerUnit.Original.Full.Hosting.Tax = item.MonthlyPricesPerUnit.Original.Full.Hosting.Tax
		cloudServerItem.MonthlyPricesPerUnit.Original.Full.Hosting.Total = item.MonthlyPricesPerUnit.Original.Full.Hosting.Total
		cloudServers = append(cloudServers, cloudServerItem)
	}
	return cloudServers
}

type createCloudServer struct {
	ImageID        string `json:"image_id"`
	FlavorID       string `json:"flavor_id"` // configID
	Period         string `json:"period"`    // monthly
	SetupPassword  bool   `json:"setup_password"`
	KeyFingerprint string `json:"key_fingerprint"`
	Name           string `json:"name"`
	GPEnabled      bool   `json:"gp_enabled"`
	BackupEnabled  bool   `json:"backup_enabled"`
	BackupCopies   uint   `json:"backup_copies"`
}

// CreateCloudServer ..
func CreateCloudServer(c *cli.Context) {
	if len(c.String("fingerprint")) > 0 && c.Bool("password") {
		log.Fatalln("use -fingerprint or -password")
	}

	newCloudServer := createCloudServer{
		ImageID:   c.String("imageID"),
		FlavorID:  c.String("configID"),
		Period:    "monthly",
		Name:      c.String("name"),
		GPEnabled: c.Bool("gpn"),
	}

	if c.Uint("backups") > 0 {
		newCloudServer.BackupEnabled = true
		newCloudServer.BackupCopies = c.Uint("backups")
	}

	if c.Bool("password") {
		newCloudServer.SetupPassword = true
	} else {
		newCloudServer.KeyFingerprint = common.GetUserFingerprint(c.String("fingerprint"))
	}

	var response common.Response
	tools.PostRequest(fmt.Sprintf("cloud_computing/regions/%d/instances", c.Uint("regionID")), &response, newCloudServer)

	if !response.Success {
		log.Fatalf("%v\n", response.Success)
	}
}

// RegionInstance ..
type RegionInstance struct {
	ID         string
	Name       string
	Status     string
	FlavorName string `json:"flavor_name"`
	RegionID   string
	RegionName string
	LocalIP    string `json:"local_ip"`
	InternalIP string `json:"internal_ip"`
	ExternalIP string `json:"external_ip"`
}

func getExistsCloudServersRegion(regionIns region) []RegionInstance {
	var response struct {
		Data    []RegionInstance
		NumFoud uint `json:"num_found"`
	}
	tools.GetRequest(fmt.Sprintf("cloud_computing/regions/%s/instances", regionIns.ID), &response)

	var regionInstances []RegionInstance
	for _, item := range response.Data {
		regionInstances = append(regionInstances, RegionInstance{
			ID:         item.ID,
			Name:       item.Name,
			Status:     item.Status,
			FlavorName: item.FlavorName,
			RegionID:   regionIns.ID,
			RegionName: regionIns.Name,
			LocalIP:    item.LocalIP,
			InternalIP: item.InternalIP,
			ExternalIP: item.ExternalIP,
		})
	}

	return regionInstances
}

func getExistsCloudServers() []RegionInstance {
	var existsInstances []RegionInstance
	for _, region := range getRegions() {
		for _, instanse := range getExistsCloudServersRegion(region) {
			existsInstances = append(existsInstances, instanse)
		}
	}
	return existsInstances
}
