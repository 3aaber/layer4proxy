package firewall

import (
	"errors"
	"fmt"
	"layer4proxy/config"
	"net"
	"strconv"

	"github.com/biter777/countries"
)

/**
 * FirewallRule defines order (firewall, deny)
 * and IP or Network
 */
type FirewallRule struct {
	Allow   bool
	Type    RuleType
	Country string
	Ip      *net.IP
	Network *net.IPNet
	ASN     uint64
}

/**
 * Rule Type (IP,CIDR,Country,ASN)
 */
type RuleType int

const (
	NOTASSIGNED RuleType = 0
	IP          RuleType = 1
	CIDR        RuleType = 2
	COUNTRY     RuleType = 3
	ASN         RuleType = 4
)

/**
 * ParseAccessRule Parse AccessRule string
 */
func ParseAccessRule(rule config.Rule) (*FirewallRule, error) {

	r := rule.Access
	cidrOrIpOrCountry := rule.Match

	if r != "allow" && r != "deny" {
		return nil, fmt.Errorf("cant parse rule definition : %v ", rule)
	}

	// try check if cidrOrIpOrCountry is country and handle
	countryShould := countries.ByName(cidrOrIpOrCountry)
	if countryShould != countries.Unknown {
		return &FirewallRule{
			Allow:   r == "allow",
			Type:    COUNTRY,
			Country: cidrOrIpOrCountry,
			Ip:      nil,
			Network: nil,
			ASN:     0,
		}, nil
	}

	// try check if cidrOrIpOrCountry is ip and handle

	ipShould := net.ParseIP(cidrOrIpOrCountry)
	if ipShould != nil {
		return &FirewallRule{
			Allow:   r == "allow",
			Type:    IP,
			Country: "",
			Ip:      &ipShould,
			Network: nil,
			ASN:     0,
		}, nil
	}

	_, ipNetShould, _ := net.ParseCIDR(cidrOrIpOrCountry)
	if ipNetShould != nil {
		return &FirewallRule{
			Allow:   r == "allow",
			Type:    CIDR,
			Country: "",
			Ip:      nil,
			Network: ipNetShould,
			ASN:     0,
		}, nil
	}

	// try check if cidrOrIpOrCountry is asn and handle
	asn, isASN := checkASN(cidrOrIpOrCountry)
	if isASN {
		return &FirewallRule{
			Allow:   r == "allow",
			Type:    ASN,
			Country: "",
			Ip:      nil,
			Network: nil,
			ASN:     asn,
		}, nil
	}

	return nil, errors.New("Cant parse acces rule target, not an ip or cidr: " + cidrOrIpOrCountry)

}

func checkASN(asn string) (uint64, bool) {
	parsedASN, err := strconv.ParseUint(asn, 10, 32)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	if parsedASN > 0 && parsedASN < 4294967295 {
		return parsedASN, true
	}
	return 0, false
}

/**
 * Checks if ip matches access rule
 */
func (f *FirewallRule) Matches(ip *net.IP) bool {

	switch f.Type {
	case IP:
		return (*f.Ip).Equal(*ip)
	case CIDR:
		return f.Network.Contains(*ip)
	case COUNTRY:
		country, err := GetCountry(*ip)
		if err != nil {
			fmt.Printf("Cant get country for ip: %v", *ip)
			return true
		}
		return f.Country == country
	case ASN:
		asn, err := GetASN(*ip)
		if err != nil {
			fmt.Printf("Cant get asn for ip: %v", *ip)
			return true
		}
		return f.ASN == asn
	default:
		fmt.Printf("Unknown rule type: %v", f.Type)
	}

	return false
}

/**
 * Checks is it's allow or deny rule
 */
func (f *FirewallRule) Allows() bool {
	return f.Allow
}
