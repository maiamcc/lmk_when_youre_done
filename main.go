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
var requiredTimesUnder = flag.Int("times-under", 8, "number of readings under the threshold before deciding that your work is probably done")
var sleepActive = flag.Duration("sleep-active", time.Second, "time to sleep between readings when waiting for work to end")
var sleepPassive = flag.Duration("sleep-passive", time.Second * 5, "time to sleep between readings when waiting for work to start")

func psCmdForPID(pid int) string {
	return fmt.Sprintf("ps -p %d -o pcpu | tail +2", pid)
}

func notify(ctx context.Context, msg string) error {
	cmd := exec.CommandContext(ctx, "terminal-notifier", "-message", msg)
	return cmd.Run()

}
func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("Need exactly one positional arg (pid), got: %+v", os.Args[1:])
	}
	pid, err := strconv.Atoi(flag.Args()[0])
	if err != nil {
		log.Fatalf("Error converting arg %s to int: %v", os.Args[1], err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


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
			time.Sleep(*sleepActive) // zzz then check again
		} else {
			if overThreshold {
				// We just crossed the threshold going down

				// Make sure it's real
				timesUnder += 1
				if timesUnder >= *requiredTimesUnder {
					err = notify(ctx, "Okay, your thing is ready for you!")
					if err != nil {
						log.Fatalf("Error sending notification: %v", err)
					}
					overThreshold = false
					timesUnder = 0
				}
			} else {
				// we're still under the threshold, nothing to do
			}
			time.Sleep(*sleepPassive)
		}

	}
}
