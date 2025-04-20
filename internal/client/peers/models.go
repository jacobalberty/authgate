package peers

import "time"

type Peer struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	IP            string    `json:"ip"`
	ConnectionIP  string    `json:"connection_ip"`
	Connected     bool      `json:"connected"`
	LastSeen      time.Time `json:"last_seen"`
	Os            string    `json:"os"`
	KernelVersion string    `json:"kernel_version"`
	GeonameID     int       `json:"geoname_id"`
	Version       string    `json:"version"`
	Groups        []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		PeersCount     int    `json:"peers_count"`
		ResourcesCount int    `json:"resources_count"`
		Issued         string `json:"issued"`
	} `json:"groups"`
	SSHEnabled                  bool      `json:"ssh_enabled"`
	UserID                      string    `json:"user_id"`
	Hostname                    string    `json:"hostname"`
	UIVersion                   string    `json:"ui_version"`
	DNSLabel                    string    `json:"dns_label"`
	LoginExpirationEnabled      bool      `json:"login_expiration_enabled"`
	LoginExpired                bool      `json:"login_expired"`
	LastLogin                   time.Time `json:"last_login"`
	InactivityExpirationEnabled bool      `json:"inactivity_expiration_enabled"`
	ApprovalRequired            bool      `json:"approval_required"`
	CountryCode                 string    `json:"country_code"`
	CityName                    string    `json:"city_name"`
	SerialNumber                string    `json:"serial_number"`
	ExtraDNSLabels              []string  `json:"extra_dns_labels"`
	AccessiblePeersCount        int       `json:"accessible_peers_count"`
}
