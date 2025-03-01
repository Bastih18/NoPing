package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Bastih18/NoPing/globals"
	"github.com/fatih/color"
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

func getLatestVersion(current string) (string, string, string) {
	const url = "https://api.github.com/repos/Bastih18/NoPing/tags"
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "Could not fetch", color.YellowString("are maybe"), "the latest"
	}
	defer resp.Body.Close()
	var tags []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil || len(tags) == 0 {
		return "Could not decode", color.YellowString("are maybe"), "the latest"
	}
	latest := tags[0].Name
	if current == latest {
		return latest, color.GreenString("are"), "the latest"
	}
	if current == "dev" {
		return latest, color.GreenString("are"), color.CyanString("the dev")
	}
	const pattern = `^v\d+\.\d+\.\d+$`
	re := regexp.MustCompile(pattern)
	if re.MatchString(current) && re.MatchString(latest) {
		if cmp := compareVersions(current, latest); cmp < 0 {
			return latest, color.RedString("are"), color.RedString("an outdated")
		} else if cmp > 0 {
			return latest, color.GreenString("are"), color.CyanString("an ahead")
		}
	}
	return latest, color.YellowString("are maybe"), "unable to compare versions"
}

func compareVersions(v1, v2 string) int {
	p1, p2 := strings.Split(v1[1:], "."), strings.Split(v2[1:], ".")
	for i := 0; i < 3; i++ {
		n1, _ := strconv.Atoi(p1[i])
		n2, _ := strconv.Atoi(p2[i])
		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}
	return 0
}
