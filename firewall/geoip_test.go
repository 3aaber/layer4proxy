package firewall

import (
	"net"
	"testing"
)

func TestCountry(t *testing.T) {
	err := Setupdb("dbfiles/geoCountry.mmdb", "dbfiles/geoIsp.mmdb")
	if err != nil {
		t.Error(err)
	}
	defer Closedb()

	ips := []struct {
		ip      string
		country string
		err     error
	}{
		{"5.123.4.5", "IR", nil},
		{"5.123.4.6", "IR", nil},
		{"5.123.4.7", "IR", nil},
		{"5.123.4.8", "IR", nil},
	}

	for _, table := range ips {
		country, err := GetCountry(net.ParseIP(table.ip))
		if err != nil { // if there is an error, check if it is the expected error
			if err.Error() != table.err.Error() {
				t.Errorf("Expected error: %s, got: %s", table.err, err)
			} else {
				t.Logf("Expected error: %s, got: %s", table.err, err)
			} // if there is no error, check if the country is the expected one
		} else {
			if country != table.country {
				t.Errorf("expected %s, got %s", table.country, country)
			}
		}
	}
}

func TestASN(t *testing.T) {
	err := Setupdb("dbfiles/geoCountry.mmdb", "dbfiles/geoIsp.mmdb")
	if err != nil {
		t.Error(err)
	}
	defer Closedb()

	ips := []struct {
		ip  string
		asn uint64
		err error
	}{
		{"5.123.4.5", 44244, nil},
		{"5.123.4.6", 44244, nil},
		{"5.123.4.7", 44244, nil},
		{"5.123.4.8", 44244, nil},
	}

	for _, table := range ips {
		asn, err := GetASN(net.ParseIP(table.ip))
		if err != nil { // if there is an error, check if it is the expected error
			if err.Error() != table.err.Error() {
				t.Errorf("Expected error: %s, got: %s", table.err, err)
			} else {
				t.Logf("Expected error: %s, got: %s", table.err, err)
			} // if there is no error, check if the country is the expected one
		} else {
			if asn != table.asn {
				t.Errorf("expected %d, got %d", table.asn, asn)
			}
		}
	}
}
