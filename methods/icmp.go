package methods

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Bastih18/NoPing/globals"

	"github.com/fatih/color"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"

	ps "github.com/Bastih18/NoPing/printStats"
)

func ICMPPing(host net.Addr, count int, timeoutMs int) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Printf("Error creating listener: %v\n", err)
		return
	}
	defer conn.Close()

	// Use lower 16 bits of the process ID as identifier.
	pid := os.Getpid() & 0xffff
	payload := bytes.Repeat([]byte("PING"), 14) // Create a 56-byte payload

	greenBold := color.New(color.Bold, color.FgGreen)
	redBold := color.New(color.Bold, color.FgRed)
	normalGreen := color.New(color.FgGreen)

	stats := globals.Stats{}

	for i := 0; i < count; i++ {
		if globals.Stop {
			break
		}

		stats.Attempted++
		message := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   pid,
				Seq:  i,
				Data: payload,
			},
		}

		messageBytes, err := message.Marshal(nil)
		if err != nil {
			fmt.Printf("Error marshalling message: %v\n", err)
			continue
		}

		flushBuffer(conn)
		startTime := time.Now()
		_, err = conn.WriteTo(messageBytes, host)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			continue
		}

		if err = conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMs) * time.Millisecond)); err != nil {
			fmt.Printf("Error setting read deadline: %v\n", err)
			continue
		}

		var n int
		var peer net.Addr
		reply := make([]byte, 1500)
		for {
			n, peer, err = conn.ReadFrom(reply)
			if err != nil {
				stats.Failed++
				fmt.Printf("%s", redBold.Sprintf("❌ Request timeout for icmp_seq %d\n", i))
				break
			}

			parsedMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
			if err != nil {
				continue
			}

			if parsedMessage.Type == ipv4.ICMPTypeEchoReply {
				if echo, ok := parsedMessage.Body.(*icmp.Echo); ok {
					if echo.ID == pid && echo.Seq == i {
						elapsed := time.Since(startTime)
						stats.Connected++
						stats.TotalTime += elapsed

						if stats.Connected == 1 || elapsed < stats.Min {
							stats.Min = elapsed
						}
						if elapsed > stats.Max {
							stats.Max = elapsed
						}

						fmt.Printf("✅ Connected to %v: icmp_seq=%s time=%s\n",
							greenBold.Sprint(peer),
							normalGreen.Sprintf("%d", i),
							greenBold.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6))
						break
					}
				}
			}
		}
		if i < count-1 {
			select {
			case <-time.After(time.Second):
			case <-globals.StopChan:
				break
			}
		}
	}
	ps.PrintStats("ICMP", host.String(), 0, stats)
}

func flushBuffer(conn *icmp.PacketConn) {
	conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
	buf := make([]byte, 1500)
	for {
		_, _, err := conn.ReadFrom(buf)
		if err != nil {
			break
		}
	}
	conn.SetReadDeadline(time.Time{})
}
