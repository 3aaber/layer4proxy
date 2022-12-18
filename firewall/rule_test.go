package firewall

import (
	"fmt"
	"layer4proxy/config"
	"net"
	"reflect"
	"testing"
)

func TestRule(t *testing.T) {

	ip1 := net.ParseIP("127.0.0.1")
	_, ipnet2, _ := net.ParseCIDR("127.0.0.1/24")

	rules := []struct {
		rule   config.Rule
		parsed FirewallRule
		err    error
	}{
		{
			rule:   config.Rule{Access: "allow", Match: "IR"},
			parsed: FirewallRule{Type: COUNTRY, Country: "US", Ip: nil, Network: nil},
			err:    fmt.Errorf("problem in parse country, expected: %s , got: %s", "IR", "IR"),
		},
		{
			rule:   config.Rule{Access: "allow", Match: "IR"},
			parsed: FirewallRule{Type: COUNTRY, Country: "IR", Ip: nil, Network: nil},
			err:    nil,
		},
		{
			rule:   config.Rule{Access: "allow", Match: "127.0.0.1"},
			parsed: FirewallRule{Type: NOTASSIGNED, Country: "", Ip: nil, Network: ipnet2},
			err:    nil,
		},
		{
			rule:   config.Rule{Access: "allow", Match: "127.0.0.1/24"},
			parsed: FirewallRule{Type: NOTASSIGNED, Country: "", Ip: &ip1, Network: nil},
			err:    nil,
		},
	}

	for _, item := range rules {
		parsed, err := ParseAccessRule(item.rule)
		if err == nil {
			if reflect.DeepEqual(parsed, item.parsed) && item.err == nil {
				t.Errorf("Expected: %v, got: %v", item.parsed, parsed)
			} else {
				t.Logf("Expected: %v, got: %v", item.parsed, parsed)
			}
		} else {
			if err.Error() != item.err.Error() {
				t.Errorf("Expected error: %s, got: %s", item.err, err)
			} else {
				t.Logf("Expected error: %s, got: %s", item.err, err)
			}
		}
	}
}
