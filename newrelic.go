package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// TODO: Read the URI and the license key for environment variables - DONE
// TODO: Start a timer and flush events after timer ends
// TODO: When the counter hits a batch size i.e.  200 or the timer hits flush time - send to New Relic
// TODO: Read events and send events on seperate threads?
// TODO: Check Logs buffer size - API will reject payloads over 1MB
// TODO: Write to a log file for INFO and DEBUG messages
// TODO: Implement gzip compression on JSON

func init() {
	// setup logging
	log.SetOutput(os.Stdout)
}

// how many logs will we hold in the buffer before transmitting?
const bufferSize = 1

// how long in ms before we transmit the logs if the buffer doesn't fill up first
const flushMillis = 2000

// LogPayload - a struct to hold the payload for an individual log message
type LogPayload struct {
	Logs [bufferSize]json.RawMessage `json:"logs"`
}

// func readEnvironmentVariable(key string) string {
// 	value := os.Getenv(key)
// 	if len(value) == 0 {
// 		log.Error(key + " environment variable is not set!")
// 		panic(key + " environment variable is not set!")
// 	} else {
// 		return value
// 	}
// }

func main() {
	// Parse command line args
	args := parseArgs()

	// a counter to keep track of how many logs in the buffer
	var logCounter int
	// an array to store the JSON strings we are sent by NXLog om_exec module
	var logBuffer [bufferSize]json.RawMessage

	// reader for stdin
	// reader := bufio.NewReader(os.Stdin)
	// read a string from the reader delimited by newline char
	// message, _ := reader.ReadString('\n')

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		// log the message for debugging purposes
		log.Info(message)
		rawIn := json.RawMessage(message)

		// store the JSON string, log, in the logs array
		logBuffer[logCounter] = rawIn

		// increment the log counter
		logCounter++
		log.Info("incremented the buffer counter to: " + strconv.Itoa(logCounter))

		if logCounter == bufferSize {
			var lp = LogPayload{Logs: logBuffer}
			send(lp, args[0], args[1])
			// reset the log counter back to 0
			logCounter = 0
		}
	}
}

func parseArgs() [2]*string {
	log.Info("Parsing command line args")

	// Define two command line flags
	logAPIURIPtr := flag.String("NEW_RELIC_LOG_URI", "https://log-api.newrelic.com/log/v1", "New Relic Log API URI")
	licenseKeyPtr := flag.String("NEW_RELIC_LICENSE_KEY", "", "New Relic License Key")
	flag.Parse()

	if len(*licenseKeyPtr) == 0 {
		log.Info("--NEW_RELIC_LICENSE_KEY is a mandatory argument")
		panic("--NEW_RELIC_LICENSE_KEY is a mandatory argument")
	}

	if *logAPIURIPtr == "https://log-api.newrelic.com/log/v1" {
		log.Info("Using US Logging API Endpoint")
	} else {
		log.Info("Using EU Logging API Endpoint")
	}

	var args [2]*string
	args[0] = logAPIURIPtr
	args[1] = licenseKeyPtr

	log.Info("Finished parsing args")
	return args
}

func send(logStruct LogPayload, logAPIURI *string, licenseKey *string) {
	log.Info("sending payload of size:", len(logStruct.Logs))

	// API expects an array of LogPayload
	var array []LogPayload
	array = append(array, logStruct)

	// Marshal the struct into JSON
	arr, err := json.Marshal(array)
	if err != nil {
		log.Error(err)
	} else {
		log.Info("parsed JSON: " + string(arr))
	}

	// Read API_URI and LICENSE_KEY from env vars
	// logAPIURI := readEnvironmentVariable("NEW_RELIC_LOG_URI")
	// licenseKey := readEnvironmentVariable("NEW_RELIC_LICENSE_KEY")

	req, err := http.NewRequest("POST", *logAPIURI, bytes.NewBuffer(arr))
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-License-Key", *licenseKey)

	client := &http.Client{}
	resp, httpErr := client.Do(req)
	if httpErr != nil {
		log.Info("There was an HTTP error", httpErr)
		log.Fatal(httpErr)
	} else {
		log.Info("HTTP response code: ", strconv.Itoa(resp.StatusCode))
	}
}
