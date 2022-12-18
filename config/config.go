package config

type Backend struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Priority int    `json:"priority"`
	Weight   int    `json:"weight"`
	Sni      string `json:"sni"`
}

/**
 * Server udp options
 * for protocol = "udp"
 */
type Udp struct {
	MaxRequests  uint64 `json:"max_requests"`
	MaxResponses uint64 `json:"max_responses"`
	Transparent  bool   `json:"transparent"`
}

/**
 * Access configuration
 */
type AccessConfig struct {
	Default string `json:"default"`
	Rules   []Rule `json:"rules"`
}

// AccessRule defines access rule
type Rule struct {
	Access string `json:"access"` // Access Rule : allow | deny
	Match  string `json:"match"`  // Match Rule : IP | CIDR | Country | ASN
}
