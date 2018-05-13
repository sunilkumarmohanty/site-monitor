package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/sunilkumarmohanty/site-monitor/monitor"
)

func main() {
	log.Println("Monitor started")
	configFilePath := readEnvironVar("MONITOR_CONFIG_FILE_PATH", "./monitor.config")
	config, err := ioutil.ReadFile(configFilePath)
	handleErr(err)

	logDirPath := readEnvironVar("MONITOR_LOG_DIR_PATH", ".")
	logFilePath := fmt.Sprintf("%v/monitor_%v.log", logDirPath, time.Now().Format("20060102150405"))
	f, err := os.Create(logFilePath)
	handleErr(err)
	defer f.Close()

	m, err := monitor.New(config,
		f,
		readEnvironVar("MONITOR_USER_AGENT", "monitorv1.0.0"),
		readEnvironVar("MONITOR_TIME_OUT", "30"),
	)
	handleErr(err)
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	go func() {
		<-osChan
		log.Println("Interrupted by OS")
		m.Stop()
	}()

	m.Run()
	log.Println("Monitor stopped")

}

func readEnvironVar(key string, defaultVal string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultVal
	}
	return val
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
