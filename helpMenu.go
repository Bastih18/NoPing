package main

import "fmt"

func printHelpMenu() {
	fmt.Println(headerText)
	fmt.Println("")
	fmt.Println("Syntax: noping <ip> [OPTIONS]")
	fmt.Println("")
	fmt.Println("ARGS:")
	fmt.Println("    ip: IP address to ping")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("    -h, --help                     Print the Help Menu")
	fmt.Println("    -p, --port <port>              port to ping (default: ICMP) [OPTIONAL]")
	fmt.Println("    -c, --count <count>            number of pings (default: 65535) [OPTIONAL]")
	fmt.Println("    -t, --timeout <timeout>        timeout in milliseconds (default: 1000) [OPTIONAL]")
	fmt.Println("    -m, --minimal                  print only the minimum information (default: false) [OPTIONAL]")
	fmt.Println("    -v, --version                  print detailed version information")
	fmt.Println("    --proto <protocol (tcp/udp)>   protocol to use (default: tcp) [OPTIONAL]")
	fmt.Println("    --update [version]             update noping to the specified version (when empty, it updates to the latest version)")
	fmt.Println("")
}
