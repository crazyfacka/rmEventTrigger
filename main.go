package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/phuslu/log"
)

// Build flags
var (
	BuildTime  string
	CommitHash string
	GoVersion  string
)

type Event int

const (
	Sync Event = iota
	Connection
)

var eventType = map[string]Event{
	"Sync":       Sync,
	"Connection": Connection,
}

var mapConfig map[Event]EventConfig

type EventConfig struct {
	Event       string   `json:"Event"`
	Actions     []string `json:"Actions"`
	LastTrigger time.Time
}

type Config struct {
	Conf []EventConfig `json:"conf"`
}

func executeEvent(ch <-chan Event) {
	for ev := range ch {
		if v, ok := mapConfig[ev]; ok {
			lastTriggerDiff := time.Since(v.LastTrigger)
			if lastTriggerDiff.Seconds() < 30 {
				log.Info().Str("event", v.Event).Str("lastTrigger", lastTriggerDiff.String()).Msg("Last event happend shortly before")
				continue
			}

			for _, item := range v.Actions {
				log.Info().Str("cmd", item).Msg("Executing command")
				cmdStart := time.Now()

				splitCmd := strings.Fields(item) // TODO Optimize this
				cmd := exec.Command(splitCmd[0], splitCmd[1:]...)

				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Error().Err(err).Msg("Error opening stdout pipe")
					return
				}

				if err := cmd.Start(); err != nil {
					log.Error().Err(err).Msg("Error running command")
					return
				}

				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					line := scanner.Text()
					log.Debug().Msg(line)
				}

				if err := scanner.Err(); err != nil {
					log.Error().Err(err).Msg("Error reading output")
				}

				cmd.Process.Kill()

				log.Info().Str("duration", time.Since(cmdStart).String()).Msg("Execution complete")
			}

			v.LastTrigger = time.Now()
			mapConfig[ev] = v
		}
	}
}

func loadConfiguration(confFile string) (map[Event]EventConfig, error) {
	file, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	mapConfig = make(map[Event]EventConfig)
	for _, val := range config.Conf {
		if _, ok := eventType[val.Event]; !ok {
			return nil, fmt.Errorf("event type not available: %s", val.Event)
		}

		if _, ok := mapConfig[eventType[val.Event]]; !ok {
			mapConfig[eventType[val.Event]] = val
		} else {
			tmpEvent := mapConfig[eventType[val.Event]]
			tmpEvent.Actions = append(mapConfig[eventType[val.Event]].Actions, val.Actions...)
			mapConfig[eventType[val.Event]] = tmpEvent
		}
	}

	return mapConfig, nil
}

func showHelpAndExit() {
	fmt.Printf("Usage: %s [OPTIONS]\n", filepath.Base(os.Args[0]))
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()

	os.Exit(0)
}

func main() {
	log.Info().Str("BuildTime", BuildTime).Str("CommitHash", CommitHash).Str("GoVersion", GoVersion).Msg("Starting rmEventTrigger")

	help := flag.Bool("h", false, "Whether to show help output or not")
	debug := flag.Bool("d", false, "Debug mode (will tail ./debug.log)")
	confFile := flag.String("c", "conf.json", "Path to configuration file")
	flag.Parse()

	if *help {
		showHelpAndExit()
	}

	log.Info().Str("confFile", *confFile).Msg("Loading config file")
	_, err := loadConfiguration(*confFile)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding configuration file")
		return
	}
	log.Info().Str("confFile", *confFile).Msg("Loaded config file")

	var cmd *exec.Cmd
	if *debug {
		cmd = exec.Command("tail", "-f", "./debug.log")
	} else {
		cmd = exec.Command("journalctl", "-f")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error().Err(err).Msg("Error opening stdout pipe")
		return
	}

	if err := cmd.Start(); err != nil {
		log.Error().Err(err).Msg("Error starting command to listen in to events")
		return
	}

	syncRe, _ := regexp.Compile(`.*rm.synchronizer.*execute ended.*`)

	evs := make(chan Event, 5)
	go executeEvent(evs)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if syncRe.Match([]byte(line)) {
			evs <- Sync
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("Error reading output")
	}

	cmd.Process.Kill()
}
