package methods

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Bastih18/NoPing/globals"
	ps "github.com/Bastih18/NoPing/printStats"
	"github.com/fatih/color"
)

// Define a list of protocols and their corresponding ports and packet formats
var udpMessages = []struct {
	Protocol string
	Port     int
	Packet   []byte
}{
	{"dns", 53, []byte{
		0xaa, 0xbb, // Transaction ID
		0x01, 0x00, // Standard query
		0x00, 0x01, // One question
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x03, 'w', 'w', 'w',
		0x06, 'g', 'o', 'o', 'g', 'l', 'e',
		0x03, 'c', 'o', 'm', 0x00, // "www.google.com"
		0x00, 0x01, // Type A query
		0x00, 0x01, // Class IN (Internet)
	}},
	{"ntp", 123, []byte{
		0x1b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}},
	{"snmp", 161, []byte{
		0x30, 0x26, 0x02, 0x01, 0x00, 0x04, 0x14, 0x2b, 0x06, 0x01,
		0x04, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}},
	{"upnp", 1900, []byte{
		0x00, 0x00, 0x00, 0x00, // Example SSDP request for UPnP
	}},
	{"steam", 27030, []byte{
		0xFF, 0xFF, 0xFF, 0xFF, // Magic prefix for Steam query
		0x11,                   // Steam query request type
		0x00, 0x00, 0x00, 0x00, // Empty header
	}},
	{"syslog", 514, []byte{
		0x3c, 0x01, // Header
		0x00, 0x00, // Facility and Severity
		0x00, 0x00, 0x00, 0x00, // Time and Hostname
		0x00, 0x00, 0x00, 0x00, // AppName and Message
	}},
	{"tftp", 69, []byte{
		0x00, 0x01, // Request Packet
		0x00, 0x00, 0x00, 0x01, // Request to read file
	}},
	{"radius", 1812, []byte{
		0x01, 0x00, 0x00, 0x06, // Request type
		0x00, 0x00, 0x00, 0x01, // Identifier and length
		// More RADIUS-specific fields...
	}},
	{"chargen", 19, []byte{0x00}}, // Simple message
	{"nfs", 2049, []byte{
		0x00, 0x00, 0x00, 0x00, // NFS request packet
	}},
	{"sip", 5060, []byte{
		0x00, 0x01, 0x00, 0x00, // Request for SIP service
		// Further SIP message fields...
	}},
	{"dhcp", 67, []byte{
		0x01, 0x01, 0x06, 0x00, // DHCP Discover
		0x00, 0x00, 0x00, 0x00, // Client hardware address
		// Further DHCP fields...
	}},
	{"bgp", 179, []byte{
		0x00, 0x01, 0x00, 0x00, // Open message fields
	}},
}

// Auto-detect protocol based on port
func autoDetectProtocol(port int) string {
	for _, entry := range udpMessages {
		if entry.Port == port {
			return entry.Protocol
		}
	}
	return "" // Unknown protocol for the given port
}

// Generalized function to send UDP packets based on protocol
func UDPPing(host string, port int, timeout time.Duration, protocol string, count int) {
	// If protocol is "auto", detect it based on the port
	if protocol == "auto" {
		protocol = autoDetectProtocol(port)
		if protocol == "" {
			fmt.Printf("❌ Unknown protocol for port %d\n", port)
			return
		}
	}

	// Convert protocol to lowercase to ensure consistency
	protocol = strings.ToLower(protocol)

	// Find the matching protocol in the udpMessages list
	var packet []byte
	for _, entry := range udpMessages {
		if entry.Protocol == protocol {
			packet = entry.Packet
			break
		}
	}

	// Check if the protocol is valid
	if packet == nil {
		fmt.Printf("❌ Unknown protocol: %s\n", protocol)
		return
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Printf("❌ Error resolving UDP address: %v\n", err)
		return
	}

	// Colorized outputs
	greenBold := color.New(color.Bold, color.FgGreen)
	redBold := color.New(color.Bold, color.FgRed)

	// Stats tracking
	stats := globals.Stats{}

	// Loop to send multiple pings
	for i := 0; i < count; i++ {
		if globals.Stop {
			break
		}
		stats.Attempted++
		start := time.Now()

		// Timeout context for this UDP request
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		go func() {
			select {
			case <-globals.StopChan:
				cancel()
				break
			case <-ctx.Done():
			}
		}()

		// Bind to a local address
		localAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:0")
		conn, err := net.DialUDP("udp", localAddr, udpAddr)
		elapsed := time.Since(start)
		defer cancel()

		// Check if the connection was successful
		if err != nil {
			stats.Failed++
			fmt.Printf("%s\n", redBold.Sprint("❌ Connection timed out"))
		} else {
			// Send the query packet
			_, err = conn.Write(packet)
			if err != nil {
				stats.Failed++
				fmt.Printf("%s\n", redBold.Sprint("❌ Failed to send UDP packet"))
				continue
			}

			// Set read deadline for response
			conn.SetReadDeadline(time.Now().Add(timeout))
			buffer := make([]byte, 1024)
			n, _, err := conn.ReadFrom(buffer)

			// If we received any data, it's valid and the service is online
			if err == nil && n > 0 {
				stats.Connected++
				stats.TotalTime += elapsed

				// Track min/max times
				if stats.Connected == 1 || elapsed < stats.Min {
					stats.Min = elapsed
				}
				if elapsed > stats.Max {
					stats.Max = elapsed
				}

				// Convert time to milliseconds for correct output
				fmt.Printf("✅ UDP ping to %s: time=%.3fms protocol=%s port=%s\n",
					greenBold.Sprint(host),
					float64(elapsed.Nanoseconds())/1e6, // Convert nanoseconds to milliseconds
					greenBold.Sprint(protocol),
					greenBold.Sprint(strconv.Itoa(port)),
				)
			} else {
				// If no response is received but the connection succeeded, treat it as "open"
				stats.Connected++
				fmt.Printf("✅ UDP ping to %s: port open but no response received\n", greenBold.Sprint(address))
			}
		}

		// Wait before sending the next packet
		if i < count-1 {
			select {
			case <-time.After(time.Second):
			case <-globals.StopChan:
				break
			}
		}
	}

	// Print statistics after the loop
	ps.PrintStats("UDP", host, port, stats)
}
