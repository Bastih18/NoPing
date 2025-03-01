package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/Bastih18/NoPing/globals"
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

// func getReverseDNS(ip string) string {
// 	resolver := &net.Resolver{
// 		PreferGo: true,
// 		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
// 			d := net.Dialer{Timeout: 2 * time.Second}
// 			return d.DialContext(ctx, "udp", "1.1.1.1:53")
// 		},
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()
// 	names, err := resolver.LookupAddr(ctx, ip)
// 	if err != nil || len(names) == 0 {
// 		return "No hostname found"
// 	}
// 	return names[0]
// }

func getIpFromDomain(host string) net.Addr {
	ips, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Println("❌ Could not resolve host:", err)
		os.Exit(1)
		return nil
	}
	return ips
}
