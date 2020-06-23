package cloud

import (
	"fmt"

	"github.com/bidease/spl"
)

// Instance ..
type Instance struct {
	ID                 string
	OpenstackUUID      string `json:"openstack_uuid"`
	Status             string
	Name               string
	RegionID           uint   `json:"region_id"`
	RegionCode         string `json:"region_code"`
	FlavorID           string `json:"flavor_id"`
	FlavorName         string `json:"flavor_name"`
	ImageID            string `json:"image_id"`
	ImageName          string `json:"image_name"`
	PublicIPv4Address  string `json:"public_ipv4_address"`
	PublicIPv6Address  string `json:"public_ipv6_address"`
	PrivateIPv4Address string `json:"private_ipv4_address"`
	IPv6Enabled        bool   `json:"ipv6_enabled"`
	GPNEnabled         bool   `json:"gpn_enabled"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// GetInstances ..
func GetInstances() (*[]Instance, error) {
	var instances []Instance
	page := 1

	for {
		var cInstancesTmp []Instance
		res, err := spl.RequestGet(fmt.Sprintf("cloud_computing/instances?per_page=100&page=%d", page), &cInstancesTmp)
		if err != nil {
			return nil, err
		}

		total := spl.StrToInt(res.Header.Get("X-Total"))

		for _, instance := range cInstancesTmp {
			instances = append(instances, instance)
		}

		if len(instances) == total {
			break
		}

		page++
	}

	return &instances, nil
}

// GetInstance ..
func GetInstance(id string) (*Instance, error) {
	var instnce Instance

	_, err := spl.RequestGet(fmt.Sprintf("cloud_computing/instances/%s", id), &instnce)
	if err != nil {
		return nil, err
	}

	return &instnce, nil
}
