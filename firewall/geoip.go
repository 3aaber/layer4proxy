package firewall

import (
	"fmt"
	"layer4proxy/env"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

var (
	dbCountry *geoip2.Reader
	dbIsp     *geoip2.Reader
)

func checkfileexist(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func Setupdb(dbCountryAddress, dbIspAddress string) error { // address is optional and can be used to override the default database file
	var err error
	// Setup Country Database
	if dbCountryAddress == "" { // if no address is given, use the default database file from the config file (if it exists)
		dbCountryAddress = env.GeoipDatabase // default database file is in the config file
	}
	if !checkfileexist(dbCountryAddress) { // if the database file does not exist, return an error
		log.Println("Geoip file not found")
		return fmt.Errorf("geoip file not found")
	}

	dbCountry, err = geoip2.Open(dbCountryAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Setup Isp Database
	if dbIspAddress == "" { // if no address is given, use the default database file from the config file (if it exists)
		dbIspAddress = env.GeoipDatabase // default database file is in the config file
	}
	if !checkfileexist(dbIspAddress) { // if the database file does not exist, return an error
		log.Println("Geoip file not found")
		return fmt.Errorf("geoip file not found")
	}

	dbIsp, err = geoip2.Open(dbIspAddress)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Closedb() {
	dbCountry.Close()
}

func GetCountry(ip net.IP) (string, error) {

	if dbCountry == nil {
		Setupdb("", "")
	}

	// If you are using strings that may be invalid, check that ip is not nil

	if ip == nil {
		return "", fmt.Errorf("invalid ip")
	}

	record, err := dbCountry.City(ip)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return record.Country.IsoCode, nil

}

func GetASN(ip net.IP) (uint64, error) {

	if dbIsp == nil {
		Setupdb("", "")
	}

	if ip == nil {
		return 0, fmt.Errorf("invalid ip")
	}

	record, err := dbIsp.ASN(ip)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return uint64(record.AutonomousSystemNumber), nil

}
