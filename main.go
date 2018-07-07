package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

const timeFormat = "2006-01-02_15-04-05"

var debug *bool

type Notifications struct {
	ScanToNotifications struct {
		ScanToDeviceDisplay string
		ScanToHostID        string
		ScanToNotSetup      int
		ADFLoaded           int
	}
	StartScanNotifications struct {
		StartScan int
		ADFLoaded int
	}

	FaxNotifications struct {
		FaxReceiveFunction int
		FaxPrinting        int
		LastFaxLogEntry    struct {
			EntryID    int
			Type       int
			FaxNumber  string
			TimeDate   string
			NumPages   int
			ResultCode int
		}
		FaxMasterHostID       string
		FaxUploadState        int
		FaxLogChangeIndicator int
		FaxForwardEnabled     int
		FaxForwardNumber      string
	}
}

func notifications(address string) (n Notifications, err error) {
	res, err := http.Get(address)
	if err != nil {
		return
	}

	dec := xml.NewDecoder(res.Body)
	err = dec.Decode(&n)
	return
}

func scan(path string, options []string) error {
	c := exec.Command("scanimage", options...)
	c.Dir = path
	if *debug {
		log.Println(c)
	}
	return c.Run()
}

func mkdir(prefix string) (string, error) {
	p := path.Join(prefix, time.Now().Format(timeFormat))
	err := os.Mkdir(p, 0700)
	if err != nil {
		return "", err
	}

	return p, nil
}

type config struct {
	Address  string
	Sleep    string
	MaxSleep string
	Path     string
	Sane     struct {
		Flatbed []string
		ADF     []string
	}
}

func (c *config) load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(c)
	return err
}

func (c config) write(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(c)
	return err
}

func main() {
	configpath := flag.String("config", "config.json", "config file")
	example := flag.Bool("example", false, "write example config")
	debug = flag.Bool("debug", false, "debug output")
	flag.Parse()

	if *configpath == "" {
		log.Fatal("no config")
	}

	if *example {
		err := config{}.write(*configpath)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	c := &config{}
	err := c.load(*configpath)
	if err != nil {
		log.Fatal(err)
	}

	sleep, err := time.ParseDuration(c.Sleep)
	if err != nil {
		log.Fatal(err)
	}

	maxsleep, err := time.ParseDuration(c.MaxSleep)
	if err != nil {
		log.Fatal(err)
	}

	currentsleep := sleep

	for {
		if currentsleep > maxsleep {
			log.Fatal("max sleep reached")
		}

		time.Sleep(currentsleep)

		n, err := notifications(c.Address)
		if err != nil {
			currentsleep = currentsleep * 2
			if *debug {
				log.Println(err)
				log.Println("backing off to %v", currentsleep)
			}
			continue
		}

		currentsleep = sleep

		if n.StartScanNotifications.StartScan == 1 {
			p, err := mkdir(c.Path)
			if err != nil {
				log.Fatal(err)
			}

			options := c.Sane.Flatbed
			if n.StartScanNotifications.ADFLoaded == 1 {
				options = c.Sane.ADF
			}

			err = scan(p, options)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
