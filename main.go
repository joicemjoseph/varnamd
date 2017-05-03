package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

var (
	port                   int
	version                bool
	maxHandleCount         int
	host                   string
	uiDir                  string
	enableInternalApis     bool // internal APIs are not exposed to public
	enableSSL              bool
	certFilePath           string
	keyFilePath            string
	logToFile              bool    // logs will be written to file when true
	varnamdConfig          *config // config instance used across the application
	startedAt              time.Time
	downloadEnabledSchemes string // comma separated list of scheme identifier for which download will be performed
	syncIntervalInSecs     int
	//upstreamURL            string
	syncDispatcherRunning bool
)

const (
	downloadPageSize = 100
)

func getLogsDir() string {
	d := getConfigDir()
	logsDir := path.Join(d, "logs")
	err := os.MkdirAll(logsDir, 0777)
	if err != nil {
		panic(err)
	}

	return logsDir
}

func redirectLogToFile() {
	year, month, day := time.Now().Date()
	logfile := path.Join(getLogsDir(), fmt.Sprintf("%d-%d-%d.log", year, month, day))
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
}

func init() {
	flag.IntVar(&port, "p", 8080, "Run daemon in specified port")
	flag.IntVar(&maxHandleCount, "max-handle-count", 10, "Maximum number of handles can be opened for each language")
	flag.StringVar(&host, "host", "", "Host for the varnam daemon server")
	flag.StringVar(&uiDir, "ui", "", "UI directory path")
	flag.BoolVar(&enableInternalApis, "enable-internal-apis", false, "Enable internal APIs")
	flag.BoolVar(&enableSSL, "enable-ssl", false, "Enables SSL")
	flag.StringVar(&certFilePath, "cert-file-path", "", "Certificate file path")
	flag.StringVar(&keyFilePath, "key-file-path", "", "Key file path")
	//flag.StringVar(&upstreamurl, "upstream", "https://api.varnamproject.com", "Provide an upstream server")
	flag.StringVar(&downloadEnabledSchemes, "enable-download", "", "Comma separated language identifier for which varnamd will download words from upstream")
	flag.IntVar(&syncIntervalInSecs, "sync-interval", 30, "Download interval in seconds")
	flag.BoolVar(&logToFile, "log-to-file", true, "If true, logs will be written to a file")
	flag.BoolVar(&version, "version", false, "Print the version and exit")
}

//func syncRequired() bool {
//return len(varnamdConfig.schemesToDownload) > 0
//}

// Starts the sync process only if it is not running
//func startSyncDispatcher() {
//if syncRequired() && !syncDispatcherRunning {
//sync := newSyncDispatcher(varnamdConfig.syncIntervalInSecs * time.Second)
//sync.start()
//sync.runNow() // run one round of sync immediatly rather than waiting for the next interval to occur
//syncDispatcherRunning = true
//}
//}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	//varnamdConfig = initConfig()
	startedAt = time.Now()
	if version {
		fmt.Println(varnamdVersion)
		os.Exit(0)
	}
	if logToFile {
		redirectLogToFile()
	}

	log.Printf("varnamd %s", varnamdVersion)

	//startSyncDispatcher()
	startDaemon()
}
