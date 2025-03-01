package methods

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/Bastih18/NoPing/globals"
	ps "github.com/Bastih18/NoPing/printStats"

	"github.com/fatih/color"
)

func TCPPing(host string, port int, count int, timeout time.Duration) {
	address := net.JoinHostPort(host, strconv.Itoa(port))

	greenBold := color.New(color.Bold, color.FgGreen)
	redBold := color.New(color.Bold, color.FgRed)

	stats := globals.Stats{}

	for i := 0; i < count; i++ {
		if globals.Stop {
			break
		}
		stats.Attempted++
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		go func() {
			select {
			case <-globals.StopChan:
				cancel()
				break
			case <-ctx.Done():
			}
		}()

		dialer := net.Dialer{}
		conn, err := dialer.DialContext(ctx, "tcp", address)
		elapsed := time.Since(start)
		defer cancel()

		if err != nil {
			stats.Failed++
			fmt.Printf("%s\n", redBold.Sprint("❌ Connection timed out"))
		} else {
			defer conn.Close()

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
				greenBold.Sprint("TCP"),
				greenBold.Sprint(strconv.Itoa(port)),
			)
		}
		if i < count-1 {
			select {
			case <-time.After(time.Second):
			case <-globals.StopChan:
				break
			}
		}
	}
	ps.PrintStats("TCP", host, port, stats)
}
