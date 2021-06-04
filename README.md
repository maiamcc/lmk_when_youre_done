# LMK When You're Done

Say you've got some process that occasionally does some prolonged, boring, CPU-intensive work. Ideally you want to be able to start this work, tab away and do something more interesting, and come back when the work is done. This util will monitor the given process for CPU usage over a given threshold (presumably the work) and notify you once CPU usage is back below the threshold (presumably when the work is done).

What CPU levels do/don't indicate the "work" will vary based on process, machine, etc. Play with the levers until you've found a combination of settings that work for you.

## Requirements
* MacOS
* [Terminal Notifier](https://github.com/julienXX/terminal-notifier) (`brew install terminal-notifier`)

## Usage
* find the PID of the process to monitor (e.g. via `ps -ef | grep MyProcessName`)
* invoke this util: `go run main.go [flags] <pid>`
* stuck? `go run main.go --help`
