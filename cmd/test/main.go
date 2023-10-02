package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"winsleepd"
)

var (
	ActionSleep      = "sleep"
	ActionHibernate  = "hibernate"
	ActionScreenOff  = "screenoff"
	ActionLockScreen = "lockscreen"
)

type Configuration struct {
	Timeout string `json:"timeout"`
	Action  string `json:"action"`
}

func GetConfiguration() (*Configuration, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(dir + "/.winsleepd.json")
	if err != nil {
		if os.IsNotExist(err) {
			return NewConfiguration()
		}
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func NewConfiguration() (*Configuration, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filepath.Join(dir, ".winsleepd.json"))
	if err != nil {
		return nil, err
	}
	encoder := json.NewEncoder(file)
	configuration := Configuration{
		Timeout: "30m", // 30 minutes
		Action:  ActionScreenOff,
	}
	encoder.SetIndent("", "\t")
	err = encoder.Encode(configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

type daemon struct {
	mouseX     int
	mouseY     int
	lastChange time.Time
	lastLogMsg string
}

func main() {
	m := &daemon{}
	for {
		x, y := winsleepd.GetMousePos()
		if x != m.mouseX || y != m.mouseY {
			m.mouseX = x
			m.mouseY = y
			m.lastChange = time.Now()
		}
		log.Println("Query mouse position:", x, y)
		log.Println("Last change:", m.lastChange)
		config, err := GetConfiguration()
		if err != nil {
			log.Println(1, fmt.Sprintf("error getting configuration: %v", err))
		}
		timeout, err := time.ParseDuration(config.Timeout)
		if err != nil {
			log.Println(1, fmt.Sprintf("error parsing timeout: %v", err))
		}
		fmt.Println("Timeout:", timeout)
		fmt.Println("Time till action:", time.Now().Sub(m.lastChange))
		if time.Now().Sub(m.lastChange) > timeout {
			if time.Now().Sub(m.lastChange) > timeout {
				switch config.Action {
				case ActionSleep:
					logmsg := "Sleeping..."
					if logmsg != m.lastLogMsg {
						m.lastLogMsg = logmsg
						log.Println(1, logmsg)
					}
					winsleepd.Sleep()
				case ActionHibernate:
					logmsg := "Hibernating..."
					if logmsg != m.lastLogMsg {
						m.lastLogMsg = logmsg
						log.Println(1, logmsg)
					}
					winsleepd.Hibernate()
				case ActionScreenOff:
					logmsg := "Turning screen off..."
					if logmsg != m.lastLogMsg {
						m.lastLogMsg = logmsg
						log.Println(1, logmsg)
					}
					winsleepd.ScreenOff()
				case ActionLockScreen:
					logmsg := "Locking Screen..."
					if logmsg != m.lastLogMsg {
						m.lastLogMsg = logmsg
						log.Println(1, logmsg)
					}
					winsleepd.LockScreen()
				}
			}
		}
		time.Sleep(5 * time.Second)
	}

}
