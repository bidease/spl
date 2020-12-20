package dedicated

import (
	"fmt"

	"github.com/bidease/spl"
)

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

// GetServices ..
func GetServices(id string) (*[]Service, error) {
	var services []Service

	_, err := spl.RequestGet(fmt.Sprintf("hosts/dedicated_servers/%s/services", id), &services)
	if err != nil {
		return nil, err
	}

	return &services, nil
}
