package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"noping/globals"
	"os"
)

func getASNGeoInfo(ip string) (string, globals.GeoInfo) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return "Unknown ASN", globals.GeoInfo{City: "Unknown", Region: "Unknown", Country: "Unknown"}
	}
	defer resp.Body.Close()

	var data struct {
		Org     string `json:"org"`
		City    string `json:"city"`
		Region  string `json:"region"`
		Country string `json:"country"`
		Bogon   bool   `json:"bogon"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)
	if data.Bogon {
		return "Reserved IP", globals.GeoInfo{City: "nil", Region: "nil", Country: "nil"}
	}
	return data.Org, globals.GeoInfo{City: data.City, Region: data.Region, Country: data.Country}
}

func getReverseDNS(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "No hostname found"
	}
	return names[0]
}

func getIpFromDomain(host string) net.Addr {
	ips, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Println("❌ Could not resolve host:", err)
		os.Exit(1)
		return nil
	}
	return ips
}
