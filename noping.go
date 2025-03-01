package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"noping/globals"
	"noping/methods"
)

const headerText = "noping v1.0.0 - Copyright (c) 2025 bastih18"

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		globals.StopProgram()
	}()

	if len(os.Args) < 2 {
		printHelpMenu()
		os.Exit(0)
	}

	ip := os.Args[1]
	port := -1
	count := 65535
	timeout := 1000
	minimal := false
	protocol := "tcp"

	args := os.Args[2:]
	for i, arg := range args {
		switch arg {
		case "-h", "--help":
			printHelpMenu()
			os.Exit(0)
		case "-p", "--port":
			if i+1 < len(args) && args[i+1][0] != '-' {
				r, err := strconv.Atoi(args[i+1])
				if err != nil {
					fmt.Println("❌ The port must be a number")
					os.Exit(0)
				}
				if r < 1 || r > 65535 {
					fmt.Println("❌ The port must be between 1 and 65535")
					os.Exit(0)
				}
				port = r
			} else {
				fmt.Println("❌ You must specify a valid port after", arg)
				os.Exit(0)
			}
		case "-c", "--count":
			if i+1 < len(args) && args[i+1][0] != '-' {
				r, err := strconv.Atoi(args[i+1])
				if err != nil || r < 1 {
					fmt.Println("❌ The count must be a number greater than 0")
					os.Exit(0)
				}
				count = r
			} else {
				fmt.Println("❌ You must specify a valid count after", arg)
				os.Exit(0)
			}
		case "-t", "--timeout":
			if i+1 < len(args) && args[i+1][0] != '-' {
				r, err := strconv.Atoi(args[i+1])
				if err != nil || r < 1 {
					fmt.Println("❌ The timeout must be a number greater than 0")
					os.Exit(0)
				}
				timeout = r
			} else {
				fmt.Printf("❌ You must specify a timeout when using the %s argument\n", arg)
				os.Exit(0)
			}
		case "-m", "--minimal":
			minimal = true
		case "--proto":
			if i+1 < len(args) && args[i+1][0] != '-' {
				if args[i+1] == "tcp" || args[i+1] == "udp" {
					protocol = strings.ToLower(args[i+1])
				} else {
					fmt.Println("❌ The protocol must be either 'tcp' or 'udp'")
					os.Exit(0)
				}
			} else {
				fmt.Printf("❌ You must specify a protocol when using the %s argument\n", arg)
				os.Exit(0)
			}
		}
	}
	if port == -1 {
		protocol = "icmp"
	}

	boldYellow := color.New(color.Bold, color.FgYellow)

	fmt.Println(headerText)
	fmt.Println("")
	rawIp := getIpFromDomain(ip)
	if minimal {
		fmt.Printf("Connecting to %s on %s %s:\n\n", boldYellow.Sprint(ip), boldYellow.Sprint(strings.ToUpper(protocol)), boldYellow.Sprint(map[bool]string{true: strconv.Itoa(port), false: ""}[protocol != "icmp"]))
	} else {
		asnInfo, geoInfo := getASNGeoInfo(rawIp.String())
		fmt.Printf("🌍 Target: %s (%s)\n", map[bool]string{true: fmt.Sprintf("%s [%s]", boldYellow.Sprint(ip), boldYellow.Sprint(rawIp)), false: boldYellow.Sprint(rawIp)}[rawIp.String() != ip], color.BlueString(asnInfo))
		fmt.Printf("📍 Location: %s\n", boldYellow.Sprint(map[bool]string{true: geoInfo.City + ", " + geoInfo.Region + ", " + geoInfo.Country, false: "Unknown"}[geoInfo.City != "nil" && geoInfo.City != ""]))
		fmt.Printf("🔄 Reverse DNS: %s\n\n", boldYellow.Sprint(getReverseDNS(rawIp.String())))
	}

	if protocol == "icmp" {
		methods.ICMPPing(rawIp, count, timeout)
	} else if protocol == "tcp" {
		methods.TCPPing(rawIp.String(), port, count, time.Duration(timeout)*time.Millisecond)
	} else if protocol == "udp" {
		methods.UDPPing(rawIp.String(), port, count, time.Duration(timeout)*time.Millisecond)
	}
}
