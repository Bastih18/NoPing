package methods

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"noping/globals"

	ps "noping/printStats"

	"github.com/fatih/color"
)

func UDPPing(host string, port int, count int, timeout time.Duration) {
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("udp", address, timeout)

	greenBold := color.New(color.Bold, color.FgGreen)
	redBold := color.New(color.Bold, color.FgRed)

	if err != nil {
		fmt.Printf("%s\n", redBold.Sprintf("Error dialing UDP: %v", err))
		return
	}
	defer conn.Close()

	stats := globals.Stats{}

	for i := 0; i < count; i++ {
		if globals.Stop {
			break
		}

		stats.Attempted++
		message := []byte("ping")
		start := time.Now()

		_, err = conn.Write(message)
		if err != nil {
			stats.Failed++
			fmt.Printf("%s\n", redBold.Sprint("❌ Connection timed out"))
			time.Sleep(time.Second)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(timeout))
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		elapsed := time.Since(start)

		if err != nil {
			stats.Failed++
			fmt.Printf("❌ UDP ping to %s: timeout: %v\n", redBold.Sprint(address), redBold.Sprint(err))
		} else {
			stats.Connected++
			stats.TotalTime += elapsed

			if stats.Connected == 1 || elapsed < stats.Min {
				stats.Min = elapsed
			}
			if elapsed > stats.Max {
				stats.Max = elapsed
			}

			fmt.Printf("✅ Connected to %s: time=%s protocol=%s port=%s\n",
				greenBold.Sprint(host),
				greenBold.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
				greenBold.Sprint("UDP"),
				greenBold.Sprint(strconv.Itoa(port)))
		}
		if i < count-1 {
			select {
			case <-time.After(time.Second):
			case <-globals.StopChan:
				break
			}
		}
	}
	ps.PrintStats("UDP", host, port, stats)
}
