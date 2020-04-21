package spl

// Location ..
type Location struct {
	ID                   int
	Name                 string
	SupportedFeatures    []string `json:"supported_features"`
	L2SegmentsEnabled    bool     `json:"l2_segments_enabled"`
	PrivateRacksEnabled  bool     `json:"private_racks_enabled"`
	LoadBalancersEnabled bool     `json:"load_balancers_enabled"`
}

// GetLocations ..
func GetLocations() (*[]Location, error) {
	var location []Location

	_, err := RequestGet("locations", &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}
