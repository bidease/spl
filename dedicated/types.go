package dedicated

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

// Service ..
type Service struct {
	ID            string
	Name          string
	Type          string
	Currency      string
	Lable         string
	DateFrom      string  `json:"date_from"`
	DateTo        string  `json:"date_to"`
	UsageQuantity float64 `json:"usage_quantity"`
	Tax           float64
	Total         float64
	Subtotal      float64
	DiscountRate  float64 `json:"discount_rate"`
}
