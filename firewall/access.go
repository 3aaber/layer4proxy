package firewall

/**
 * firewall.go - firewall
 *
 */

import (
	"errors"
	"net"

	"layer4proxy/config"
)

/**
 * Access defines firewall rules chain
 */
type Access struct {
	AllowDefault bool
	Rules        []FirewallRule
}

/**
 * Creates new Access based on config
 */
func NewAccess(cfg config.AccessConfig) (*Access, error) {

	/// TODO
	// if cfg == nil {
	// 	return nil, errors.New("AccessConfig is nil")
	// }

	if cfg.Default == "" {
		cfg.Default = "allow"
	}

	if cfg.Default != "allow" && cfg.Default != "deny" {
		return nil, errors.New("AccessConfig Unexpected Default: " + cfg.Default)
	}

	access := Access{
		AllowDefault: cfg.Default == "allow",
		Rules:        []FirewallRule{},
	}

	// Parse rules
	for _, r := range cfg.Rules {
		rule, err := ParseAccessRule(r)
		if err != nil {
			return nil, err
		}
		access.Rules = append(access.Rules, *rule)
	}

	return &access, nil
}

/**
 * Checks if ip is allowed
 */
func (a *Access) Allows(ip *net.IP) bool {

	for _, r := range a.Rules {
		if r.Matches(ip) {
			return r.Allows()
		}
	}

	return a.AllowDefault
}
