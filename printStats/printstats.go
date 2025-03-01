package printstats

import (
	"fmt"
	"time"

	"noping/globals"

	"github.com/fatih/color"
)

func PrintStats(proto, host string, port int, stats globals.Stats) {
	if stats.Attempted == 0 {
		return
	}
	greenBold := color.New(color.Bold, color.FgGreen)
	blueBold := color.New(color.Bold, color.FgBlue)
	percentage := (float64(stats.Failed) / float64(stats.Attempted)) * 100
	greenBold.Println("\nConnection statistics:")
	if port != 0 {
		fmt.Printf("    Protocol = %s, Host = %s, Port = %s\n", blueBold.Sprint(proto), blueBold.Sprint(host), blueBold.Sprint(port))
	} else {
		fmt.Printf("    Protocol = %s, Host = %s\n", blueBold.Sprint(proto), blueBold.Sprint(host))
	}
	fmt.Printf("    Attempted = %s, Connected = %s, Failed = %s (%s)\n",
		blueBold.Sprint(stats.Attempted),
		blueBold.Sprint(stats.Connected),
		blueBold.Sprint(stats.Failed),
		blueBold.Sprintf("%.2f%%", percentage))
	if stats.Connected > 0 {
		avg := stats.TotalTime / time.Duration(stats.Connected)
		fmt.Printf("Approximate connection times:\n")
		fmt.Printf("    Minimum = %s, Maximum = %s, Average = %s\n",
			blueBold.Sprintf("%.2fms", float64(stats.Min.Nanoseconds())/1e6),
			blueBold.Sprintf("%.2fms", float64(stats.Max.Nanoseconds())/1e6),
			blueBold.Sprintf("%.2fms", float64(avg.Nanoseconds())/1e6),
		)
	}
}
