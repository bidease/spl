package dedicated

import (
	"fmt"

	"github.com/bidease/spl"
)

// Server ..
type Server struct {
	ID                 string
	Title              string
	LocationID         int `json:"location_id"`
	Status             string
	Configuration      string
	PrivateIPv4Address string `json:"private_ipv4_address"`
	PublicIPv4Address  string `json:"public_ipv4_address"`
	LeaseStartAt       string `json:"lease_start_at"`
	ScheduledReleaseAt string `json:"scheduled_release_at"`
	Type               string
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
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
