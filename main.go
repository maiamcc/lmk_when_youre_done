package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var cpuThreshold = flag.Float64("threshold", 100, "CPU threshold to alert on")
var requiredTimesUnder = 8

func psCmdForPID(pid int) string {
	return fmt.Sprintf("ps -p %d -o pcpu | tail +2", pid)
}

func notify(ctx context.Context, msg string) error {
	// terminal-notifier -message basdf
	cmd := exec.CommandContext(ctx, "terminal-notifier", "-message", msg)
	return cmd.Run()

}
func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Need exactly one arg (pid), got: %+v", os.Args[1:])
	}
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Error converting arg %s to int: %v", os.Args[1], err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	flag.Parse()

	var timesUnder int
	var overThreshold bool
	for {
		cmd := exec.CommandContext(ctx, "/bin/sh", "-c", psCmdForPID(pid))
		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("cmd failed with: %v", err)
		}
		outstr := strings.TrimSpace(string(out))

		cpu, err := strconv.ParseFloat(outstr, 64)
		if err != nil {
			log.Fatalf("Error converting cmd output %s to float64: %v", outstr, err)
		}

		if cpu >= *cpuThreshold {
			if !overThreshold {
				// We just crossed the threshold!
				overThreshold = true
			} else {
				// we're still over the threshold, nothing to do
			}
			time.Sleep(time.Second) // zzz then check again
		} else {
			if overThreshold {
				// We just crossed the threshold going down

				// Make sure it's real
				timesUnder += 1
				if timesUnder >= requiredTimesUnder {
					err = notify(ctx, "Okay Crossfire is ready for you!")
					if err != nil {
						log.Fatalf("Error sending notification: %v", err)
					}
					overThreshold = false
					timesUnder = 0
				}
			} else {
				// we're still under the threshold, nothing to do
			}
			time.Sleep(time.Second * 3) // sleep for longer if nothing interesting is happening
		}

	}
}
