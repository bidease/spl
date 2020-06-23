package dedicated

import (
	"fmt"

	"github.com/bidease/spl"
)

// Server ..
type Server struct {
	ID                   string
	Title                string
	LocationID           uint   `json:"location_id"`
	LocationCode         string `json:"location_code"`
	Status               string
	OperationalStatus    string `json:"operational_status"`
	PowerStatus          string `json:"power_status"`
	Configuration        string
	ConfigurationDetails ServerConfigurationDetails
	PrivateIPv4Address   string `json:"private_ipv4_address"`
	PublicIPv4Address    string `json:"public_ipv4_address"`
	LeaseStartAt         string `json:"lease_start_at"`
	ScheduledReleaseAt   string `json:"scheduled_release_at"`
	Type                 string
	RackID               string `json:"rack_id"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

// ServerConfigurationDetails ..
type ServerConfigurationDetails struct {
	RAMSize                 uint   `json:"ram_size"`
	ServerModelID           uint   `json:"server_model_id"`
	ServerModelName         string `json:"server_model_name"`
	BandwidthID             uint   `json:"bandwidth_id"`
	BandwidthName           string `json:"bandwidth_name"`
	PrivateUplinkID         uint   `json:"private_uplink_id"`
	PrivateUplinkName       string `json:"private_uplink_name"`
	PublicUplinkID          uint   `json:"public_uplink_id"`
	PublicUplinkName        string `json:"public_uplink_name"`
	OperatingSystemID       string `json:"operating_system_id"`
	OperatingSystemFullName string `json:"operating_system_full_name"`
}

// GetServers ..
func GetServers() (*[]Server, error) {
	var servers []Server
	page := 1

	for {
		var serversTmp []Server

		res, err := spl.RequestGet(fmt.Sprintf("hosts?per_page=100&page=%d", page), &serversTmp)
		if err != nil {
			return nil, err
		}

		total := spl.StrToInt(res.Header.Get("X-Total"))

		for _, server := range serversTmp {
			servers = append(servers, server)
		}

		if len(servers) == total {
			break
		}

		page++
	}

	return &servers, nil
}

// GetServer ..
func GetServer(id string) (*Server, error) {
	var server Server

	_, err := spl.RequestGet(fmt.Sprintf("hosts/dedicated_servers/%s", id), &server)
	if err != nil {
		return nil, err
	}

	return &server, nil
}
