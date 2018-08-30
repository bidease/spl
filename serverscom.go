package main

const (
	hostsURL    = "/hosts"
	hostURL     = "/hosts/%d"
	balanceURL  = "/statement/balance"
	servicesURL = "/hosts/%d/services"
)

type commonInfoHost struct {
	ID                 uint64
	Title              string
	Conf               string
	ScheduledReleaseAt string `json:"scheduled_release_at"`
	Location           struct {
		Name     string
		Timezone string
	}
	Networks []struct {
		ID       uint64
		HostIP   string `json:"host_ip"`
		PoolType string `json:"pool_type"` // drac, private, public
		Netmask  string
	}
	OSReinstallation bool `json:"os_reinstallation"`
}

type hosts struct {
	Data []struct {
		commonInfoHost
	}
	NumFound uint64 `json:"num_found"`
}

type host struct {
	Data struct {
		commonInfoHost
		OS struct {
			Arch    string
			Name    string
			Version string
		}
		Server struct {
			ID                  uint64
			Configuration       string
			OriginID            uint64 `json:"origin_id"`
			ChassisModelID      uint64 `json:"chassis_model_id"`
			ChassisModelName    string `json:"chassis_model_name"`
			ChassisModelCPUName string `json:"chassis_model_cpu_name"`
			RaidModelName       string `json:"raid_model_name"`
		}
	}
}

type balance struct {
	Data struct {
		Balance          string
		EstimatedBalance string `json:"estimated_balance"`
	}
}

type services struct {
	Data []struct {
		Currency         string
		OriginalCurrency string  `json:"original_currency"`
		OriginalPrice    float64 `json:"original_price"`
		Price            float64
	}
}
